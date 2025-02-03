# Используем multi-stage build для оптимизации размера финального образа

# Этап 1: Сборка приложения
FROM golang:1.23.5-alpine3.20 AS builder

# Копируем исходный код проекта в контейнер
COPY . /github.com/vadimkiryanov/tg-bot-lebowski/
# Устанавливаем рабочую директорию
WORKDIR /github.com/vadimkiryanov/tg-bot-lebowski/

# Загружаем зависимости проекта
RUN go mod download
# Собираем бинарный файл для Linux
RUN GOOS=linux go build -o ./.bin/bot ./cmd/bot/main.go

# Этап 2: Создание финального образа
FROM alpine:latest

# Устанавливаем рабочую директорию в корневой каталог
WORKDIR /root/

# Устанавливаем tzdata для корректной работы с часовыми поясами
# Это важно для логирования и временных операций
RUN apk add --no-cache tzdata

# Копируем только собранный бинарный файл из предыдущего этапа
# --from=0 указывает на первый этап сборки (builder)
COPY --from=0 /github.com/vadimkiryanov/tg-bot-lebowski/.bin/bot .
# Копируем конфигурационный файл
COPY .env .

# Объявляем порт, который будет использоваться приложением
EXPOSE 80

# Запускаем бот при старте контейнера
CMD ["./bot"]