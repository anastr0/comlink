package message

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func (h *MessagesAPIHandler) RetrieveConversationHandler(c *gin.Context) {
	// get conversation history

	// TODO : pagination
	if c.Query("user1") == "" || c.Query("user2") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user1 and user2 are required"})
		return
	}
	user1, err1 := strconv.Atoi(c.Query("user1"))
	user2, err2 := strconv.Atoi(c.Query("user2"))

	if err1 != nil || err2 != nil {
		log.Printf("errors: %v %v\n", err1, err2)
	}

	// get conversation ids, query by conversation ids
	conv_key1 := GetConversationID(user1, user2)
	conv_key2 := GetConversationID(user2, user1)

	var messages []Message
	result := h.db.Limit(10).Where(
		"conversation = ? OR conversation = ?", conv_key1, conv_key2,
	).Order("id	desc").Find(&messages)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}

func GetConversationID(sender_id, receiver_id int) string {
	// generate md5 conversation_id from sender_id and receiver_id
	// uses cantor pairing func https://www.cantorsparadise.com/cantor-pairing-function-e213a8a89c2b

	h := md5.New()
	cant_id := (sender_id+receiver_id)*(sender_id+receiver_id+1)/2 + receiver_id
	h.Write([]byte(strconv.Itoa(cant_id)))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *MessagesAPIHandler) SendMessageHandler(c *gin.Context) {
	// send message

	// TODO skip: check if sender and receiver exists
	var message_json Message

	if err := c.ShouldBindJSON(&message_json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if message_json.Sender == message_json.Receiver {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver are same"})
	} else {
		// producer : push message to kafka
		message := &Message{
			Content:   message_json.Content,
			Sender:    message_json.Sender,
			Receiver:  message_json.Receiver,
			Read:      false,
			Timestamp: time.Now(),
		}

		// TODO skip : We will use the conversation_id as key. This will cause
		// all messages from the same conversation to end up
		// on the same partition (order is preserved).
		// conversation_id := sarama.StringEncoder(message.Conversation)
		// Key: conversation_id
		h.producer.Input() <- &sarama.ProducerMessage{
			Topic: "messages",
			Value: message,
		}
		log.Printf("Message sent to Kafka: %v\n", message)
		c.JSON(http.StatusCreated, gin.H{"status": "message sent"})
	}
}

func (h *MessagesAPIHandler) MarkMessageAsReadHandler(c *gin.Context) {
	// mark message as read
	id := c.Param("id")
	msg_id, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var message Message

	// check if message exists, return error if not
	h.db.Where("id = ?", msg_id).Limit(1).Find(&message)

	if message.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
	} else {
		// message_id exists, mark message as read
		result := h.db.Model(&message).Update("read", true)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		} else {
			c.JSON(http.StatusNoContent, gin.H{"status": "read"})
		}
	}
}
