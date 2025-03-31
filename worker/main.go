package main

import (
	"comlink/message"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func storeMessageInDB(db *gorm.DB, msg []byte) {
	var message_json message.Message
	if err := json.Unmarshal(msg, &message_json); err != nil {
		log.Printf("error: %v\n", err)
	} else {
		// write consumed message to db with conversation key
		// creates new conversation if none existent between sender and receiver

		// get pair of possible conversation_ids for sender and receiver
		conv_id1 := message.GetConversationID(message_json.Sender, message_json.Receiver)
		conv_id2 := message.GetConversationID(message_json.Receiver, message_json.Sender)

		var conv message.Conversation

		// check if a conversation exists in db between sender and receiver
		// if found, use existing conversation key
		// else, create new conversation with one of the generated conversation_ids
		conv_result := db.Limit(1).Find(&conv, "key = ? OR key = ?", conv_id1, conv_id2)
		if conv_result.Error != nil {
			log.Printf("Error fetching conversation: %v\n", conv_result.Error)
		}

		if conv.ID == 0 {
			// create conversation
			log.Printf("Conversation not found, creating conversation: %v\n", conv)
			conv = message.Conversation{
				User1: message_json.Sender,
				User2: message_json.Receiver,
				Key:   conv_id1,
			}
			conv_result = db.Create(&conv)
			if conv_result.Error != nil {
				log.Printf("Error creating conversation: %v\n", conv_result.Error)
			}
		} else {
			// conversation found
			log.Printf("Conversation found: %v\n", conv)
		}

		message_json.Conversation = conv.Key
		// get possible conversation ids
		result := db.Create(&message_json)
		if result.Error != nil {
			log.Printf("Error creating message: %v\n", result.Error)
		} else {
			log.Printf("Message written in db successfuly.\n")
		}
	}
}

func worker() {
	keepRunning := true

	// get db and consumer
	db := message.GetDB()
	consumer := message.GetPartitionConsumer()

	// Run infinite loop : keep consuming, until a kill signal is received
	// TODO skip : keep track of consumed messages, (exactly once)

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
	for keepRunning {
		select {
		case msg := <-consumer.Messages():
			log.Printf("Message claimed: topic=%q partition=%d offset=%d key=%q value=%q consumer=%v\n",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value), consumed)
			// Consume messages. write messages to db
			storeMessageInDB(db, msg.Value)
			consumed++
		case err := <-consumer.Errors():
			log.Printf("Error from consumer: %v\n", err)
		case s := <-signals:
			// Graceful shutdown
			log.Printf("Signal received: %s, shutting down...\n", s)
			if err := consumer.Close(); err != nil {
				log.Fatalln(err)
			}
			keepRunning = false
		}
	}

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	// start worker service
	worker()
}
