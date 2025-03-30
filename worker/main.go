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

// TODO : run worker in docker
// TODO : pass brokers, topic as flags

// input1 := strconv.Itoa(1) + strconv.Itoa(2)
// first := sha256.New()
// first.Write([]byte(input1))

func storeMessageInDB(db *gorm.DB, msg []byte) {
	// TODO : store message in db
	var message_json message.Message
	if err := json.Unmarshal(msg, &message_json); err != nil {
		log.Printf("error: %v\n", err)
	} else {
		// write consumed message to db
		// TODO : create conversation id if not exists

		result := db.Create(&message_json)
		if result.Error != nil {
			log.Printf("error: %v\n", result.Error)
		} else {
			log.Printf("message created successfully\n")
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
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err.Error())
	}

	// start worker service
	worker()
}
