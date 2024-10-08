package repositories

import (
	"chat-app-api/internal/models"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type FriendsList struct {
	Profile struct {
		ID              uint   `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
		Username        string `json:"username"`
	} `json:"profile"`
	LastSeen    string `json:"last_seen"`
	UnreadCount int    `json:"unread_count"`
	LastMessage struct {
		Content string `json:"content"`
		Time    string `json:"time"`
	} `json:"last_message"`
}

type MessageRepository interface {
	CreateMessage(message *models.Message) (*models.Message, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
	DeleteMessage(id uint) error
	FindBySenderIdAndReceiverId(senderID uint, receiverID uint) ([]models.Message, error)
	GetFriendListWithLastMessage(userID uint) ([]FriendsList, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(message *models.Message) (*models.Message, error) {
	if err := r.db.Create(message).Error; err != nil {
		return nil, err
	}

	// Preload the sender information
	if err := r.db.Preload("Sender").First(message, message.ID).Error; err != nil {
		return nil, err
	}

	return message, nil
}

func (r *messageRepository) UpdateMessage(message *models.Message) (*models.Message, error) {
	if err := r.db.Save(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *messageRepository) DeleteMessage(id uint) error {
	if err := r.db.Delete(&models.Message{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) FindBySenderIdAndReceiverId(currentID uint, friendID uint) ([]models.Message, error) {
	var messages []models.Message

	// Query to fetch all messages between current user and friend, regardless of who sent it
	if err := r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		currentID, friendID, friendID, currentID).
		Order("created_at desc").
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) GetFriendListWithLastMessage(userID uint) ([]FriendsList, error) {
	var friendsList []FriendsList

	// Query to fetch the friend and the last message details
	rows, err := r.db.Raw(`
		SELECT 
			u.id, u.first_name, u.last_name, u.profile_image_url, u.username,
			m1.content AS last_message_content, 
			m1.created_at AS last_message_time
		FROM messages m1
		INNER JOIN (
			SELECT MAX(id) AS id
			FROM messages
			WHERE sender_id = ? OR receiver_id = ?
			GROUP BY LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id)
		) m2 ON m1.id = m2.id
		INNER JOIN users u ON (u.id = m1.sender_id OR u.id = m1.receiver_id)
		WHERE u.id != ?
	`, userID, userID, userID).Rows()

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Process rows and map to FriendsList struct
	for rows.Next() {
		var friend FriendsList
		var lastMessageTime time.Time

		// Scan the results into the FriendsList structure
		if err := rows.Scan(
			&friend.Profile.ID,
			&friend.Profile.FirstName,
			&friend.Profile.LastName,
			&friend.Profile.ProfileImageUrl,
			&friend.Profile.Username,
			&friend.LastMessage.Content,
			&lastMessageTime,
		); err != nil {
			return nil, err
		}

		// Format the last message time as a string (if you need a specific format)
		friend.LastMessage.Time = lastMessageTime.Format("2006-01-02 15:04")

		// sample data
		friend.UnreadCount = 0
		friend.LastSeen = "2021-09-01 15:04"

		// Append the result to the friends list
		friendsList = append(friendsList, friend)
	}

	return friendsList, nil
}
