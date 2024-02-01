package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	TX_SUCCESS = 1
	TX_FAIL    = -1
	TX_PENDING = 0
)

type TransactionLog struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OrderID           string             `json:"order_id" bson:"order_id"`
	TransactionStatus int                `json:"transaction_status" bson:"transaction_status"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
	Commits           []Commits          `json:"commits" bson:"commits"`
}

type Commits struct {
	ServiceName string    `json:"service_name" bson:"service_name"`
	TxStatus    int       `json:"tx_status" bson:"tx_status"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
