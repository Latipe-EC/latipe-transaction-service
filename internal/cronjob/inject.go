package cronjob

import (
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

var Set = wire.NewSet(
	NewCronInstance,
	NewCheckingTxStatusCronJ,
)

func NewCronInstance() *cron.Cron {
	return cron.New()
}
