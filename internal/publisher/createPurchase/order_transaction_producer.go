package createPurchase

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/domain/message"
	"latipe-transaction-service/pkgs/rabbitclient"

	"time"
)

type OrderOrchestratorPub struct {
	channel *amqp.Channel
	cfg     *config.Config
}

func NewOrderOrchestratorPub(cfg *config.Config, conn *amqp.Connection) *OrderOrchestratorPub {
	producer := OrderOrchestratorPub{
		cfg: cfg,
	}

	ch, err := conn.Channel()
	if err != nil {
		rabbitclient.FailOnError(err, "Failed to open a channel")
		return nil
	}
	producer.channel = ch

	return &producer
}

func (pub *OrderOrchestratorPub) ReplyPurchaseMessage(message *message.CreateOrderReplyMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEvent.ReplyRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEvent.ReplyRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) PublishPurchaseProductMessage(message *message.OrderProductMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderProductEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderProductEvent.CommitRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderProductEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderProductEvent.CommitRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) PublishPurchasePromotionMessage(message *message.ApplyVoucherMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.CommitRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.CommitRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) PublishPurchaseDeliveryMessage(message *message.DeliveryMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.CommitRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.CommitRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) PublishPurchasePaymentMessage(message *message.PaymentMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.CommitRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.CommitRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) PublishPurchaseEmailMessage(message *message.EmailPurchaseMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.CommitRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.CommitRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) PublishPurchaseCartMessage(message *message.CartMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderCartEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderCartEvent.CommitRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderCartEvent.Exchange,         // exchange
		pub.cfg.RabbitMQ.SagaOrderCartEvent.CommitRoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) RollbackPurchaseProductMessage(message *message.RollbackPurchaseMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderProductEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderProductEvent.RollbackRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderProductEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderProductEvent.RollbackRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) RollbackPurchasePromotionMessage(message *message.RollbackPurchaseMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.RollbackRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPromotionEvent.RollbackRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) RollbackPurchaseDeliveryMessage(message *message.RollbackPurchaseMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.RollbackRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderDeliveryEvent.RollbackRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) RollbackPurchaseEmailMessage(message *message.RollbackPurchaseMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.RollbackRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEmailEvent.RollbackRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}

func (pub *OrderOrchestratorPub) RollbackPurchasePaymentMessage(message *message.RollbackPurchaseMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.CommitRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderPaymentEvent.CommitRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}
