package repos

import (
	"context"
	"latipe-transaction-service/internal/domain/entities"
)

type TransactionRepository interface {
	CreateTransactionData(ctx context.Context, dao *entities.TransactionLog) error
	FindByOrderID(ctx context.Context, orderID string) (*entities.TransactionLog, error)
	FindAllPendingTransaction(ctx context.Context) ([]*entities.TransactionLog, error)
	UpdateTransaction(ctx context.Context, dao *entities.TransactionLog, commit *entities.Commits) error
}
