FROM golang:1.23.5-alpine3.20 AS  builder

COPY . /github.com/vadimkiryanov/tg-bot-lebowski/
WORKDIR /github.com/vadimkiryanov/tg-bot-lebowski/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/bot ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

# Установка tzdata для поддержки часовых поясов, чтобы бот работал в любом часовом поясе
RUN apk add --no-cache tzdata

# Копируем бинарный файл из предыдущего этапа
COPY --from=0 /github.com/vadimkiryanov/tg-bot-lebowski/.bin/bot .
# Копируем .env файл в контейнер
COPY .env .

EXPOSE 80

CMD ["./bot"]