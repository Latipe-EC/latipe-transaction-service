package orderserv

import (
	"context"
	"latipe-transaction-service/internal/domain/message"
)

type OrderService interface {
	HandleTransactionPurchaseReply(ctx context.Context, message *message.CreateOrderReplyMessage, serviceType int) error
}
