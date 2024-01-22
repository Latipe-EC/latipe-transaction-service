package message

const (
	COMMIT_SUCCESS = 1
	COMMIT_FAIL    = 0
)

type CreateOrderReplyMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}
