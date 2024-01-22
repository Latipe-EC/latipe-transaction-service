package message

type OrderPendingMessage struct {
	UserRequest      UserRequest         `json:"user_request"`
	Status           int                 `json:"status"`
	OrderID          string              `json:"order_id"`
	Amount           int                 `json:"amount"`
	ShippingCost     int                 `json:"shipping_cost"`
	ShippingDiscount int                 `json:"shipping_discount"`
	ItemDiscount     int                 `json:"item_discount"`
	SubTotal         int                 `json:"sub_total"`
	PaymentMethod    int                 `json:"payment_method" `
	Vouchers         string              `json:"vouchers,omitempty"`
	Address          OrderAddress        `json:"address,omitempty" `
	Delivery         Delivery            `json:"delivery,omitempty" `
	OrderItems       []OrderItemsMessage `json:"order_items,omitempty"`
}

type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type OrderItemsMessage struct {
	CartId      string      `json:"cart_id"`
	ProductItem ProductItem `json:"product_item"`
}

type ProductItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	NameOption  string `json:"name_option"`
	StoreID     string `json:"store_id"`
	OptionID    string `json:"option_id" `
	Image       string `json:"image"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	NetPrice    int    `json:"net_price"`
}

type OrderAddress struct {
	AddressId string `json:"address_id"`
}

type Delivery struct {
	DeliveryId    string `json:"delivery_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Cost          int    `json:"cost" validate:"required"`
	ReceivingDate string `json:"receiving_date" validate:"required"`
}
