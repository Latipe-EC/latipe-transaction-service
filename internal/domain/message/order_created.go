package message

import "time"

type OrderPendingMessage struct {
	UserRequest     UserRequest     `json:"user_request,omitempty"`
	Address         OrderAddress    `json:"address,omitempty" validate:"required"`
	CheckoutMessage CheckoutMessage `json:"checkout_data"`
	OrderDetail     []OrderDetail   `json:"order_detail"`
}
type OrderDetail struct {
	Status           int                 `json:"status"`
	OrderID          string              `json:"order_id,omitempty"`
	StoreID          string              `json:"store_id" validate:"required"`
	Amount           int                 `json:"amount,omitempty" validate:"required"`
	ShippingCost     int                 `json:"shipping_cost,omitempty"`
	ShippingDiscount int                 `json:"shipping_discount,omitempty" validate:"required"`
	ItemDiscount     int                 `json:"item_discount,omitempty" validate:"required"`
	SubTotal         int                 `json:"sub_total,omitempty" validate:"required"`
	PaymentMethod    int                 `json:"payment_method,omitempty" validate:"required"`
	Vouchers         string              `json:"vouchers,omitempty"`
	Delivery         Delivery            `json:"delivery,omitempty" validate:"required"`
	OrderItems       []OrderItemsMessage `json:"order_items,omitempty" validate:"required"`
	CartIds          []string            `json:"cart_ids"`
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
	AddressId     string `json:"address_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	AddressDetail string `json:"address_detail"`
}

type Delivery struct {
	DeliveryId    string    `json:"delivery_id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Cost          int       `json:"cost" validate:"required"`
	ReceivingDate time.Time `json:"receiving_date" validate:"required"`
}

type PromotionMessage struct {
	FreeShippingVoucher string   `json:"free_shipping_voucher"`
	PaymentVoucher      string   `json:"payment_voucher"`
	ShopVoucher         []string `json:"shop_vouchers"`
}

type CheckoutMessage struct {
	CheckoutID    string      `json:"checkout_id"`
	UserID        string      `json:"user_id"`
	OrderData     []OrderData `json:"order_data"`
	TotalAmount   uint        `json:"total_amount"`
	PaymentMethod int         `json:"payment_method"`
}

type OrderData struct {
	OrderID string `json:"order_id"`
	Amount  uint   `json:"amount"`
}
