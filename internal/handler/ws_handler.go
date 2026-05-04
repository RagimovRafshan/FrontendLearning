package handler

import (
	"context"

	"messenger/internal/model"
	"messenger/internal/service"
	"messenger/internal/websocket"
)

// WSHandler обрабатывает все WebSocket сообщения
type WSHandler struct {
	authService  *service.AuthService
	chatService  *service.ChatService
	hub          *websocket.Hub
}

// NewWSHandler создает обработчик WebSocket
func NewWSHandler(authSvc *service.AuthService, chatSvc *service.ChatService, hub *websocket.Hub) *WSHandler {
	return &WSHandler{
		authService: authSvc,
		chatService: chatSvc,
		hub:         hub,
	}
}

// HandleMessage обрабатывает входящее сообщение
func (h *WSHandler) HandleMessage(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	switch req.Type {
	// Авторизация
	case model.MsgTypeRegister:
		return h.handleRegister(ctx, client, req)
	case model.MsgTypeLogin:
		return h.handleLogin(ctx, client, req)
	case model.MsgTypeForgotPassword:
		return h.handleForgotPassword(ctx, client, req)
	case model.MsgTypeRefreshToken:
		return h.handleRefreshToken(ctx, client, req)

	// Пользователи
	case model.MsgTypeSearchUser:
		return h.handleSearchUser(ctx, client, req)
	case model.MsgTypeUpdateProfile:
		return h.handleUpdateProfile(ctx, client, req)
	case model.MsgTypeGetProfile:
		return h.handleGetProfile(ctx, client, req)

	// Чаты
	case model.MsgTypeCreateChatByLink:
		return h.handleCreateChatByLink(ctx, client, req)
	case model.MsgTypeCreateChatWithUser:
		return h.handleCreateChatWithUser(ctx, client, req)
	case model.MsgTypeCreateGroupChat:
		return h.handleCreateGroupChat(ctx, client, req)
	case model.MsgTypeCreateChannel:
		return h.handleCreateChannel(ctx, client, req)
	case model.MsgTypeEditGroupChat:
		return h.handleEditGroupChat(ctx, client, req)
	case model.MsgTypeEditChannel:
		return h.handleEditChannel(ctx, client, req)
	case model.MsgTypeGetChats:
		return h.handleGetChats(ctx, client, req)
	case model.MsgTypeGetChatInfo:
		return h.handleGetChatInfo(ctx, client, req)

	// Сообщения
	case model.MsgTypeSendMessage:
		return h.handleSendMessage(ctx, client, req)
	case model.MsgTypeEditMessage:
		return h.handleEditMessage(ctx, client, req)
	case model.MsgTypeGetMessages:
		return h.handleGetMessages(ctx, client, req)
	case model.MsgTypeDeleteMessage:
		return h.handleDeleteMessage(ctx, client, req)

	// Статусы
	case model.MsgTypeGetOnlineStatus:
		return h.handleGetOnlineStatus(ctx, client, req)
	case model.MsgTypeSubscribeEvents:
		return h.handleSubscribeEvents(ctx, client, req)
	case model.MsgTypeTypingStatus:
		return h.handleTypingStatus(ctx, client, req)

	default:
		return model.NewErrorResponse(req.RequestID, 400, "unknown message type"), nil
	}
}

// ========== Авторизация ==========

