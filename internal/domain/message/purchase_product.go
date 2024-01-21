package message

type OrderProductMessage struct {
	Items []ProductMessage `json:"items"`
}

type ProductMessage struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}
