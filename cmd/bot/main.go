package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

var (
	botToken      string // токен бота
	chatID        int64  // ID основного чата
	chatIDLogs    int64  // ID чата для ошибок
	debtors       = make([]string, 0)
	stickerFileID string

	bot *tgbotapi.BotAPI
)

func init() {
	// Initialize bot and configuration first
	if err := initBot(); err != nil {
		logrus.Fatalf("Failed to initialize bot: %v", err)
	}

	// Validate debtors after bot is initialized
	if err := validateDebtors(); err != nil {
		errorMsg := fmt.Sprintf("Error validating debtors: %v", err)
		logrus.Error(errorMsg)
		sendLogsMessage(errorMsg)
		logrus.Fatal("Terminating due to invalid debtors configuration")
	}

	// Debug logging
	logrus.WithFields(logrus.Fields{
		"chatID":  chatID,
		"debtors": debtors,
	}).Info("Bot initialized successfully")
}

func main() {
	// Устанавливаем московское время
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logrus.Fatalf("Ошибка при установке временной зоны: %v", err)
		location = time.UTC
	}

	for {
		now := time.Now().In(location)
		// Вычисляем время до следующих 15:00
		next := time.Date(now.Year(), now.Month(), now.Day(), 16, 00, 0, 0, location)
		if now.After(next) {
			// Если текущее время после 15:00, переходим на следующий день
			next = next.Add(24 * time.Hour)
		}

		// Ждем до следующего запланированного времени
		duration := next.Sub(now)
		logrus.Infof("Следующая отправка через: %v", duration)
		time.Sleep(duration)

		// Отправляем сообщения всем должникам
		for _, username := range debtors {
			if err := sendStickerAndMessage(bot, chatID, username); err != nil {
				logrus.Errorf("Ошибка при отправке: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// Инициализация бота
func initBot() error {
	// Load environment variables first
	if err := gotenv.Load(); err != nil {
		return fmt.Errorf("error loading env variables: %w", err)
	}

	// Get bot token first and initialize bot
	botToken = os.Getenv("TG_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("TG_BOT_TOKEN is not set")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("failed to create bot instance: %w", err)
	}

	// Get and parse chat IDs
	chatIDStr := os.Getenv("TG_CHAT_ID")
	chatIDLogsStr := os.Getenv("TG_CHAT_ID_LOGS")
	stickerFileID = os.Getenv("TG_STICKER_FILE_ID")

	if chatID, err = strconv.ParseInt(chatIDStr, 10, 64); err != nil {
		return fmt.Errorf("failed to parse chat ID: %w", err)
	}
	if chatIDLogs, err = strconv.ParseInt(chatIDLogsStr, 10, 64); err != nil {
		return fmt.Errorf("failed to parse error chat ID: %w", err)
	}

	bot.Debug = true                                            // включаем режим отладки
	logrus.Infof("Authorized on account %s", bot.Self.UserName) // выводим в консоль имя бота

	return nil
}

// Отправляем стикер и сообщение
func sendStickerAndMessage(bot *tgbotapi.BotAPI, chatID int64, username string) error {

	// Отправляем сообщение
	msgText := fmt.Sprintf("Как дела, <a href=\"https://t.me/%s\">Лебовски</a>?", username) // формируем сообщение", username)

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = "HTML" // Включаем поддержку HTML-разметки
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	// Отправляем стикер
	sticker := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerFileID))
	_, err = bot.Send(sticker)
	if err != nil {
		return fmt.Errorf("error sending sticker: %v", err)
	}

	return nil
}

// Валидация debtors
func validateDebtors() error {
	debtors = []string{
		os.Getenv("DEBTOR_1"),
		os.Getenv("DEBTOR_2"),
	}

	if len(debtors) != 2 {
		return fmt.Errorf("expected exactly 2 debtors, got %d", len(debtors))
	}

	for i, debtor := range debtors {
		if debtor == "" {
			return fmt.Errorf("debtor %d is empty", i+1)
		}
	}

	return nil
}

// Отправляем сообщение об ошибке
func sendLogsMessage(message string) {
	if bot == nil {
		logrus.Error("Cannot send error message: bot is not initialized")
		return
	}

	msg := tgbotapi.NewMessage(chatIDLogs, message)
	if _, err := bot.Send(msg); err != nil {
		logrus.Errorf("Failed to send error message: %v", err)
	}
}
