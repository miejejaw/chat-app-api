package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"chat-app-api/internal/models"
	"chat-app-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket Upgrader configuration
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Mutex for concurrent map access
var mu sync.Mutex

// Map to track connected clients by user ID
var clients = make(map[uint]*websocket.Conn)

// MessageHandler handles real-time chat through WebSocket
type MessageHandler struct {
	messageService services.MessageService
}

// NewMessageHandler creates a new instance of MessageHandler
func NewMessageHandler(messageService services.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// HandleConnections handles incoming WebSocket connections
func (h *MessageHandler) HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(ws)

	// Extract the user ID from query parameters
	senderID, err := getUserIDFromQuery(c)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return
	}

	// Add the WebSocket connection to the map
	mu.Lock()
	clients[senderID] = ws
	mu.Unlock()

	// Listen for incoming messages
	for {
		// Define a structure for receiving only content, receiver_id, and sender_id
		var incomingMessage struct {
			Content    string `json:"content"`
			ReceiverID uint   `json:"receiver_id"`
		}

		// Read the message from the WebSocket connection
		err := ws.ReadJSON(&incomingMessage)
		if err != nil {
			log.Printf("Error reading JSON message: %v", err)
			mu.Lock()
			delete(clients, senderID)
			mu.Unlock()
			break
		}

		// Create a Message model with the extracted information
		message := models.Message{
			Content:    incomingMessage.Content,
			SenderID:   senderID, // Use the sender's ID from the WebSocket connection query
			ReceiverID: incomingMessage.ReceiverID,
		}

		// Process and save the message
		h.processMessage(&message)
	}
}

// processMessage saves the message using the service and forwards it to the recipient
func (h *MessageHandler) processMessage(msg *models.Message) {
	// Save the message to the database using the message service
	createdMsg, err := h.messageService.CreateMessage(msg)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		return
	}

	// Send the message to the recipient if they are connected
	h.sendMessageToRecipient(createdMsg, msg.ReceiverID)
	createdMsg.IsSelf = true
	h.sendMessageToSender(createdMsg, msg.Sender.ID)
}

func (h *MessageHandler) sendMessageToRecipient(msg *services.RealTimeMessageResponse, recipientID uint) {
	mu.Lock()
	recipientWS, exists := clients[recipientID]
	mu.Unlock()

	if exists {
		// Send the message to the recipient
		err := recipientWS.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending message to recipient %d: %v", recipientID, err)
		}
	} else {
		log.Printf("Recipient %d is not online", recipientID)
	}
}
func (h *MessageHandler) sendMessageToSender(msg *services.RealTimeMessageResponse, senderID uint) {
	mu.Lock()
	senderWS, exists := clients[senderID]
	mu.Unlock()

	if exists {
		// Send the message to the recipient
		err := senderWS.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending message to sender %d: %v", senderID, err)
		}
	} else {
		log.Printf("Sender %d is not online", senderID)
	}
}

func (h *MessageHandler) GetFriendsWithLastMessage(c *gin.Context) {
	currentUserIDString, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	currentUserID, err := strconv.ParseUint(currentUserIDString.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	messages, err := h.messageService.GetFriendListWithLastMessage(uint(currentUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) GetMessagesBySenderIdAndReceiverId(c *gin.Context) {
	currentUserIDString, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	currentUserID, err := strconv.ParseUint(currentUserIDString.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	friendID, err := getUserIDFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
		return
	}

	messages, err := h.messageService.GetMessagesBySenderIdAndReceiverId(uint(currentUserID), friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// Helper function to extract user ID from the query parameters
func getUserIDFromQuery(c *gin.Context) (uint, error) {
	userIDStr := c.Query("user_id")
	var userID uint
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	return userID, err
}
