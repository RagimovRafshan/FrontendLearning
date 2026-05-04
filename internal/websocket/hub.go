package websocket

import (
	"context"
	"encoding/binary"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"messenger/internal/model"
	"messenger/pkg/codec"
)

const (
	// Максимальный размер сообщения (1MB)
	MaxMessageSize = 1024 * 1024
	
	// Таймауты
	WriteTimeout   = 10 * time.Second
	ReadTimeout    = 60 * time.Second
	PingInterval   = 30 * time.Second
	PongWait       = PingInterval + 10*time.Second
	
	// Буферы
	SendBufferSize = 256
)

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidMessage   = errors.New("invalid message format")
	ErrConnectionClosed = errors.New("connection closed")
)

// upgrader для WebSocket соединений
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Разрешаем все origins (в продакшене нужно ограничить)
		return true
	},
}

// Client представляет подключенного клиента
type Client struct {
	ID        string
	UserID    string
	Username  string
	conn      *websocket.Conn
	send      chan []byte
	hub       *Hub
	mu        sync.RWMutex
	isAuth    bool
	lastSeen  time.Time
	ctx       context.Context
	cancel    context.CancelFunc
}

// Hub управляет всеми подключениями
type Hub struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	
	// Индексы для быстрого поиска
	userClients map[string]map[string]*Client // user_id -> client_id -> client
}

// MessageHandler обрабатывает входящие сообщения
type MessageHandler func(ctx context.Context, client *Client, req *model.Request) (*model.Response, error)

// NewHub создает новый хаб
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		broadcast:   make(chan []byte, SendBufferSize),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		userClients: make(map[string]map[string]*Client),
	}
}

// Run запускает хаб
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			if client.UserID != "" {
				if _, ok := h.userClients[client.UserID]; !ok {
					h.userClients[client.UserID] = make(map[string]*Client)
				}
				h.userClients[client.UserID][client.ID] = client
			}
			h.mu.Unlock()
			
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				if client.UserID != "" {
					if clients, ok := h.userClients[client.UserID]; ok {
						delete(clients, client.ID)
						if len(clients) == 0 {
							delete(h.userClients, client.UserID)
						}
					}
				}
			}
			h.mu.Unlock()
			close(client.send)
			
		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Буфер переполнен, закрываем соединение
					go client.disconnect()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// ServeWS обрабатывает WebSocket подключение
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, handler MessageHandler) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}

	client := &Client{
		ID:       generateClientID(),
		conn:     conn,
		send:     make(chan []byte, SendBufferSize),
		hub:      h,
		lastSeen: time.Now(),
	}
	
	client.ctx, client.cancel = context.WithCancel(context.Background())
	
	h.register <- client
	
	// Запускаем горутину для записи
	go client.writePump()
	// Запускаем горутину для чтения с обработкой
	go client.readPump(handler)
}

// GetOnlineUsers возвращает список онлайн пользователей
func (h *Hub) GetOnlineUsers(userIDs []string) map[string]bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	result := make(map[string]bool)
	for _, uid := range userIDs {
		if clients, ok := h.userClients[uid]; ok {
			result[uid] = len(clients) > 0
		} else {
			result[uid] = false
		}
	}
	return result
}

// GetUserClients возвращает всех клиентов пользователя
func (h *Hub) GetUserClients(userID string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var clients []*Client
	if userClients, ok := h.userClients[userID]; ok {
		for _, c := range userClients {
			clients = append(clients, c)
		}
	}
	return clients
}

// BroadcastToUser отправляет сообщение всем клиентам пользователя
func (h *Hub) BroadcastToUser(userID string, response *model.Response) error {
	data, err := codec.Encode(response)
	if err != nil {
		return err
	}
	
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if clients, ok := h.userClients[userID]; ok {
		for _, client := range clients {
			select {
			case client.send <- data:
			default:
				go client.disconnect()
			}
		}
	}
	return nil
}

// BroadcastToChat отправляет сообщение всем участникам чата
func (h *Hub) BroadcastToChat(chatID string, memberIDs []string, response *model.Response) error {
	data, err := codec.Encode(response)
	if err != nil {
		return err
	}
	
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	for _, uid := range memberIDs {
		if clients, ok := h.userClients[uid]; ok {
			for _, client := range clients {
				select {
				case client.send <- data:
				default:
					go client.disconnect()
				}
			}
		}
	}
	return nil
}

// IsUserOnline проверяет, онлайн ли пользователь
func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	clients, ok := h.userClients[userID]
	return ok && len(clients) > 0
}

// Client методы

func (c *Client) disconnect() {
	c.cancel()
	c.hub.unregister <- c
	c.conn.Close()
}

func (c *Client) readPump(handler MessageHandler) {
	defer func() {
		c.disconnect()
	}()
	
	c.conn.SetReadLimit(MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.lastSeen = time.Now()
		c.conn.SetReadDeadline(time.Now().Add(ReadTimeout))
		return nil
	})
	
	for {
		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Логируем ошибку
			}
			break
		}
		
		if messageType != websocket.BinaryMessage {
			// Ожидаем только бинарные сообщения (MessagePack)
			continue
		}
		
		// Декодируем длину сообщения (первые 4 байта)
		if len(data) < 4 {
			continue
		}
		
		msgLen := binary.BigEndian.Uint32(data[:4])
		if len(data) < 4+int(msgLen) {
			continue
		}
		
		msgData := data[4 : 4+msgLen]
		
		var req model.Request
		if err := codec.Decode(msgData, &req); err != nil {
			// Отправляем ошибку клиенту
			resp := model.NewErrorResponse("", 400, "Invalid message format")
			respData, _ := codec.Encode(resp)
			c.send <- respData
			continue
		}
		
		// Обрабатываем запрос
		ctx := context.WithValue(c.ctx, "client", c)
		response, err := handler(ctx, c, &req)
		if err != nil {
			response = model.NewErrorResponse(req.RequestID, 500, err.Error())
		}
		
		// Отправляем ответ
		if response != nil {
			respData, err := codec.Encode(response)
			if err == nil {
				c.send <- respData
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(PingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(WriteTimeout))
			if !ok {
				// Канал закрыт
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			// Формируем пакет: [4 байта длина][данные]
			packet := make([]byte, 4+len(message))
			binary.BigEndian.PutUint32(packet[:4], uint32(len(message)))
			copy(packet[4:], message)
			
			if err := c.conn.WriteMessage(websocket.BinaryMessage, packet); err != nil {
				return
			}
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(WriteTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SetAuth устанавливает аутентификацию для клиента
func (c *Client) SetAuth(userID, username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.UserID = userID
	c.Username = username
	c.isAuth = true
}

// IsAuthenticated проверяет аутентификацию
func (c *Client) IsAuthenticated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isAuth
}

// GenerateClientID генерирует уникальный ID клиента
func generateClientID() string {
	return time.Now().Format("20060102150405.000000000")
}
