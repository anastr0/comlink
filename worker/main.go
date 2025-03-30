package main

import (
	"comlink/message"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/IBM/sarama"
)

// TODO : run worker in docker
// TODO : pass brokers, topic as flags

// Sarama configuration options
var (
	brokers = "localhost:9092,"
	topic   = "messages"
)

func worker() {
	keepRunning := true

	// get db
	db := message.GetDB()

	config := sarama.NewConfig()
	brokerList := strings.Split(brokers, ",")

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		panic(err)
	}
	log.Println("consumer created successfully")

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Consumer initialization error: %v", err)
	}

	// Run infinite loop : keep consuming, until a kill signal is received

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Consume messages. write messages to db
	consumed := 0
	for keepRunning {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Message claimed: topic=%q partition=%d offset=%d key=%q value=%q consumer=%v\n",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value), consumed)

			var message_json message.Message

			if err := json.Unmarshal(msg.Value, &message_json); err != nil {
				log.Printf("error: %v\n", err)
			} else {
				// write consumed message to db
				// TODO : clear or offset message after consumption
				// TODO : create conversation id if not exists
				result := db.Create(&message_json)
				if result.Error != nil {
					log.Printf("error: %v\n", result.Error)
				} else {
					log.Printf("message created successfully\n")
				}
			}
			consumed++
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from consumer: %v\n", err)
		case s := <-signals:
			log.Printf("Signal received: %s, shutting down...\n", s)
			if err := partitionConsumer.Close(); err != nil {
				log.Fatalln(err)
			}
			keepRunning = false
		}
	}

}

func main() {
	// start worker service
	worker()
}
