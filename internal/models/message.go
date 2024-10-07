package models

import (
	"time"
)

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Content    string    `gorm:"not null" json:"content"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	Sender     User      `gorm:"foreignKey:SenderID" json:"sender"` // Define sender relationship
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Receiver   User      `gorm:"foreignKey:ReceiverID" json:"receiver"` // Define receiver relationship
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
