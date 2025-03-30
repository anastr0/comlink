# comlink

Real time messaging platform

![https://starwars.fandom.com/wiki/C1_personal_comlink](https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQiUCgC5WObaqB8zTWtzfQXJIdztAVVoqtT2MBylAPwn_JVJLbl15iRQIIwX1ui30p4PL8&usqp=CAU)


# Run in local

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
