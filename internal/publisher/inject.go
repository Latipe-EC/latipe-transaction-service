package publisher

import (
	"github.com/google/wire"
	"latipe-transaction-service/internal/publisher/createPurchase"
)

var Set = wire.NewSet(createPurchase.NewOrderOrchestratorPub)
