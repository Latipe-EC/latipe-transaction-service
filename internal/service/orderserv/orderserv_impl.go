package orderserv

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"latipe-transaction-service/internal/domain/entities"
	"latipe-transaction-service/internal/domain/message"
	"latipe-transaction-service/internal/domain/repos"
	msgqueue "latipe-transaction-service/internal/publisher/createPurchase"
	"strings"
	"sync"
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
	fmt.Println()
	log.Infof("starting purchase transaction for order [%v]", msg.OrderID)

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
	errCh := make(chan error, 5) // number of messages to send

	// create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(5) // number of goroutines is equal to the number of message types

	// define a function to send messages and handle errors
	sendMessage := func(fn func() error) {
		defer wg.Done() // decrement the wait group counter when the goroutine finishes
		if err := fn(); err != nil {
			errCh <- err
		}
	}

	// start goroutines for each message type

	go sendMessage(func() error {
		// Create and send product message
		productMsg := message.OrderProductMessage{
			OrderId: msg.OrderID,
			StoreId: msg.StoreId,
		}

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
		// Create and send voucher message
		voucherMsg := message.VoucherMessage{
			OrderID:  msg.OrderID,
			Vouchers: strings.Split(msg.Vouchers, ";"),
		}
		return o.transactionOrchestrator.PublishPurchasePromotionMessage(&voucherMsg)
	})

	go sendMessage(func() error {
		// Create and send payment message
		paymentMsg := message.PaymentMessage{
			OrderID:       msg.OrderID,
			PaymentMethod: msg.PaymentMethod,
			SubTotal:      msg.SubTotal,
			Amount:        msg.Amount,
			UserId:        msg.UserRequest.UserId,
			Status:        msg.Status,
		}
		return o.transactionOrchestrator.PublishPurchasePaymentMessage(&paymentMsg)
	})

	go sendMessage(func() error {
		// Create and send email message
		return o.transactionOrchestrator.PublishPurchaseEmailMessage(&message.EmailPurchaseMessage{
			EmailRecipient: msg.UserRequest.Username,
			Name:           msg.Address.Name,
			OrderId:        msg.OrderID,
		})
	})

	go sendMessage(func() error {
		// Create and send delivery message
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

	// wait for all goroutines to finish
	wg.Wait()

	// close the error channel after all goroutines have finished
	close(errCh)

	// handle any remaining errors in the error channel
	for err := range errCh {
		log.Errorf("Publish message failed: %v", err)
	}
	log.Infof("purchase transaction is finished [%v]", msg.OrderID)
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
