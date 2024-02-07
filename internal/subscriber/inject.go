package subscriber

import (
	"github.com/google/wire"
	"latipe-transaction-service/internal/subscriber/cancelPurchase"
	"latipe-transaction-service/internal/subscriber/createPurchase"
)

var Set = wire.NewSet(
	createPurchase.NewPurchaseSubscriberReply,
	createPurchase.NewPurchaseCreateOrchestratorSubscriber,
	cancelPurchase.NewPurchaseCancelOrchestratorSubscriber,
)
