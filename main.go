package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type Message struct {
	ID          int       `json:"id"`
	Content     string    `json:"content" binding:"required"`
	Sender_ID   int       `json:"sender_id" binding:"required"`
	Receiver_ID int       `json:"receiver_id" binding:"required"`
	Read        bool      `json:"read"`
	Timestamp   time.Time `json:"timestamp"`
}

var db = make(map[string]User)
var msg_db = make(map[string]Message)
var user_id = 1
var msg_id = 7070

func createUserHandler(c *gin.Context) {
	// create user
	var user_json User

	if err := c.ShouldBindJSON(&user_json); err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		user_json.ID = user_id
		db[user_json.Name] = user_json
		user_id = user_id + 1
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func sendMessageHandler(c *gin.Context) {
	// send message

	// TODO : check if sender and receiver exists
	var message_json Message

	if err := c.ShouldBindJSON(&message_json); err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if message_json.Sender_ID == message_json.Receiver_ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver are same"})
	} else {
		message_json.ID = msg_id
		message_json.Timestamp = time.Now()
		msg_id = msg_id + 1
		id := strconv.Itoa(msg_id)
		msg_db[id] = message_json
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": message_json})
	}
}

func markMessageAsReadHandler(c *gin.Context) {
	// mark message as read
	id := c.Param("id")
	msg := msg_db[id]
	msg.Read = true
	msg_db[id] = msg
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": msg})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/user", createUserHandler)
	r.GET("/user", func(c *gin.Context) {
		c.JSON(http.StatusOK, db)
	})

	r.POST("/message", sendMessageHandler)
	r.GET("/message", func(c *gin.Context) {
		c.JSON(http.StatusOK, msg_db)
	})
	r.PATCH("/message/:id", markMessageAsReadHandler)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
