package message

type RollbackOrderProductMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}

type RollbackOrderPromotionMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}

type RollbackOrderDeliveryMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}
