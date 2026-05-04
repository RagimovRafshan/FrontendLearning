package handlers

import (
"fmt"
"net/http"
"sync"
"time"

"github.com/google/uuid"
"messenger/internal/models"
"messenger/internal/storage"
"messenger/pkg/msgpack"
)

type MessageHandler struct {
scylla       *storage.ScyllaDBClient
messages     sync.Map // Map[messageID]messageData
chatMessages sync.Map // Map[chatID][]messageID
}

type messageData struct {
MessageID   string
ChatID      string
SenderID    string
SenderName  string
Content     string
MessageType string // text, image, file, etc.
MediaURL    string
CreatedAt   int64
EditedAt    int64
}

func NewMessageHandler(scylla *storage.ScyllaDBClient) *MessageHandler {
return &MessageHandler{
scylla: scylla,
}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
ctx := r.Context()
userID, _ := ctx.Value("user_id").(string)
login, _ := ctx.Value("login").(string)

if userID == "" {
msgpack.Respond(w, http.StatusUnauthorized, models.Response{
Success: false,
Error:   "Unauthorized",
})
return
}

var req models.SendMessageRequest
if err := msgpack.DecodeRequest(r, &req); err != nil {
msgpack.Respond(w, http.StatusBadRequest, models.Response{
Success: false,
Error:   "Invalid request body: " + err.Error(),
})
return
}

if req.ChatID == "" || req.Content == "" {
msgpack.Respond(w, http.StatusBadRequest, models.Response{
Success: false,
Error:   "chat_id and content are required",
})
return
}

messageType := req.MessageType
if messageType == "" {
messageType = "text"
}

messageID := uuid.New().String()

message := &messageData{
MessageID:   messageID,
ChatID:      req.ChatID,
SenderID:    userID,
SenderName:  login,
Content:     req.Content,
MessageType: messageType,
MediaURL:    req.MediaURL,
CreatedAt:   time.Now().UnixMilli(),
}

h.messages.Store(messageID, message)

// Add to chat messages list
var chatMsgs []string
if existing, ok := h.chatMessages.Load(req.ChatID); ok {
chatMsgs = existing.([]string)
}
chatMsgs = append(chatMsgs, messageID)
h.chatMessages.Store(req.ChatID, chatMsgs)

msgpack.Respond(w, http.StatusOK, models.Response{
Success: true,
Data: models.MessageResponse{
MessageID:   message.MessageID,
ChatID:      message.ChatID,
SenderID:    message.SenderID,
SenderName:  message.SenderName,
Content:     message.Content,
MessageType: message.MessageType,
MediaURL:    message.MediaURL,
CreatedAt:   message.CreatedAt,
},
})
}

func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
ctx := r.Context()
userID, _ := ctx.Value("user_id").(string)
if userID == "" {
msgpack.Respond(w, http.StatusUnauthorized, models.Response{
Success: false,
Error:   "Unauthorized",
})
return
}

var req models.GetMessageRequest
if err := msgpack.DecodeRequest(r, &req); err != nil {
msgpack.Respond(w, http.StatusBadRequest, models.Response{
Success: false,
Error:   "Invalid request body: " + err.Error(),
})
return
}

if req.MessageID == "" {
msgpack.Respond(w, http.StatusBadRequest, models.Response{
Success: false,
Error:   "message_id is required",
})
return
}

messageInterface, exists := h.messages.Load(req.MessageID)
if !exists {
msgpack.Respond(w, http.StatusNotFound, models.Response{
Success: false,
Error:   "Message not found",
})
return
}

message := messageInterface.(*messageData)

msgpack.Respond(w, http.StatusOK, models.Response{
Success: true,
Data: models.MessageResponse{
MessageID:   message.MessageID,
ChatID:      message.ChatID,
SenderID:    message.SenderID,
SenderName:  message.SenderName,
Content:     message.Content,
MessageType: message.MessageType,
MediaURL:    message.MediaURL,
CreatedAt:   message.CreatedAt,
EditedAt:    message.EditedAt,
},
})
}

