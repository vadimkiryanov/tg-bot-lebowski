package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

var (
	botToken string
	chatID   int64
	debtors  = make([]string, 0)
	// Можете заменить на ID вашего стикера
	stickerFileID string
)

func init() {
	// Инициализация переменных окружения
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: [%s]\n", err)
	}

	// Convert string chat ID to int64
	chatIDStr := os.Getenv("TG_CHAT_ID")
	botToken = os.Getenv("TG_BOT_TOKEN")
	stickerFileID = os.Getenv("TG_STICKER_FILE_ID")

	debtors = append(debtors, os.Getenv("DEBTOR_1"), os.Getenv("DEBTOR_2")) // добавляем должников

	var err error
	fmt.Printf("chatIDStr: %v\n", chatIDStr)
	fmt.Printf("debtors: %v\n", debtors)
	fmt.Printf("botToken: %v\n", botToken)

	if chatID, err = strconv.ParseInt(chatIDStr, 10, 64); err != nil {
		logrus.Fatalf("failed to parse chat ID: %v", err)
	}

}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken) // получаем экземпляр бота
	if err != nil {
		logrus.Fatalf("error creating bot: %v", err)
	}

	bot.Debug = true                                          // включаем режим отладки
	log.Printf("Authorized on account %s", bot.Self.UserName) // выводим в консоль имя бота

	// Устанавливаем московское время
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Printf("Ошибка при установке временной зоны: %v", err)
		location = time.UTC
	}

	for {
		now := time.Now().In(location)
		// Вычисляем время до следующих 15:00
		next := time.Date(now.Year(), now.Month(), now.Day(), 20, 50, 0, 0, location)
		if now.After(next) {
			// Если текущее время после 15:00, переходим на следующий день
			next = next.Add(24 * time.Hour)
		}

		// Ждем до следующего запланированного времени
		duration := next.Sub(now)
		log.Printf("Следующая отправка через: %v", duration)
		time.Sleep(duration)

		// Отправляем сообщения всем должникам
		for _, username := range debtors {
			if err := sendStickerAndMessage(bot, chatID, username); err != nil {
				log.Printf("Ошибка при отправке: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func sendStickerAndMessage(bot *tgbotapi.BotAPI, chatID int64, username string) error {

	// Отправляем сообщение
	msgText := fmt.Sprintf("%s Как дела, Лебовски?", username)
	msg := tgbotapi.NewMessage(chatID, msgText)
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
