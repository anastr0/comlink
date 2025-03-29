package message

import (
	"crypto/sha256"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func (h *MessageHandler) RetrieveConversationHandler(c *gin.Context) {
	// get all messages
	var messages []Message
	result := h.db.Find(&messages)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "messages": messages})
	}
}

func cantorFunc(a, b int) string {
	// https://www.cantorsparadise.com/cantor-pairing-function-e213a8a89c2b
	cant_id := (a+b)*(a+b+1)/2 + b
	return strconv.Itoa(cant_id)
}

func (h *MessageHandler) SendMessageHandler(c *gin.Context) {
	// send message

	// TODO : check if sender and receiver exists
	// TODO : sender and receiver are query params
	// TODO : create conversation id, move to worker
	input1 := strconv.Itoa(1) + strconv.Itoa(2)
	first := sha256.New()
	first.Write([]byte(input1))

	var message_json Message

	if err := c.ShouldBindJSON(&message_json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if message_json.Sender == message_json.Receiver {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver are same"})
	} else {
		// result := h.db.Create(&message_json)

		message_json.Conversation = cantorFunc(message_json.Sender, message_json.Receiver)
		message := &Message{
			Content:      message_json.Content,
			Sender:       message_json.Sender,
			Receiver:     message_json.Receiver,
			Read:         false,
			Timestamp:    time.Now(),
			Conversation: message_json.Conversation,
		}

		// We will use the conversation_id as key. This will cause
		// all messages from the same conversation to end up
		// on the same partition (order is preserved).
		conversation_id := sarama.StringEncoder(message.Conversation)
		h.producer.Input() <- &sarama.ProducerMessage{
			Topic: "messages",
			Key:   conversation_id,
			Value: message,
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": message_json})
	}
}

func (h *MessageHandler) MarkMessageAsReadHandler(c *gin.Context) {
	// mark message as read
	id := c.Param("id")
	msg_id, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := Message{ID: uint(msg_id)}
	result := h.db.Model(&message).Update("read", true)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": message})
	}
}
