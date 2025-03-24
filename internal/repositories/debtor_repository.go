package repositories

import (
	"fmt"

	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/entities"
	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/ports"
)

// debtorRepository - репозиторий для работы с должниками
type debtorRepository struct {
	debtors []string
}

// NewDebtorRepository - создает новый экземпляр репозитория должников
func NewDebtorRepository(debtors []string) ports.DebtorRepository {
	return &debtorRepository{
		debtors: debtors,
	}
}

// GetAll - возвращает список всех должников
func (r *debtorRepository) GetAll() ([]*entities.Debtor, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	// Создаем список должников из имен пользователей
	result := make([]*entities.Debtor, len(r.debtors))
	for i, username := range r.debtors {
		result[i] = &entities.Debtor{Username: username}
	}
	return result, nil
}

// Validate - проверяет корректность данных о должниках
func (r *debtorRepository) Validate() error {
	if len(r.debtors) != 2 {
		return fmt.Errorf("ожидалось ровно 2 должника, получено %d", len(r.debtors))
	}

	for i, debtor := range r.debtors {
		if debtor == "" {
			return fmt.Errorf("должник %d пустой", i+1)
		}
	}

	return nil
}
