package subscriber

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPurchaseSubscriberReply,
	NewPurchaseCreateOrchestratorSubscriber,
)
