package message

type DeliveryMessage struct {
	OrderID       string `json:"order_id"`
	Name          string `json:"name"`
	Cost          int    `json:"cost"`
	ReceivingDate string `json:"receiving_date"`
	AddressId     string `json:"address_id"`
}
