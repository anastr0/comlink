package message

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type MessagesAPIHandler struct {
	db       *gorm.DB
	producer sarama.AsyncProducer
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

	encoded []byte
	err     error
}

func (msg *Message) ensureEncoded() {
	if msg.encoded == nil && msg.err == nil {
		msg.encoded, msg.err = json.Marshal(msg)
	}
}

func (msg *Message) Length() int {
	msg.ensureEncoded()
	return len(msg.encoded)
}

func (msg *Message) Encode() ([]byte, error) {
	msg.ensureEncoded()
	return msg.encoded, msg.err
}

type Conversation struct {
	ID    uint   `json:"id"`
	Key   string // hash_func ( user1_id, user2_id )
	User1 uint   `json:"user1"`
	User2 uint   `json:"user2"`
}
