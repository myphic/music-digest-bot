# Music digest telegram bot
Телеграм бот на GO, который отправляет актуальные новинки музыки в тг канал

## Содержание
- [Технологии](#технологии)
- [Развертывание](#развертывание)
- [To do](#to-do)


## Технологии
- Golang
- Docker
- Postgres
- Goose
- golangci-lint, lefthook

## Развертывание

1) Создать .env файл по аналогии с .env-example
```sh
$ cd build & docker-compose up
$ goose -dir internal/db/migrations up
$ go run cmd/main.go
```

## To do
- [x] Apple-music, soundcloud, AI summary
