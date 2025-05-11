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
