package services

import (
	"chat-app-api/internal/models"
	"chat-app-api/internal/repositories"
)

type messageResponse struct {
	ID      uint   `json:"id"`
	IsSelf  bool   `json:"is_self"`
	IsRead  bool   `json:"is_read"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type RealTimeMessageResponse struct {
	ID         uint   `json:"id"`
	IsSelf     bool   `json:"is_self"`
	IsRead     bool   `json:"is_read"`
	Message    string `json:"message"`
	Time       string `json:"time"`
	ReceiverID uint   `json:"receiver_id"`
	Sender     struct {
		ID              uint   `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	} `json:"sender"`
}

type MessageService interface {
	CreateMessage(message *models.Message) (*RealTimeMessageResponse, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
	DeleteMessage(id uint) error
	GetMessagesBySenderIdAndReceiverId(senderId, receiverId uint) ([]messageResponse, error)
	GetFriendListWithLastMessage(userID uint) ([]repositories.FriendsList, error)
}

type messageService struct {
	messageRepository repositories.MessageRepository
}

func NewMessageService(repo repositories.MessageRepository) MessageService {
	return &messageService{messageRepository: repo}
}

func (s *messageService) CreateMessage(message *models.Message) (*RealTimeMessageResponse, error) {
	response, _ := s.messageRepository.CreateMessage(message)

	realTimeResponse := &RealTimeMessageResponse{
		ID:         response.ID,
		IsSelf:     false,
		IsRead:     false,
		Message:    response.Content,
		Time:       response.CreatedAt.Format("2006-01-02 15:04"),
		ReceiverID: response.ReceiverID,
		Sender: struct {
			ID              uint   `json:"id"`
			FirstName       string `json:"first_name"`
			LastName        string `json:"last_name"`
			ProfileImageUrl string `json:"profile_image_url"`
		}{
			ID:              response.Sender.ID,
			FirstName:       response.Sender.FirstName,
			LastName:        response.Sender.LastName,
			ProfileImageUrl: response.Sender.ProfileImageUrl,
		},
	}

	return realTimeResponse, nil
}

func (s *messageService) UpdateMessage(message *models.Message) (*models.Message, error) {
	return s.messageRepository.UpdateMessage(message)
}

func (s *messageService) DeleteMessage(id uint) error {
	return s.messageRepository.DeleteMessage(id)
}

func (s *messageService) GetMessagesBySenderIdAndReceiverId(currentID uint, friendID uint) ([]messageResponse, error) {
	messages, err := s.messageRepository.FindBySenderIdAndReceiverId(currentID, friendID)
	if err != nil {
		return nil, err
	}

	var response []messageResponse
	for _, message := range messages {
		response = append(response, messageResponse{
			ID:      message.ID,
			IsSelf:  message.SenderID == currentID,
			Message: message.Content,
			IsRead:  true,
			Time:    message.CreatedAt.Format("2006-01-02 15:04:05"),
		})

	}
	return response, nil
}

func (s *messageService) GetFriendListWithLastMessage(userID uint) ([]repositories.FriendsList, error) {
	return s.messageRepository.GetFriendListWithLastMessage(userID)
}
