// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/domain/repos"
	"latipe-transaction-service/internal/publisher"
	"latipe-transaction-service/internal/service/orderserv"
	"latipe-transaction-service/internal/subscriber"
	"latipe-transaction-service/pkgs/db/mongodb"
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
	orderOrchestratorPub := publisher.NewOrderOrchestratorPub(configConfig)
	orderService := orderserv.NewOrderService(transactionRepository, orderOrchestratorPub)
	purchaseReplySubscriber := subscriber.NewPurchaseSubscriberReply(configConfig, orderService)
	purchaseCreateOrchestratorSubscriber := subscriber.NewPurchaseCreateOrchestratorSubscriber(configConfig, orderService)
	server := NewServer(configConfig, purchaseReplySubscriber, purchaseCreateOrchestratorSubscriber)
	return server, nil
}

// server.go:

type Server struct {
	app               *fiber.App
	globalCfg         *config.Config
	purchaseReplySub  *subscriber.PurchaseReplySubscriber
	purchaseCreateSub *subscriber.PurchaseCreateOrchestratorSubscriber
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

func NewServer(
	cfg *config.Config,
	purchaseReplySub *subscriber.PurchaseReplySubscriber,
	purchaseCreateSub *subscriber.PurchaseCreateOrchestratorSubscriber) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	prometheus := fiberprometheus.New("latipe-order-service-v2")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

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
	}
}
