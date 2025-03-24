package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/entities"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/ports"
)

// cryptoService реализует сервис для работы с криптовалютами
type cryptoService struct {
	logger ports.Logger
}

// NewCryptoService создает новый экземпляр сервиса криптовалют
func NewCryptoService(logger ports.Logger) ports.CryptoService {
	return &cryptoService{
		logger: logger,
	}
}

// GetTopCurrencies получает информацию о топовых криптовалютах
func (s *cryptoService) GetTopCurrencies(limit int) ([]*entities.CryptoCurrency, error) {
	s.logger.Info("Запрос данных о криптовалютах")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,binancecoin,ripple,cardano&vs_currencies=usd&include_24hr_change=true")

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Ошибка при получении данных о криптовалютах", err)
		return nil, fmt.Errorf("ошибка получения данных криптовалют: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		s.logger.Error("Ошибка при декодировании ответа API", err)
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	
	s.logger.Info("Данные о криптовалютах успешно получены")

	// Формируем список криптовалют с полученными данными
	currencies := []*entities.CryptoCurrency{
		{Symbol: "BTC", Price: data["bitcoin"]["usd"], Change24h: data["bitcoin"]["usd_24h_change"]},
		{Symbol: "ETH", Price: data["ethereum"]["usd"], Change24h: data["ethereum"]["usd_24h_change"]},
		{Symbol: "BNB", Price: data["binancecoin"]["usd"], Change24h: data["binancecoin"]["usd_24h_change"]},
		{Symbol: "XRP", Price: data["ripple"]["usd"], Change24h: data["ripple"]["usd_24h_change"]},
		{Symbol: "ADA", Price: data["cardano"]["usd"], Change24h: data["cardano"]["usd_24h_change"]},
	}

	s.logger.Info("Возвращаем список из 5 криптовалют")
	return currencies, nil
}
