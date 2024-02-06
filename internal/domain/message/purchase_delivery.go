package message

import "time"

type DeliveryMessage struct {
	OrderID       string          `json:"order_id"`
	DeliveryID    string          `json:"delivery_id"`
	Address       Address         `json:"address"`
	PaymentMethod int             `json:"payment_method"`
	Total         int             `json:"total"`
	ShippingCost  int             `json:"shipping_cost"`
	ReceiveDate   time.Time       `json:"receive_date"`
	ShippingItems []ShippingItems `json:"items"`
}

type Address struct {
	AddressId     string `json:"address_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	AddressDetail string `json:"address_detail"`
}

type ShippingItems struct {
	ProductName  string `json:"product_name"`
	ProductID    string `json:"product_id"`
	ProductImage string `json:"product_image"`
	Quantity     string `json:"quantity"`
}
