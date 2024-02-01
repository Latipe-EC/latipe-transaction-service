package repos

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"latipe-transaction-service/internal/domain/entities"
	"latipe-transaction-service/pkgs/db/mongodb"
	"latipe-transaction-service/pkgs/util/pagable"
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

func (t *transactionRepository) UpdateTransactionCommit(ctx context.Context, dao *entities.TransactionLog, commit *entities.Commits) error {
	filter := bson.M{
		"order_id":             dao.OrderID,
		"commits.service_name": commit.ServiceName,
	}

	update := bson.D{
		{"$set", bson.D{
			{"commits.$.tx_status", commit.TxStatus},
			{"commits.$.updated_at", time.Now()},
		}},
	}

	_, err := t.transCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}

	return nil
}

func (t *transactionRepository) UpdateTransactionStatus(ctx context.Context, dao *entities.TransactionLog) error {
	filter := bson.M{"order_id": dao.OrderID}

	update := bson.D{
		{"$set", bson.D{
			{"transaction_status", dao.TransactionStatus},
			{"updated_at", time.Now()},
		}},
	}

	_, err := t.transCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}

	return nil
}

func (t *transactionRepository) FindAllPendingTransaction(ctx context.Context) ([]*entities.TransactionLog, error) {
	var txs []*entities.TransactionLog
	filter := bson.M{
		"created_at": bson.M{"$gt": time.Now()},
		"status":     entities.TX_PENDING,
	}

	cursor, err := t.transCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &txs); err != nil {
		return nil, err
	}

	return txs, nil
}

func (t *transactionRepository) CountTxSuccess(ctx context.Context, orderId string) (int, error) {

	filter := bson.D{{"order_id", orderId}}
	// Truy vấn MongoDB
	cursor, err := t.transCollection.Aggregate(ctx, mongo.Pipeline{
		{{"$match", filter}},
		{{"$project", bson.D{
			{"countCommitsSuccess", bson.D{
				{"$size", bson.D{
					{"$filter", bson.D{
						{"input", "$commits"},
						{"as", "commit"},
						{"cond", bson.D{{"$eq", []interface{}{"$$commit.tx_status", 1}}}},
					}},
				}},
			}},
		}}},
	})
	if err != nil {
		return -1, err
	}

	var result entities.ResultCursor
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			log.Error(err)
			return -1, err
		}

	}

	return int(result.CountCommitsSuccess), nil
}

func (t *transactionRepository) FindByTransactionId(ctx context.Context, transactionId string) (*entities.TransactionLog, error) {
	var entity entities.TransactionLog
	id, _ := primitive.ObjectIDFromHex(transactionId)

	err := t.transCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, err
}

func (t *transactionRepository) FindAll(ctx context.Context, query *pagable.Query) ([]entities.TransactionLog, int, error) {
	var trans []entities.TransactionLog

	filter, err := query.ConvertQueryToFilter()
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetOffset()))
	cursor, err := t.transCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(ctx, &trans); err != nil {
		return nil, 0, err
	}

	total, err := t.transCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return trans, int(total), err
}
