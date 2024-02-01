package orderserv

import "time"

const (
	PRODUCT_SERVICE   = 1
	PROMOTION_SERVICE = 2
	DELIVERY_SERVICE  = 3
	PAYMENT_SERVICE   = 4
)

const (
	PRODUCT_SERVICE_NAME   = "product_service"
	PROMOTION_SERVICE_NAME = "promotion_service"
	DELIVERY_SERVICE_NAME  = "delivery_service"
	PAYMENT_SERVICE_NAME   = "payment_service"
)

const NUMBER_OF_SERVICES_COMMIT = 3

const TTL_ORDER_MESSAGE = 10 * time.Minute
