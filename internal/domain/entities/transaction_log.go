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
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderID           string             `json:"order_id" bson:"order_id"`
	TransactionStatus int                `json:"transaction_status" bson:"transaction_status"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	SuccessAt         time.Time          `json:"success_at" bson:"success_at"`
	ProductCommit     int                `json:"product_commit" bson:"product_commit"`
	PromotionCommit   int                `json:"promotion_commit" bson:"promotion_commit"`
	DeliveryCommit    int                `json:"delivery_commit" bson:"delivery_commit"`
	PaymentCommit     int                `json:"payment_commit" bson:"payment_commit"`
	EmailCommit       int                `json:"email_commit" bson:"email_commit"`
}
