package service

import (
	"github.com/google/wire"
	"latipe-transaction-service/internal/service/orderserv"
	"latipe-transaction-service/internal/service/transactionserv"
)

var Set = wire.NewSet(orderserv.NewOrderService,
	transactionserv.NewTransactionService)
