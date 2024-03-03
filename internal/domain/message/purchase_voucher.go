package message

type ApplyVoucherMessage struct {
	CheckoutID   string   `json:"checkout_id"`
	UserID       string   `json:"user_id"`
	OrderID      string   `json:"order_id"`
	VoucherCodes []string `json:"voucher_codes"`
}
