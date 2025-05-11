package message

import (
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
		db:       GetDB(),              // add database to message handler
		producer: getMessageProducer(), // add kafka message producer to message handler
	}
}

// func (h *MessageHandler) CloseProducer() {
// 	_ = h.producer.Close()
// }

func GetDB() *gorm.DB {
	// TODO : better usage of env variables
	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_sslmode := os.Getenv("POSTGRES_SSLMODE")
	postgres_timezone := os.Getenv("POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		postgres_host, postgres_user, postgres_password, postgres_db, postgres_port, postgres_sslmode, postgres_timezone)

	// open db connection with gorm
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
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
		log.Fatalf("Failed to start Sarama producer: %v\n", err)
	}

	// messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to write messages to kafka: %v", err)
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
		log.Fatalf("Failed to start Sarama consumer: %v\n", err)
	}
	log.Println("consumer created successfully")

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failure to initialise partition consumer: %v", err)
	}
	return partitionConsumer
}
