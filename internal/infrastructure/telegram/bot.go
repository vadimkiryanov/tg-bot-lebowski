package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/ports"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(token string, debug bool) (ports.BotService, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания экземпляра бота: %w", err)
	}

	api.Debug = debug

	return &Bot{api: api}, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true

	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}
	return nil
}

func (b *Bot) SendSticker(chatID int64, stickerID string) error {
	sticker := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID))
	_, err := b.api.Send(sticker)
	if err != nil {
		return fmt.Errorf("ошибка отправки стикера: %w", err)
	}
	return nil
}
