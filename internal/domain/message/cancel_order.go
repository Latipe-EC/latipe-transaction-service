package message

type CancelOrderMessage struct {
	OrderID      string `json:"order_id"`
	CancelStatus int    `json:"cancel_status"`
}
