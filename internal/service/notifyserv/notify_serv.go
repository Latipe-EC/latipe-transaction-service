package notifyserv

import "latipe-transaction-service/internal/domain/entities"

type NotifyService interface {
	SendMessageToTelegram(trans *entities.TransactionLog, reason string) error
}
