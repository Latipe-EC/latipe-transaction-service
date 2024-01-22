package message

import "time"

type PaymentMessage struct {
	OrderID       string    `json:"order_id"`
	PaymentMethod int       `json:"payment_method"`
	SubTotal      int       `json:"sub_total"`
	TotalDiscount int       `json:"total_discount"`
	Amount        int       `json:"amount"`
	UserId        string    `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}
