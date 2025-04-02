package main

import (
	"log"
	"time"

	"github.com/vadimkiryanov/tg-bot-lebowski/internal/config"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/infrastructure/logger"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/infrastructure/telegram"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/repositories"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/services"
)

func main() {
	// Инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инцициализация логгера.
	log := logger.NewLogger()

	// Инцициализация бота
	bot, err := telegram.NewBot(cfg.BotToken, cfg.Debug)
	if err != nil {
		log.Fatal("Ошибка инициализации бота: %v", err)
	}

	// Инцициализация репозиториев
	debtorRepo := repositories.NewDebtorRepository(cfg.Debtors)

	// Инцициализация сервисов
	cryptoService := services.NewCryptoService(log)
	notificationService := services.NewNotificationService(
		bot,
		debtorRepo,
		cryptoService,
		cfg.ChatID,
		cfg.ChatIDLogs,
		cfg.StickerFileID,
		log,
	)

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal("Ошибка получения локации: %v", err)
	}

	// Запуск уведомлений
	notificationService.StartScheduledNotifications(location)
}
