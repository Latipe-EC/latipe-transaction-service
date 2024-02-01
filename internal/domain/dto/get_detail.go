package dto

import (
	"time"
)

type GetTransactionLogByIdRequest struct {
	AuthorizationHeader
	TransactionID string `params:"id"`
}

type GetTransactionLogByIdResponse struct {
	ID                string            `json:"id,omitempty"`
	OrderID           string            `json:"order_id"`
	TransactionStatus int               `json:"transaction_status"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	Commits           []CommitsResponse `json:"commits"`
}

type CommitsResponse struct {
	ServiceName string    `json:"service_name"`
	TxStatus    int       `json:"tx_status"`
	UpdatedAt   time.Time `json:"updated_at"`
}
