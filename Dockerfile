# Используем multi-stage build для оптимизации размера финального образа

# Этап 1: Сборка приложения
FROM node:20-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем исходный код
COPY . .

# Собираем TypeScript
RUN npm run build

# Этап 2: Создание финального образа
FROM node:20-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только необходимые файлы из этапа сборки
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/.env ./

# Устанавливаем tzdata для корректной работы с часовыми поясами
RUN apk add --no-cache tzdata

# Запускаем бот при старте контейнера
CMD ["npm", "start"]