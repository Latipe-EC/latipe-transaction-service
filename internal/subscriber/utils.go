package subscriber

import (
	"github.com/gofiber/fiber/v2/log"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/service/orderserv"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func MappingRoutingKeyToService(routingKey string, config *config.Config) int {
	switch routingKey {
	case config.RabbitMQ.SagaOrderProductEvent.ReplyRoutingKey:
		return orderserv.PRODUCT_SERVICE
	case config.RabbitMQ.SagaOrderDeliveryEvent.ReplyRoutingKey:
		return orderserv.DELIVERY_SERVICE
	case config.RabbitMQ.SagaOrderPromotionEvent.ReplyRoutingKey:
		return orderserv.PROMOTION_SERVICE
	case config.RabbitMQ.SagaOrderPaymentEvent.ReplyRoutingKey:
		return orderserv.PAYMENT_SERVICE
	}
	return -1
}
