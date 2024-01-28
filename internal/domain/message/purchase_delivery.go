package message

type DeliveryMessage struct {
	OrderID       string          `json:"order_id"`
	DeliveryID    string          `json:"delivery_id"`
	Address       Address         `json:"address"`
	PaymentMethod int             `json:"payment_method"`
	Total         int             `json:"total"`
	ShippingCost  int             `json:"shipping_cost"`
	ReceiveDate   string          `json:"receive_date"`
	ShippingItems []ShippingItems `json:"items"`
}

type Address struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	AddressDetail string `json:"address_detail"`
	ProvinceCode  string `json:"province_code"`
	DistrictCode  string `json:"district_code"`
	WardCode      string `json:"ward_code"`
}

type ShippingItems struct {
	ProductName  string `json:"product_name"`
	ProductID    string `json:"product_id"`
	ProductImage string `json:"product_image"`
	Quantity     string `json:"quantity"`
}
