// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/adapter/userserv"
	"latipe-transaction-service/internal/api/handler"
	"latipe-transaction-service/internal/api/middeware"
	"latipe-transaction-service/internal/api/router"
	"latipe-transaction-service/internal/cronjob"
	"latipe-transaction-service/internal/domain/repos"
	"latipe-transaction-service/internal/publisher/createPurchase"
	"latipe-transaction-service/internal/service/notifyserv"
	"latipe-transaction-service/internal/service/orderserv"
	"latipe-transaction-service/internal/service/transactionserv"
	"latipe-transaction-service/internal/subscriber/cancelPurchase"
	createPurchase2 "latipe-transaction-service/internal/subscriber/createPurchase"
	"latipe-transaction-service/pkgs/cache"
	"latipe-transaction-service/pkgs/db/mongodb"
	"latipe-transaction-service/pkgs/rabbitclient"
)

// Injectors from server.go:

func New() (*Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	mongoClient, err := mongodb.OpenMongoDBConnection(configConfig)
	if err != nil {
		return nil, err
	}
	transactionRepository := repos.NewTransactionRepository(mongoClient)
	cacheEngine, err := cache.NewCacheEngine(configConfig)
	if err != nil {
		return nil, err
	}
	transactionService := transactionserv.NewTransactionService(transactionRepository, cacheEngine)
	transactionApiHandler := handler.NewTransactionHandler(transactionService)
	userService := userserv.NewUserService(configConfig)
	authMiddleware := middleware.NewAuthMiddleware(userService)
	transactionRouter := router.NewTransactionRouter(transactionApiHandler, authMiddleware)
	connection := rabbitclient.NewRabbitClientConnection(configConfig)
	orderOrchestratorPub := createPurchase.NewOrderOrchestratorPub(configConfig, connection)
	notifyService := notifyserv.NewTelegramBot(configConfig)
	orderService := orderserv.NewOrderService(transactionRepository, orderOrchestratorPub, notifyService, cacheEngine)
	purchaseReplySubscriber := createPurchase2.NewPurchaseSubscriberReply(configConfig, orderService, connection)
	purchaseCreateOrchestratorSubscriber := createPurchase2.NewPurchaseCreateOrchestratorSubscriber(configConfig, orderService, connection)
	purchaseCancelOrchestratorSubscriber := cancelPurchase.NewPurchaseCancelOrchestratorSubscriber(configConfig, orderService, connection)
	cron := cronjob.NewCronInstance()
	checkingTxStatusCronJ := cronjob.NewCheckingTxStatusCronJ(configConfig, cron, orderService)
	server := NewServer(configConfig, transactionRouter, purchaseReplySubscriber, purchaseCreateOrchestratorSubscriber, purchaseCancelOrchestratorSubscriber, checkingTxStatusCronJ)
	return server, nil
}

// server.go:

type Server struct {
	app               *fiber.App
	globalCfg         *config.Config
	purchaseReplySub  *createPurchase2.PurchaseReplySubscriber
	purchaseCreateSub *createPurchase2.PurchaseCreateOrchestratorSubscriber
	purchaseCancelSub *cancelPurchase.PurchaseCancelOrchestratorSubscriber
	checkTxStatus     *cronjob.CheckingTxStatusCronJ
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.globalCfg
}

func (serv Server) PurchaseReplySub() *createPurchase2.PurchaseReplySubscriber {
	return serv.purchaseReplySub
}

func (serv Server) PurchaseCreateSub() *createPurchase2.PurchaseCreateOrchestratorSubscriber {
	return serv.purchaseCreateSub
}

func (serv Server) PurchaseCancelSub() *cancelPurchase.PurchaseCancelOrchestratorSubscriber {
	return serv.purchaseCancelSub
}

func (serv Server) CheckTxStatusCron() *cronjob.CheckingTxStatusCronJ {
	return serv.checkTxStatus
}

func NewServer(
	cfg *config.Config,
	transRouter router.TransactionRouter,
	purchaseReplySub *createPurchase2.PurchaseReplySubscriber,
	purchaseCreateSub *createPurchase2.PurchaseCreateOrchestratorSubscriber,
	purchaseCancelSub *cancelPurchase.PurchaseCancelOrchestratorSubscriber,
	checkTxStatus *cronjob.CheckingTxStatusCronJ) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:5500, http://127.0.0.1:5173, http://localhost:5500, http://localhost:5173",
		AllowHeaders: "*",
		AllowMethods: "GET,HEAD,OPTIONS,POST,PUT",
	}))

	prometheus := fiberprometheus.New("latipe-order-service-v2")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(logger.New())

	recoverConfig := recover2.ConfigDefault
	app.Use(recover2.New(recoverConfig))

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
