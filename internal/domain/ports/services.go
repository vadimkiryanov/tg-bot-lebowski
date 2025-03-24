package ports

import (
	"time"

	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/entities"
)

type NotificationService interface {
	SendNotification(msg *entities.NotificationMessage) error
	StartScheduledNotifications(location *time.Location)
}

type BotService interface {
	SendMessage(chatID int64, text string) error
	SendSticker(chatID int64, stickerID string) error
}

type CryptoService interface {
	GetTopCurrencies(limit int) ([]*entities.CryptoCurrency, error)
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}
