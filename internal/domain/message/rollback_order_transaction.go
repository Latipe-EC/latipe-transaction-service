package message

type RollbackPurchaseMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}
