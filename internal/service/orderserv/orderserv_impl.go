package orderserv

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"latipe-transaction-service/internal/domain/entities"
	"latipe-transaction-service/internal/domain/message"
	"latipe-transaction-service/internal/domain/repos"
	msgqueue "latipe-transaction-service/internal/publisher/createPurchase"
	"strings"
	"time"
)

type orderService struct {
	transactionRepo         repos.TransactionRepository
	transactionOrchestrator *msgqueue.OrderOrchestratorPub
}

func NewOrderService(transRepos repos.TransactionRepository, orchestrator *msgqueue.OrderOrchestratorPub) OrderService {
	return &orderService{
		transactionRepo:         transRepos,
		transactionOrchestrator: orchestrator,
	}
}

func (o orderService) StartPurchaseTransaction(ctx context.Context, msg *message.OrderPendingMessage) error {
	dao := entities.TransactionLog{
		OrderID:           msg.OrderID,
		TransactionStatus: 0,
		CreatedAt:         time.Now(),
		SuccessAt:         time.Now(),
	}

	txs := []entities.Commits{
		{
			ServiceName: DELIVERY_SERVICE_NAME,
			TxStatus:    0,
		},
		{
			ServiceName: PAYMENT_SERVICE_NAME,
			TxStatus:    0,
		},
		{
			ServiceName: PRODUCT_SERVICE_NAME,
			TxStatus:    0,
		}, {
			ServiceName: PROMOTION_SERVICE_NAME,
			TxStatus:    0,
		},
	}
	dao.Commits = txs

	if err := o.transactionRepo.CreateTransactionData(ctx, &dao); err != nil {
		return err
	}

	// Create channels for error handling
	errCh := make(chan error, 5) // Number of messages to send

	// Define a function to send messages and handle errors
	sendMessage := func(fn func() error) {
		if err := fn(); err != nil {
			errCh <- err
		}
	}

	// Start goroutines for each message type
	go sendMessage(func() error {
		productMsg := message.OrderProductMessage{}
		var items []message.ProductMessage
		for _, i := range msg.OrderItems {
			item := message.ProductMessage{
				ProductId: i.ProductItem.ProductID,
				OptionId:  i.ProductItem.OptionID,
				Quantity:  i.ProductItem.Quantity,
			}
			items = append(items, item)
		}
		productMsg.Items = items
		return o.transactionOrchestrator.PublishPurchaseProductMessage(&productMsg)
	})

	go sendMessage(func() error {
		voucherMsg := message.VoucherMessage{
			OrderID:  msg.OrderID,
			Vouchers: strings.Split(msg.Vouchers, ";"),
		}

		return o.transactionOrchestrator.PublishPurchasePromotionMessage(&voucherMsg)
	})

	go sendMessage(func() error {
		paymentMsg := message.PaymentMessage{
			OrderID:       msg.OrderID,
			PaymentMethod: msg.PaymentMethod,
			SubTotal:      msg.SubTotal,
			TotalDiscount: msg.ItemDiscount + msg.ShippingDiscount,
			Amount:        msg.Amount,
			UserId:        msg.UserRequest.UserId,
			CreatedAt:     time.Now(),
		}
		return o.transactionOrchestrator.PublishPurchasePaymentMessage(&paymentMsg)
	})

	go sendMessage(func() error {
		return o.transactionOrchestrator.PublishPurchaseEmailMessage(&message.EmailPurchaseMessage{
			EmailRecipient: msg.UserRequest.Username,
			Name:           msg.Address.Name,
			OrderId:        msg.OrderID,
		})
	})

	go sendMessage(func() error {
		deliveryMsg := message.DeliveryMessage{
			OrderID:       msg.OrderID,
			DeliveryID:    msg.Delivery.DeliveryId,
			PaymentMethod: msg.PaymentMethod,
			Total:         msg.Amount,
			ShippingCost:  msg.ShippingCost,
			ReceiveDate:   msg.Delivery.ReceivingDate,
		}
		deliveryMsg.Address.AddressId = msg.Address.AddressId
		deliveryMsg.Address.Name = msg.Address.Name
		deliveryMsg.Address.AddressDetail = msg.Address.AddressDetail

		deliveryMsg.Address.Phone = msg.Address.Phone
		return o.transactionOrchestrator.PublishPurchaseDeliveryMessage(&deliveryMsg)
	})

	// Wait for all goroutines to finish
	for i := 0; i < 5; i++ {
		select {
		case err := <-errCh:
			// Handle error (you can log it, return it, etc.)
			log.Errorf("Publish message failed: %v", err)
		}
	}

	return nil
}

