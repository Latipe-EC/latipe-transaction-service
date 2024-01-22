package repos

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"latipe-transaction-service/internal/domain/entities"
	"latipe-transaction-service/pkgs/db/mongodb"
	"time"
)

type transactionRepository struct {
	transCollection *mongo.Collection
}

func NewTransactionRepository(client *mongodb.MongoClient) TransactionRepository {
	col := client.GetDB().Collection("transaction_logs")
	return &transactionRepository{transCollection: col}
}

func (t *transactionRepository) CreateTransactionData(ctx context.Context, dao *entities.TransactionLog) error {
	dao.CreatedAt = time.Now()

	_, err := t.transCollection.InsertOne(ctx, dao)
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionRepository) FindByOrderID(ctx context.Context, orderID string) (*entities.TransactionLog, error) {
	var entity entities.TransactionLog
	filter := bson.M{"order_id": orderID}

	err := t.transCollection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, err
}

func (t *transactionRepository) UpdateTransaction(ctx context.Context, dao *entities.TransactionLog) error {
	//TODO implement me
	panic("implement me")
}
