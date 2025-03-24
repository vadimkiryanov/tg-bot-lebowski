package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/subosito/gotenv"
)

type Config struct {
	BotToken      string
	ChatID        int64
	ChatIDLogs    int64
	Debtors       []string
	StickerFileID string
	Debug         bool
}

func Load() (*Config, error) {
	if err := gotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки переменных окружения: %w", err)
	}

	chatID, err := strconv.ParseInt(os.Getenv("TG_CHAT_ID"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга ID чата: %w", err)
	}

	chatIDLogs, err := strconv.ParseInt(os.Getenv("TG_CHAT_ID_LOGS"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга ID чата для логов: %w", err)
	}

	debtors := []string{
		os.Getenv("DEBTOR_1"),
		os.Getenv("DEBTOR_2"),
	}

	return &Config{
		BotToken:      os.Getenv("TG_BOT_TOKEN"),
		ChatID:        chatID,
		ChatIDLogs:    chatIDLogs,
		Debtors:       debtors,
		StickerFileID: os.Getenv("TG_STICKER_FILE_ID"),
		Debug:         true,
	}, nil
}
