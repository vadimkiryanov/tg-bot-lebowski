package services

import (
	"fmt"
	"time"

	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/entities"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/ports"
)

type notificationService struct {
	bot           ports.BotService
	debtorRepo    ports.DebtorRepository
	cryptoService ports.CryptoService
	chatID        int64
	chatIDLogs    int64
	stickerFileID string
	logger        ports.Logger
}

func NewNotificationService(
	bot ports.BotService,
	debtorRepo ports.DebtorRepository,
	cryptoService ports.CryptoService,
	chatID int64,
	chatIDLogs int64,
	stickerFileID string,
	logger ports.Logger,
) ports.NotificationService {
	return &notificationService{
		bot:           bot,
		debtorRepo:    debtorRepo,
		cryptoService: cryptoService,
		chatID:        chatID,
		chatIDLogs:    chatIDLogs,
		stickerFileID: stickerFileID,
		logger:        logger,
	}
}

func (s *notificationService) SendNotification(msg *entities.NotificationMessage) error {
	if err := s.bot.SendMessage(msg.ChatID, msg.Text); err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	if err := s.bot.SendSticker(msg.ChatID, msg.StickerID); err != nil {
		return fmt.Errorf("ошибка отправки стикера: %w", err)
	}

	return nil
}

func (s *notificationService) sendCryptoUpdate() error {
	currencies, err := s.cryptoService.GetTopCurrencies(5)
	if err != nil {
		return fmt.Errorf("ошибка получения данных криптовалют: %w", err)
	}

	message := "🚀 Обновление крипторынка:\n\n"
	for _, currency := range currencies {
		changeEmoji := "📈"
		if currency.Change24h < 0 {
			changeEmoji = "📉"
		}
		message += fmt.Sprintf("%s %s: $%.2f (%s%.2f%%)\n",
			changeEmoji,
			currency.Symbol,
			currency.Price,
			getChangePrefix(currency.Change24h),
			currency.Change24h,
		)
	}

	return s.bot.SendMessage(s.chatID, message)
}

func (s *notificationService) StartScheduledNotifications(location *time.Location) {
	for {
		now := time.Now().In(location)
		next := time.Date(now.Year(), now.Month(), now.Day(), 16, 00, 0, 0, location)
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		duration := next.Sub(now)
		s.logger.Info("Следующее уведомление через: ", duration)
		time.Sleep(duration)

		// Отправка уведомлений должникам
		debtors, err := s.debtorRepo.GetAll()
		if err != nil {
			s.logger.Error("Ошибка получения списка должников: ", err)
			continue
		}

		for _, debtor := range debtors {
			msg := &entities.NotificationMessage{
				ChatID:    s.chatID,
				Username:  debtor.Username,
				Text:      fmt.Sprintf("Как дела, <a href=\"https://t.me/%s\">Лебовски</a>?", debtor.Username),
				StickerID: s.stickerFileID,
			}

			if err := s.SendNotification(msg); err != nil {
				s.logger.Error("Ошибка отправки уведомления: ", err)
			}

			time.Sleep(time.Second)
		}

		// Отправка обновления криптовалют
		if err := s.sendCryptoUpdate(); err != nil {
			s.logger.Error("Ошибка отправки обновления криптовалют: ", err)
		}
	}
}

func getChangePrefix(change float64) string {
	if change > 0 {
		return "+"
	}
	return ""
}
