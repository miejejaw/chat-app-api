package models

import "time"

type User struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Username        string    `gorm:"unique; not null" json:"username"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `gorm:"unique;not null" json:"email"`
	ProfileImageUrl string    `json:"profile_image"`
	Password        string    `gorm:"not null" json:"password"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
