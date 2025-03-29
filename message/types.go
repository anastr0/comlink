package message

import (
	"time"

	"gorm.io/gorm"
)

type MessageHandler struct {
	db *gorm.DB
}

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name" binding:"required"`
}

type Message struct {
	ID           uint      `json:"id"`
	Content      string    `json:"content" binding:"required"`
	Sender       int       `json:"sender" binding:"required"`
	Receiver     int       `json:"receiver" binding:"required"`
	Read         bool      `json:"read"`
	Timestamp    time.Time `json:"timestamp"`
	Conversation string    `json:"conversation"`
}

type Conversation struct {
	ID    uint   `json:"id"`
	Key   string // hash_func ( user1_id, user2_id )
	User1 int    `json:"user1"`
	User2 int    `json:"user2"`
}
