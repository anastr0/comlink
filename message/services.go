package message

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO : move all to services pkg, along with consumer
// TODO : organise env file and vars for better usage across services, docker

// Sarama configuration options
var (
	brokers = os.Getenv("KAFKA_PEERS")
	version = sarama.DefaultVersion.String()
)

func GetMessagesHandler() *MessageHandler {
	// TODO : read secrets from env
	// init required services
	return &MessageHandler{
		db:       GetDB(),
		producer: getMessageProducer(),
	}
}

func GetDB() *gorm.DB {
	// TODO : read secrets from env
	// TODO : index message table by conversation id
	dsn := "host=localhost user=postgres password=example dbname=msg_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	// TODO : defer connection - avoid conn leaks https://go.dev/blog/defer-panic-and-recover
	if err != nil {
		log.Printf("error: %v %v\n", err, dsn)
		panic(err.Error())
	}
	return db
}

func getMessageProducer() sarama.AsyncProducer {
	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	brokers = os.Getenv("KAFKA_PEERS")
	brokerList := strings.Split(brokers, ",")

	// log.Printf("Kafka version: %s", version.String())
	config := sarama.NewConfig()
	config.Version = version
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
