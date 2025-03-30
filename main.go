package main

import (
	"comlink/message"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// TODO : run API with docker
// TODO : better logging

func setupRouter() *gin.Engine {
	// create router
	router := gin.Default()

	// get message handler - adds db access and message queue producer to handler
	comlink := message.GetMessagesHandler()

	// user routes
	{
		router.POST("/user", comlink.CreateUserHandler)
		router.GET("/user", comlink.GetUsersHandler)
	}
	// message routes
	{
		router.POST("/message", comlink.SendMessageHandler)
		router.GET("/message", comlink.RetrieveConversationHandler)
		router.PATCH("/message/:id", comlink.MarkMessageAsReadHandler)
	}

	// TODO : close producer gracefully
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
