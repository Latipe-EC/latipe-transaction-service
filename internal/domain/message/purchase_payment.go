package message

type PaymentMessage struct {
	UserId        string `json:"user_id"`
	OrderID       string `json:"order_id"`
	PaymentMethod int    `json:"payment_method"`
	SubTotal      int    `json:"sub_total"`
	Amount        int    `json:"amount"`
	Status        int    `json:"status"`
}
