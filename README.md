# comlink

A real time messaging API

![https://starwars.fandom.com/wiki/C1_personal_comlink](https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQiUCgC5WObaqB8zTWtzfQXJIdztAVVoqtT2MBylAPwn_JVJLbl15iRQIIwX1ui30p4PL8&usqp=CAU)

https://starwars.fandom.com/wiki/C1_personal_comlink/Legends

## Run in local

Run postgres db and kafka
```bash
docker-compose up -d
```

To run API
```bash
go run main.go
```

To run worker service (consumer)

```bash
go run worker/main.go
```

## Test

post messages

```bash
curl -s -X POST http://localhost:8080/message \
    -H 'content-type: application/json' \
    -d '{"content": "hi", "sender": 5, "receiver": 6}' | jq

curl -s -X POST http://localhost:8080/message \
    -H 'content-type: application/json' \
    -d '{"content": "how are you?", "sender": 5, "receiver": 6}' | jq

curl -s -X POST http://localhost:8080/message \
    -H 'content-type: application/json' \
    -d '{"content": "All good. How about you?", "sender": 6, "receiver": 5}' | jq

curl -s -X POST http://localhost:8080/message \
    -H 'content-type: application/json' \
    -d '{"content": "Same here", "sender": 5, "receiver": 6}' | jq
```

retrieve conversation

```bash
curl -X GET "http://localhost:8080/message?user1=5&user2=6" | jq
```

mark message read
```bash
curl -s -X PATCH http://localhost:8080/message/:id/read | jq
```