func (h *WSHandler) handleRegister(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	payload, ok := req.Payload.(*model.RegisterRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	resp, err := h.authService.Register(ctx, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	client.SetAuth(resp.User.ID, resp.User.Username)

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, resp, nil), nil
}

func (h *WSHandler) handleLogin(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	payload, ok := req.Payload.(*model.LoginRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	resp, err := h.authService.Login(ctx, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	client.SetAuth(resp.User.ID, resp.User.Username)

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, resp, nil), nil
}

func (h *WSHandler) handleForgotPassword(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	payload, ok := req.Payload.(*model.ForgotPasswordRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	err := h.authService.ForgotPassword(ctx, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, map[string]string{"status": "email_sent"}, nil), nil
}

func (h *WSHandler) handleRefreshToken(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	payload, ok := req.Payload.(*model.RefreshTokenRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	resp, err := h.authService.RefreshToken(ctx, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, resp, nil), nil
}

// ========== Пользователи ==========

func (h *WSHandler) handleSearchUser(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.SearchUserRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	users := h.authService.SearchUsers(payload.Query, payload.Limit)

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, users, nil), nil
}

func (h *WSHandler) handleUpdateProfile(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.UpdateProfileRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	user, err := h.authService.UpdateProfile(client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, user, nil), nil
}

func (h *WSHandler) handleGetProfile(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.GetProfileRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	userID := payload.UserID
	if userID == "" {
		userID = client.UserID
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	user.IsOnline = h.hub.IsUserOnline(userID)

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, user, nil), nil
}

// ========== Чаты ==========

func (h *WSHandler) handleCreateChatByLink(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.CreateChatByLinkRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chat, err := h.chatService.CreateChatByLink(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chat, nil), nil
}

func (h *WSHandler) handleCreateChatWithUser(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.CreateChatWithUserRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chat, err := h.chatService.CreateChatWithUser(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chat, nil), nil
}

func (h *WSHandler) handleCreateGroupChat(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.CreateGroupChatRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chat, err := h.chatService.CreateGroupChat(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chat, nil), nil
}

func (h *WSHandler) handleCreateChannel(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.CreateChannelRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chat, err := h.chatService.CreateChannel(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chat, nil), nil
}

func (h *WSHandler) handleEditGroupChat(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.EditGroupChatRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chat, err := h.chatService.EditGroupChat(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	h.notifyChatMembers(chat.ID, chat.Members, &model.Response{
		Type:      model.MsgTypeResponse,
		RequestID: req.RequestID,
		Success:   true,
		Data:      chat,
	})

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chat, nil), nil
}

func (h *WSHandler) handleEditChannel(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.EditChannelRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chat, err := h.chatService.EditChannel(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chat, nil), nil
}

func (h *WSHandler) handleGetChats(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.GetChatsRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chats, err := h.chatService.GetChats(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chats, nil), nil
}

func (h *WSHandler) handleGetChatInfo(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.GetChatInfoRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	chat, err := h.chatService.GetChatInfo(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, chat, nil), nil
}

// ========== Сообщения ==========

func (h *WSHandler) handleSendMessage(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.SendMessageRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	msg, err := h.chatService.SendMessage(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	chatInfo, _ := h.chatService.GetChatInfo(ctx, client.UserID, &model.GetChatInfoRequest{ChatID: payload.ChatID})

	if chatInfo != nil {
		sender, _ := h.authService.GetUserByID(client.UserID)
		msg.Sender = sender

		response := &model.Response{
			Type:      model.MsgTypeNewMessage,
			RequestID: req.RequestID,
			Success:   true,
			Data:      msg,
		}

		h.hub.BroadcastToChat(payload.ChatID, chatInfo.Members, response)
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, msg, nil), nil
}

func (h *WSHandler) handleEditMessage(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.EditMessageRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	msg, err := h.chatService.EditMessage(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, msg, nil), nil
}

func (h *WSHandler) handleGetMessages(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.GetMessagesRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	messages, err := h.chatService.GetMessages(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, messages, nil), nil
}

func (h *WSHandler) handleDeleteMessage(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.DeleteMessageRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	err := h.chatService.DeleteMessage(ctx, client.UserID, payload)
	if err != nil {
		return h.errorResponse(req.RequestID, err), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true, map[string]string{"status": "deleted"}, nil), nil
}

// ========== Статусы ==========

func (h *WSHandler) handleGetOnlineStatus(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.GetOnlineStatusRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	onlineStatus := h.hub.GetOnlineUsers(payload.UserIDs)

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true,
		model.OnlineStatusResponse{OnlineUsers: onlineStatus}, nil), nil
}

func (h *WSHandler) handleSubscribeEvents(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	payload, ok := req.Payload.(*model.SubscribeEventsRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true,
		map[string]interface{}{"subscribed": payload.Events}, nil), nil
}

func (h *WSHandler) handleTypingStatus(ctx context.Context, client *websocket.Client, req *model.Request) (*model.Response, error) {
	if !client.IsAuthenticated() {
		return model.NewErrorResponse(req.RequestID, 401, "unauthorized"), nil
	}

	_, ok := req.Payload.(*model.TypingStatusRequest)
	if !ok {
		return model.NewErrorResponse(req.RequestID, 400, "invalid payload"), nil
	}

	return model.NewResponse(req.RequestID, model.MsgTypeResponse, true,
		map[string]interface{}{"status": "ok"}, nil), nil
}

// ========== Вспомогательные методы ==========

func (h *WSHandler) errorResponse(reqID string, err error) *model.Response {
	code := 500
	message := err.Error()

	if err == service.ErrUserNotFound || err == service.ErrChatNotFound {
		code = 404
	} else if err == service.ErrInvalidPassword {
		code = 401
	} else if err == service.ErrUserExists || err == service.ErrAccessDenied {
		code = 403
	}

	return model.NewErrorResponse(reqID, code, message)
}

func (h *WSHandler) notifyChatMembers(chatID string, memberIDs []string, response *model.Response) {
	h.hub.BroadcastToChat(chatID, memberIDs, response)
}
