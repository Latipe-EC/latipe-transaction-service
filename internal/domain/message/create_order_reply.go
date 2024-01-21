package message

const (
	COMMIT_SUCCESS = 1
	COMMIT_FAIL    = 2
)

type CreateOrderReplyMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}
