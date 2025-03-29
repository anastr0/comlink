package message

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetMessagesHandler() *MessageHandler {
	// TODO : add kafka message queue handler
	// TODO : read secrets from env
	dsn := "host=localhost user=postgres password=example dbname=msg_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	// TODO : defer connection - avoid conn leaks
	if err != nil {
		log.Printf("error: %v %v\n", err, dsn)
		panic(err.Error())
	}
	return &MessageHandler{db: db}
}
