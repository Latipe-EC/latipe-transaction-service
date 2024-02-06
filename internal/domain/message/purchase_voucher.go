package message

type ApplyVoucherMessage struct {
	UserId       string          `json:"user_id"`
	CheckoutData CheckoutRequest `json:"checkout_data"`
	Vouchers     []string        `json:"vouchers"`
}

type CheckoutRequest struct {
	CheckoutID string               `json:"checkout_id"`
	OrderData  []PromotionOrderData `json:"order_data"`
}

type PromotionOrderData struct {
	OrderID string `json:"order_id"`
	StoreID string `json:"store_id"`
}
