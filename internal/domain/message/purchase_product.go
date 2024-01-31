package message

type OrderProductMessage struct {
	OrderId string           `json:"orderId"`
	StoreId string           `json:"storeId"`
	Items   []ProductMessage `json:"items"`
}

type ProductMessage struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}
