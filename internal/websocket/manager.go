package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"messenger/internal/models"
	"messenger/pkg/msgpack"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене нужно проверять origin
	},
}

type WSManager struct {
	clients      sync.Map // Map[userID]*Client
	broadcast    chan *WSMessage
	register     chan *Client
	unregister   chan *Client
	mu           sync.RWMutex
}

type Client struct {
	userID        string
	conn          *websocket.Conn
	send          chan []byte
	subChats      map[string]bool
	subUsers      map[string]bool
	lastHeartbeat time.Time
	mu            sync.Mutex
}

type WSMessage struct {
	Type      string      `msgpack:"type"`
	Payload   interface{} `msgpack:"payload,omitempty"`
	RequestID string      `msgpack:"request_id,omitempty"`
	UserID    string      `msgpack:"user_id,omitempty"`
}

func NewWSManager() *WSManager {
	ws := &WSManager{
		broadcast:  make(chan *WSMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go ws.run()
	return ws
}

func (ws *WSManager) run() {
	for {
		select {
		case client := <-ws.register:
			ws.clients.Store(client.userID, client)
		case client := <-ws.unregister:
			ws.clients.Delete(client.userID)
			close(client.send)
		case message := <-ws.broadcast:
			ws.handleBroadcast(message)
		}
	}
}

func (ws *WSManager) handleBroadcast(msg *WSMessage) {
	ws.clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		client.mu.Lock()
		defer client.mu.Unlock()

		// Отправляем только подписанным пользователям
		if msg.Type == "new_message" {
			if payload, ok := msg.Payload.(models.WSNewMessage); ok {
				if !client.subChats[payload.ChatID] {
					return true
				}
			}
		}

		if msg.Type == "user_status" {
			if payload, ok := msg.Payload.(models.WSUserStatus); ok {
				if !client.subUsers[payload.UserID] && payload.UserID != client.userID {
					return true
				}
			}
		}

		data, err := msgpack.Encode(msg)
		if err != nil {
			return true
		}

		select {
		case client.send <- data:
		default:
			// Клиент не готов, пропускаем
		}
		return true
	})
}

func (ws *WSManager) HandleWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	login, _ := ctx.Value("login").(string)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		userID:        userID,
		conn:          conn,
		send:          make(chan []byte, 256),
		subChats:      make(map[string]bool),
		subUsers:      make(map[string]bool),
		lastHeartbeat: time.Now(),
	}

	ws.register <- client

	// Отправляем информацию о подключении
	connInfo := models.WSConnectionInfo{
		UserID:          userID,
		SubscribedChats: []string{},
		SubscribedUsers: []string{},
		ConnectedAt:     time.Now().UnixMilli(),
	}

	initMsg := WSMessage{
		Type:    "connection_info",
		Payload: connInfo,
	}

	data, _ := msgpack.Encode(&initMsg)
	conn.WriteMessage(websocket.BinaryMessage, data)

	// Запускаем goroutines для чтения и записи
	go client.writePump(ws)
	go client.readPump(ws)

	// Обновляем статус пользователя на online
	ws.BroadcastUserStatus(userID, true)
}

func (c *Client) readPump(ws *WSManager) {
	defer func() {
		ws.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.mu.Lock()
		c.lastHeartbeat = time.Now()
		c.mu.Unlock()
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Логируем ошибку
			}
			break
		}

		// Обрабатываем входящее сообщение
		var wsMsg WSMessage
		if err := msgpack.Decode(message, &wsMsg); err != nil {
			// Пробуем JSON как fallback
			json.Unmarshal(message, &wsMsg)
		}

		c.handleIncomingMessage(&wsMsg, ws)
	}
}

func (c *Client) writePump(ws *WSManager) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.mu.Lock()
			if time.Since(c.lastHeartbeat) > 90*time.Second {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleIncomingMessage(msg *WSMessage, ws *WSManager) {
	switch msg.Type {
	case "subscribe_chats":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if chatIDs, ok := payload["chat_ids"].([]interface{}); ok {
				c.mu.Lock()
				for _, id := range chatIDs {
					if strID, ok := id.(string); ok {
						c.subChats[strID] = true
					}
				}
				c.mu.Unlock()
			}
		}

	case "subscribe_users":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if userIDs, ok := payload["user_ids"].([]interface{}); ok {
				c.mu.Lock()
				for _, id := range userIDs {
					if strID, ok := id.(string); ok {
						c.subUsers[strID] = true
					}
				}
				c.mu.Unlock()
			}
		}

	case "unsubscribe_all":
		c.mu.Lock()
		c.subChats = make(map[string]bool)
		c.subUsers = make(map[string]bool)
		c.mu.Unlock()

	case "heartbeat":
		c.mu.Lock()
		c.lastHeartbeat = time.Now()
		c.mu.Unlock()
	}
}

func (ws *WSManager) BroadcastNewMessage(msg models.WSNewMessage) {
	ws.broadcast <- &WSMessage{
		Type:    "new_message",
		Payload: msg,
	}
}

func (ws *WSManager) BroadcastUserStatus(userID string, isOnline bool) {
	lastSeen := time.Now().UnixMilli()
	if !isOnline {
		lastSeen = time.Now().Add(-5 * time.Minute).UnixMilli()
	}

	ws.broadcast <- &WSMessage{
		Type: "user_status",
		Payload: models.WSUserStatus{
			UserID:   userID,
			IsOnline: isOnline,
			LastSeen: lastSeen,
		},
	}
}

func (ws *WSManager) SendToUser(userID string, msg *WSMessage) {
	if client, ok := ws.clients.Load(userID); ok {
		c := client.(*Client)
		data, err := msgpack.Encode(msg)
		if err != nil {
			return
		}
		select {
		case c.send <- data:
		default:
			// Клиент не готов
		}
	}
}

func (ws *WSManager) IsUserOnline(userID string) bool {
	_, exists := ws.clients.Load(userID)
	return exists
}

func (ws *WSManager) GetOnlineUsers() []string {
	var users []string
	ws.clients.Range(func(key, value interface{}) bool {
		users = append(users, key.(string))
		return true
	})
	return users
}
