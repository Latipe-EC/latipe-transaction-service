package transactionserv

import (
	"context"
	"latipe-transaction-service/internal/domain/dto"
)

type TransactionService interface {
	GetAllTransaction(ctx context.Context, req *dto.GetAllTransactionLogRequest) (*dto.GetAllTransactionLogResponse, error)
	GetDetailTransaction(ctx context.Context, req *dto.GetTransactionLogByIdRequest) (*dto.GetTransactionLogByIdResponse, error)
	CheckTransactionPending(ctx context.Context) (*dto.TransactionPendingResponse, error)
}
