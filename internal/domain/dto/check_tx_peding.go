package dto

type TransactionPendingResponse struct {
	Count    uint64   `json:"count"`
	OrderIds []string `json:"order_ids"`
}
