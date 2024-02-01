package createPurchase

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/domain/message"
	"latipe-transaction-service/internal/service/orderserv"
	"sync"
	"time"
)

type PurchaseReplySubscriber struct {
	config    *config.Config
	orderServ orderserv.OrderService
	conn      *amqp.Connection
}

func NewPurchaseSubscriberReply(cfg *config.Config, orderServ orderserv.OrderService,
	conn *amqp.Connection) *PurchaseReplySubscriber {
	return &PurchaseReplySubscriber{
		config:    cfg,
		orderServ: orderServ,
		conn:      conn,
	}
}

func (mq PurchaseReplySubscriber) ListenProductPurchaseReply(wg *sync.WaitGroup) {
	channel, err := mq.conn.Channel()
	defer channel.Close()

	// define an exchange type "topic"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
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
		"product_reply",
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
		mq.config.RabbitMQ.SagaOrderProductEvent.ReplyRoutingKey,
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                         // queue
		mq.config.RabbitMQ.ServiceName, // consumer
		true,                           // auto ack
		false,                          // exclusive
		false,                          // no local
		false,                          // no wait
		nil,                            //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := mq.replyHandler(msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}
	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")
}

func (mq PurchaseReplySubscriber) ListenPromotionPurchaseReply(wg *sync.WaitGroup) {
	channel, err := mq.conn.Channel()
	defer channel.Close()

	// define an exchange type "topic"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
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
		"promotion_reply",
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
		mq.config.RabbitMQ.SagaOrderPromotionEvent.ReplyRoutingKey,
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                         // queue
		mq.config.RabbitMQ.ServiceName, // consumer
		true,                           // auto ack
		false,                          // exclusive
		false,                          // no local
		false,                          // no wait
		nil,                            //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := mq.replyHandler(msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}
	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")
}

func (mq PurchaseReplySubscriber) ListenDeliveryPurchaseReply(wg *sync.WaitGroup) {
	channel, err := mq.conn.Channel()
	defer channel.Close()

	// define an exchange type "topic"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
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
		"delivery_reply",
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
		mq.config.RabbitMQ.SagaOrderDeliveryEvent.ReplyRoutingKey,
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                         // queue
		mq.config.RabbitMQ.ServiceName, // consumer
		true,                           // auto ack
		false,                          // exclusive
		false,                          // no local
		false,                          // no wait
		nil,                            //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := mq.replyHandler(msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}
	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")
}

func (mq PurchaseReplySubscriber) ListenPaymentPurchaseReply(wg *sync.WaitGroup) {
	channel, err := mq.conn.Channel()
	defer channel.Close()

	// define an exchange type "topic"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
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
		"payment_reply",
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
		mq.config.RabbitMQ.SagaOrderPaymentEvent.ReplyRoutingKey,
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                         // queue
		mq.config.RabbitMQ.ServiceName, // consumer
		true,                           // auto ack
		false,                          // exclusive
		false,                          // no local
		false,                          // no wait
		nil,                            //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := mq.replyHandler(msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}
	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")
}

func (mq PurchaseReplySubscriber) replyHandler(msg amqp.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageDTO := message.CreateOrderReplyMessage{}

	if err := json.Unmarshal(msg.Body, &messageDTO); err != nil {
		log.Infof("Parse message to order failed cause: %s", err)
		return err
	}
	log.Infof(" order message from: %s - status %v", msg.RoutingKey, messageDTO.Status)
	err := mq.orderServ.HandleTransactionPurchaseReply(ctx, &messageDTO,
		MappingRoutingKeyToService(msg.RoutingKey, mq.config))
	if err != nil {
		log.Infof("Handling reply message was failed cause: %s", err)
		return err
	}

	return nil
}
