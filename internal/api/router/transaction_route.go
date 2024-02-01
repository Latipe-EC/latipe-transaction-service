package router

import (
	"github.com/gofiber/fiber/v2"
	"latipe-transaction-service/internal/api/handler"
	middleware "latipe-transaction-service/internal/api/middeware"
)

type TransactionRouter interface {
	Init(root *fiber.Router)
}

type transactionRouter struct {
	transHandler *handler.TransactionApiHandler
	middleware   *middleware.AuthMiddleware
}

func NewTransactionRouter(transHandler *handler.TransactionApiHandler, middleware *middleware.AuthMiddleware) TransactionRouter {
	return transactionRouter{
		transHandler: transHandler,
		middleware:   middleware,
	}
}

func (o transactionRouter) Init(root *fiber.Router) {
	trans := (*root).Group("/transactions")

	purchase := trans.Group("/purchase")
	purchase.Get("", o.middleware.RequiredRoles([]string{"ADMIN"}), o.transHandler.GetAllTransaction)
	purchase.Get("/:id", o.middleware.RequiredRoles([]string{"ADMIN"}), o.transHandler.GetTransactionById)
	purchase.Get("/queue/pending", o.middleware.RequiredRoles([]string{"ADMIN"}), o.transHandler.GetTransactionInQueue)

}
