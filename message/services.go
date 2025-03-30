package message

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetMessagesAPIHandler() *MessagesAPIHandler {
	// init required services for message API handler
	return &MessagesAPIHandler{
		db:       GetDB(),
		producer: getMessageProducer(),
	}
}

// func (h *MessageHandler) CloseProducer() {
// 	_ = h.producer.Close()
// }

func GetDB() *gorm.DB {
	flag.Parse()
	// TODO : index message table by conversation id
	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_sslmode := os.Getenv("POSTGRES_SSLMODE")
	postgres_timezone := os.Getenv("POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		postgres_host, postgres_user, postgres_password, postgres_db, postgres_port, postgres_sslmode, postgres_timezone)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Printf("error: %v %v\n", err, dsn)
		panic(err.Error())
	}
	return db
}

func getMessageProducer() sarama.AsyncProducer {
	brokerList := strings.Split(os.Getenv("KAFKA_PEERS"), ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write messages to kafka:", err)
		}
	}()

	return producer
}

func GetPartitionConsumer() sarama.PartitionConsumer {
	brokerList := strings.Split(os.Getenv("KAFKA_PEERS"), ",")
	topic := "messages"

	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("consumer created successfully")

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Consumer initialization error: %v", err)
	}
	return partitionConsumer
}
