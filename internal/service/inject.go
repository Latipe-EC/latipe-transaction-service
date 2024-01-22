package service

import (
	"github.com/google/wire"
	"latipe-transaction-service/internal/service/orderserv"
)

var Set = wire.NewSet(orderserv.NewOrderService)