func (h *MessageHandler) EditMessage(w http.ResponseWriter, r *http.Request) {
ctx := r.Context()
userID, _ := ctx.Value("user_id").(string)
if userID == "" {
msgpack.Respond(w, http.StatusUnauthorized, models.Response{
Success: false,
Error:   "Unauthorized",
})
return
}

var req models.EditMessageRequest
if err := msgpack.DecodeRequest(r, &req); err != nil {
msgpack.Respond(w, http.StatusBadRequest, models.Response{
Success: false,
Error:   "Invalid request body: " + err.Error(),
})
return
}

if req.MessageID == "" || req.Content == "" {
msgpack.Respond(w, http.StatusBadRequest, models.Response{
Success: false,
Error:   "message_id and content are required",
})
return
}

messageInterface, exists := h.messages.Load(req.MessageID)
if !exists {
msgpack.Respond(w, http.StatusNotFound, models.Response{
Success: false,
Error:   "Message not found",
})
return
}

message := messageInterface.(*messageData)

if message.SenderID != userID {
msgpack.Respond(w, http.StatusForbidden, models.Response{
Success: false,
Error:   "Only message sender can edit",
})
return
}

message.Content = req.Content
message.EditedAt = time.Now().UnixMilli()

h.messages.Store(req.MessageID, message)

msgpack.Respond(w, http.StatusOK, models.Response{
Success: true,
Data: models.MessageResponse{
MessageID:   message.MessageID,
ChatID:      message.ChatID,
SenderID:    message.SenderID,
SenderName:  message.SenderName,
Content:     message.Content,
MessageType: message.MessageType,
MediaURL:    message.MediaURL,
CreatedAt:   message.CreatedAt,
EditedAt:    message.EditedAt,
},
})
}

func (h *MessageHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
ctx := r.Context()
userID, _ := ctx.Value("user_id").(string)
if userID == "" {
msgpack.Respond(w, http.StatusUnauthorized, models.Response{
Success: false,
Error:   "Unauthorized",
})
return
}

chatID := r.URL.Query().Get("chat_id")
if chatID == "" {
msgpack.Respond(w, http.StatusBadRequest, models.Response{
Success: false,
Error:   "chat_id query parameter is required",
})
return
}

limit := 50
offset := 0

if l := r.URL.Query().Get("limit"); l != "" {
if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
limit = 50
}
}
if o := r.URL.Query().Get("offset"); o != "" {
if _, err := fmt.Sscanf(o, "%d", &offset); err != nil {
offset = 0
}
}

if limit > 100 {
limit = 100
}

messageIDsInterface, exists := h.chatMessages.Load(chatID)
if !exists {
msgpack.Respond(w, http.StatusOK, models.Response{
Success: true,
Data:    []models.MessageResponse{},
})
return
}

messageIDs := messageIDsInterface.([]string)

start := offset
if start > len(messageIDs) {
start = len(messageIDs)
}
end := start + limit
if end > len(messageIDs) {
end = len(messageIDs)
}

var messages []models.MessageResponse
for i := start; i < end; i++ {
msgID := messageIDs[i]
if msgInterface, ok := h.messages.Load(msgID); ok {
msg := msgInterface.(*messageData)
messages = append(messages, models.MessageResponse{
MessageID:   msg.MessageID,
ChatID:      msg.ChatID,
SenderID:    msg.SenderID,
SenderName:  msg.SenderName,
Content:     msg.Content,
MessageType: msg.MessageType,
MediaURL:    msg.MediaURL,
CreatedAt:   msg.CreatedAt,
EditedAt:    msg.EditedAt,
})
}
}

for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
messages[i], messages[j] = messages[j], messages[i]
}

msgpack.Respond(w, http.StatusOK, models.Response{
Success: true,
Data:    messages,
})
}
