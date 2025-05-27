package http

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/artemxgod/forum/internal/forum/domain"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Handler handles HTTP requests
type Handler struct {
	service      domain.ForumService
	logger       *zap.Logger
	upgrader     websocket.Upgrader
	clients      map[*websocket.Conn]int // websocket -> userID
	clientsMutex sync.RWMutex
}

// NewHandler creates a new HTTP handler
func NewHandler(service domain.ForumService, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, this should be more restrictive
			},
		},
		clients: make(map[*websocket.Conn]int),
	}
}

// @Summary Get chat messages
// @Description Get recent chat messages
// @Tags chat
// @Accept json
// @Produce json
// @Param limit query int false "Number of messages to return"
// @Success 200 {array} domain.Message
// @Router /api/chat/messages [get]
func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	limit := 50 // default limit
	messages, err := h.service.GetMessages(limit)
	if err != nil {
		h.logger.Error("failed to get messages", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

// @Summary Connect to chat WebSocket
// @Description Connect to chat WebSocket for real-time messages
// @Tags chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 101 {string} string "Switching Protocols"
// @Router /ws/chat [get]
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	// Upgrade connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("failed to upgrade connection", zap.Error(err))
		return
	}

	// Register client
	h.clientsMutex.Lock()
	h.clients[conn] = userID
	h.clientsMutex.Unlock()

	// Clean up on disconnect
	defer func() {
		h.clientsMutex.Lock()
		delete(h.clients, conn)
		h.clientsMutex.Unlock()
		conn.Close()
	}()

	// Handle incoming messages
	for {
		var msg domain.WebsocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("websocket error", zap.Error(err))
			}
			break
		}

		switch msg.Type {
		case "message":
			content, ok := msg.Payload["content"].(string)
			if !ok {
				continue
			}

			err = h.service.SendMessage(userID, content)
			if err != nil {
				h.logger.Error("failed to save message", zap.Error(err))
				continue
			}

			// Broadcast message to all clients
			h.broadcastMessage(msg)
		}
	}
}

func (h *Handler) broadcastMessage(msg domain.WebsocketMessage) {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()

	for conn := range h.clients {
		err := conn.WriteJSON(msg)
		if err != nil {
			h.logger.Error("failed to send message", zap.Error(err))
			conn.Close()
		}
	}
}

// RegisterRoutes registers HTTP routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/chat/messages", h.GetMessages)
	mux.HandleFunc("/ws/chat", h.HandleWebSocket)
}
