package msgqueue

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(NewOrderOrchestratorPub)