func (o orderService) HandleTransactionPurchaseReply(ctx context.Context, msg *message.CreateOrderReplyMessage, serviceType int) error {
	trans, err := o.transactionRepo.FindByOrderID(ctx, msg.OrderID)
	if err != nil {
		return err
	}

	if trans.TransactionStatus == entities.TX_PENDING {
		if msg.Status == message.COMMIT_SUCCESS {
			trans = MappingTxStatus(trans, serviceType, entities.TX_SUCCESS)
			commit := entities.Commits{
				ServiceName: MappingServiceName(serviceType),
				TxStatus:    entities.TX_SUCCESS,
			}
			if err := o.transactionRepo.UpdateTransaction(ctx, trans, &commit); err != nil {
				return err
			}
		}

		if msg.Status == message.COMMIT_FAIL {
			trans = MappingTxStatus(trans, serviceType, entities.TX_FAIL)
			commit := entities.Commits{
				ServiceName: MappingServiceName(serviceType),
				TxStatus:    entities.TX_SUCCESS,
			}

			trans.TransactionStatus = entities.TX_FAIL

			if err := o.transactionRepo.UpdateTransaction(ctx, trans, &commit); err != nil {
				return err
			}
		}
	}

	if trans.TransactionStatus == entities.TX_FAIL {
		commit := entities.Commits{
			ServiceName: MappingServiceName(serviceType),
			TxStatus:    entities.TX_FAIL,
		}
		if err := o.rollbackPurchaseToService(&commit, trans.OrderID); err != nil {
			return err
		}
	}

	return nil
}

func (o orderService) RollbackTransactionPub(dao *entities.TransactionLog) error {
	for _, i := range dao.Commits {

		if i.TxStatus == entities.TX_SUCCESS {
			if err := o.rollbackPurchaseToService(&i, dao.OrderID); err != nil {
				log.Errorf("rollback was failed cause:%v", err)
			}
		}
	}
	return nil
}

func (o orderService) rollbackPurchaseToService(dao *entities.Commits, orderId string) error {

	msg := message.RollbackPurchaseMessage{
		Status:  entities.TX_FAIL,
		OrderID: orderId,
	}

	switch dao.ServiceName {
	case DELIVERY_SERVICE_NAME:
		if err := o.transactionOrchestrator.RollbackPurchaseDeliveryMessage(&msg); err != nil {
			return err
		}
	case PRODUCT_SERVICE_NAME:
		if err := o.transactionOrchestrator.RollbackPurchaseProductMessage(&msg); err != nil {
			return err
		}
	case PROMOTION_SERVICE_NAME:
		if err := o.transactionOrchestrator.RollbackPurchasePromotionMessage(&msg); err != nil {
			return err
		}
	case PAYMENT_SERVICE_NAME:
		if err := o.transactionOrchestrator.RollbackPurchasePaymentMessage(&msg); err != nil {
			return err
		}
	}

	return nil
}

func (o orderService) CheckTransactionStatus(ctx context.Context) error {
	txs, err := o.transactionRepo.FindAllPendingTransaction(ctx)
	if err != nil {
		return err
	}

	for _, i := range txs {
		if IsCommitSuccess(i.Commits) {
			replyMessage := message.CreateOrderReplyMessage{
				Status:  entities.TX_SUCCESS,
				OrderID: i.OrderID,
			}

			err := o.transactionOrchestrator.ReplyPurchaseMessage(&replyMessage)
			if err != nil {
				log.Errorf("publish tx order [%v] message was failed %v", i.OrderID, err)
			}
		}
	}

	return err
}
