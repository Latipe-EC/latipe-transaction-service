package orderserv

import (
	"context"
	"latipe-transaction-service/internal/domain/entities"
	"latipe-transaction-service/internal/domain/message"
)

type OrderService interface {
	StartPurchaseTransaction(ctx context.Context, message *message.OrderPendingMessage) error
	CancelOrder(ctx context.Context, message *message.CancelOrderMessage) error
	HandleTransactionPurchaseReply(ctx context.Context, message *message.CreateOrderReplyMessage, serviceType int) error
	RollbackTransactionPub(dao *entities.TransactionLog) error
	CheckTransactionStatus(ctx context.Context) error
}
