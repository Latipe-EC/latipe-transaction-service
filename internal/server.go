//go:build wireinject
// +build wireinject

package server

import (
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/wire"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/cronjob"
	"latipe-transaction-service/internal/domain/repos"
	"latipe-transaction-service/internal/publisher"
	"latipe-transaction-service/internal/service"
	"latipe-transaction-service/internal/subscriber"
	"latipe-transaction-service/pkgs/db/mongodb"
	"latipe-transaction-service/pkgs/rabbitclient"
)

type Server struct {
	app               *fiber.App
	globalCfg         *config.Config
	purchaseReplySub  *subscriber.PurchaseReplySubscriber
	purchaseCreateSub *subscriber.PurchaseCreateOrchestratorSubscriber
	checkTxStatus     *cronjob.CheckingTxStatusCronJ
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.globalCfg
}

func (serv Server) PurchaseReplySub() *subscriber.PurchaseReplySubscriber {
	return serv.purchaseReplySub
}

func (serv Server) PurchaseCreateSub() *subscriber.PurchaseCreateOrchestratorSubscriber {
	return serv.purchaseCreateSub
}

func (serv Server) CheckTxStatusCron() *cronjob.CheckingTxStatusCronJ {
	return serv.checkTxStatus
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		mongodb.Set,
		rabbitclient.Set,
		repos.Set,
		service.Set,
		subscriber.Set,
		publisher.Set,
		cronjob.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	purchaseReplySub *subscriber.PurchaseReplySubscriber,
	purchaseCreateSub *subscriber.PurchaseCreateOrchestratorSubscriber,
	checkTxStatus *cronjob.CheckingTxStatusCronJ) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	prometheus := fiberprometheus.New("latipe-order-service-v2")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// Initialize default config
	app.Use(logger.New())

	app.Get("", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "transaction service was developed by TienDat",
			Version: "v1.0.0",
		}
		return ctx.JSON(s)
	})

	return &Server{
		globalCfg:         cfg,
		app:               app,
		purchaseReplySub:  purchaseReplySub,
		purchaseCreateSub: purchaseCreateSub,
		checkTxStatus:     checkTxStatus,
	}
}
