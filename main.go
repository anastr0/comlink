package main

import (
	"comlink/message"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// TODO : better logging
func setupRouter() *gin.Engine {
	// create router
	router := gin.Default()

	// get message API handler - adds db access and message queue producer to handler
	comlink := message.GetMessagesAPIHandler()

	// user routes
	{
		router.POST("/user", comlink.CreateUserHandler) // create user
		router.GET("/user", comlink.GetUsersHandler)    // get all users
	}
	// message routes
	{
		router.POST("/message", comlink.SendMessageHandler)                 // send message
		router.GET("/message", comlink.RetrieveConversationHandler)         // retrieve conversation history between two users
		router.PATCH("/message/:id/read", comlink.MarkMessageAsReadHandler) // mark message as read
	}

	// TODO skip : close producer gracefully
	return router
}

func main() {
	// load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	// setup router
	r := setupRouter()
	r.Run(":8080")
}
