package message

type VoucherMessage struct {
	OrderID  string   `json:"order_id"`
	Vouchers []string `json:"vouchers"`
}
