package orderserv

import (
	"context"
	"latipe-transaction-service/internal/domain/message"
	"latipe-transaction-service/internal/domain/repos"
	msgqueue "latipe-transaction-service/internal/publisher"
)

type orderService struct {
	transactionRepo         repos.TransactionRepository
	transactionOrchestrator *msgqueue.OrderOrchestratorPub
}

func NewOrderService(transRepos repos.TransactionRepository, orchestrator *msgqueue.OrderOrchestratorPub) OrderService {
	return &orderService{
		transactionRepo:         transRepos,
		transactionOrchestrator: orchestrator,
	}
}

func (o orderService) HandleTransactionPurchaseReply(ctx context.Context, message *message.CreateOrderReplyMessage, serviceType int) error {
	//TODO implement me
	panic("implement me")
}
