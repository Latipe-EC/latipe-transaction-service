package createPurchase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/domain/message"
	"latipe-transaction-service/internal/service/orderserv"
	"sync"
	"time"
)

type PurchaseCreateOrchestratorSubscriber struct {
	config    *config.Config
	orderServ orderserv.OrderService
	conn      *amqp.Connection
}

func NewPurchaseCreateOrchestratorSubscriber(cfg *config.Config,
	orderServ orderserv.OrderService, conn *amqp.Connection) *PurchaseCreateOrchestratorSubscriber {
	return &PurchaseCreateOrchestratorSubscriber{
		config:    cfg,
		orderServ: orderServ,
		conn:      conn,
	}
}

func (orch PurchaseCreateOrchestratorSubscriber) ListenPurchaseCreate(wg *sync.WaitGroup) {
	channel, err := orch.conn.Channel()
	defer channel.Close()

	// define an exchange type "topic"
	err = channel.ExchangeDeclare(
		orch.config.RabbitMQ.SagaOrderEvent.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare exchange: %v", err)
	}

	// create queue
	q, err := channel.QueueDeclare(
		"purchase_event",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare queue: %v", err)
	}

	err = channel.QueueBind(
		q.Name,
		orch.config.RabbitMQ.SagaOrderEvent.PublishRoutingKey,
		orch.config.RabbitMQ.SagaOrderEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                           // queue
		orch.config.RabbitMQ.ServiceName, // consumer
		true,                             // auto ack
		false,                            // exclusive
		false,                            // no local
		false,                            // no wait
		nil,                              //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := orch.handleMessage(&msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}
	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")
}

func (orch PurchaseCreateOrchestratorSubscriber) handleMessage(msg *amqp.Delivery) error {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageDTO := message.OrderPendingMessage{}

	if err := json.Unmarshal(msg.Body, &messageDTO); err != nil {
		log.Infof("Parse message to order failed cause: %s", err)
		return err
	}

	err := orch.orderServ.StartPurchaseTransaction(ctx, &messageDTO)
	if err != nil {
		log.Infof("Handling reply message was failed cause: %s", err)
		return err
	}

	endTime := time.Now()
	log.Infof("The orders [checkout_id: %v]  was processed successfully - duration:%v", messageDTO.CheckoutMessage.CheckoutID, endTime.Sub(startTime))
	fmt.Println()
	return nil
}
