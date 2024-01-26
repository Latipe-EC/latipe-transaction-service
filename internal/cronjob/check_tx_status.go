package cronjob

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/robfig/cron/v3"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/service/orderserv"
	"sync"
	"time"
)

type CheckingTxStatusCronJ struct {
	cfg          *config.Config
	cron         *cron.Cron
	orderService orderserv.OrderService
}

func NewCheckingTxStatusCronJ(cfg *config.Config, cron *cron.Cron, orderServ orderserv.OrderService) *CheckingTxStatusCronJ {
	return &CheckingTxStatusCronJ{
		cfg:          cfg,
		cron:         cron,
		orderService: orderServ,
	}
}

func (cr *CheckingTxStatusCronJ) StartJob(wg *sync.WaitGroup) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("checking transaction status")
	job, err := cr.cron.AddFunc(cr.cfg.CronJob.CheckingTxStatus, func() {
		if err := cr.orderService.CheckTransactionStatus(ctxTimeout); err != nil {
			log.Error(err)
		}
	})

	if err != nil {
		cr.cron.Remove(job)
		return err
	}

	cr.cron.Run()

	defer wg.Done()
	return nil
}
