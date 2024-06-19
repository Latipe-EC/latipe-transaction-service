//go:build wireinject
// +build wireinject

package server

import (
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	recoverFiber "github.com/gofiber/fiber/v2/middleware/recover"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/adapter"
	"latipe-transaction-service/internal/api/handler"
	middleware "latipe-transaction-service/internal/api/middeware"
	"latipe-transaction-service/internal/api/router"
	"latipe-transaction-service/internal/cronjob"
	"latipe-transaction-service/internal/domain/repos"
	"latipe-transaction-service/internal/publisher"
	"latipe-transaction-service/internal/service"
	"latipe-transaction-service/internal/service/notifyserv"
	"latipe-transaction-service/internal/subscriber"
	"latipe-transaction-service/internal/subscriber/cancelPurchase"
	"latipe-transaction-service/internal/subscriber/createPurchase"
	"latipe-transaction-service/pkgs/cache"
	"latipe-transaction-service/pkgs/db/mongodb"
	"latipe-transaction-service/pkgs/rabbitclient"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/wire"
)

type Server struct {
	app               *fiber.App
	globalCfg         *config.Config
	purchaseReplySub  *createPurchase.PurchaseReplySubscriber
	purchaseCreateSub *createPurchase.PurchaseCreateOrchestratorSubscriber
	purchaseCancelSub *cancelPurchase.PurchaseCancelOrchestratorSubscriber
	checkTxStatus     *cronjob.CheckingTxStatusCronJ
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.globalCfg
}

func (serv Server) PurchaseReplySub() *createPurchase.PurchaseReplySubscriber {
	return serv.purchaseReplySub
}

func (serv Server) PurchaseCreateSub() *createPurchase.PurchaseCreateOrchestratorSubscriber {
	return serv.purchaseCreateSub
}

func (serv Server) PurchaseCancelSub() *cancelPurchase.PurchaseCancelOrchestratorSubscriber {
	return serv.purchaseCancelSub
}

func (serv Server) CheckTxStatusCron() *cronjob.CheckingTxStatusCronJ {
	return serv.checkTxStatus
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		cache.Set,
		mongodb.Set,
		rabbitclient.Set,
		repos.Set,
		notifyserv.Set,
		service.Set,
		adapter.Set,
		handler.Set,
		middleware.Set,
		router.Set,
		subscriber.Set,
		publisher.Set,
		cronjob.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	transRouter router.TransactionRouter,
	purchaseReplySub *createPurchase.PurchaseReplySubscriber,
	purchaseCreateSub *createPurchase.PurchaseCreateOrchestratorSubscriber,
	purchaseCancelSub *cancelPurchase.PurchaseCancelOrchestratorSubscriber,
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

	recoverConfig := recoverFiber.ConfigDefault
	app.Use(recoverFiber.New(recoverConfig))

	app.Get("", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "transaction service was developed by tdatIT",
			Version: "v1.0.0",
		}
		return ctx.JSON(s)
	})
	api := app.Group("/api")
	v1 := api.Group("/v1")
	transRouter.Init(&v1)

	return &Server{
		globalCfg:         cfg,
		app:               app,
		purchaseReplySub:  purchaseReplySub,
		purchaseCreateSub: purchaseCreateSub,
		purchaseCancelSub: purchaseCancelSub,
		checkTxStatus:     checkTxStatus,
	}
}
