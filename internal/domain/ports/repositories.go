package ports

import "github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/entities"

type DebtorRepository interface {
	GetAll() ([]*entities.Debtor, error)
	Validate() error
}
