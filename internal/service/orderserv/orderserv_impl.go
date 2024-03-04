package orderserv

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"latipe-transaction-service/internal/domain/entities"
	"latipe-transaction-service/internal/domain/message"
	"latipe-transaction-service/internal/domain/repos"
	msgqueue "latipe-transaction-service/internal/publisher/createPurchase"
	"latipe-transaction-service/internal/service/notifyserv"
	redisCache "latipe-transaction-service/pkgs/cache/redis"
	"strings"
	"sync"
	"time"
)

type orderService struct {
	transactionRepo         repos.TransactionRepository
	notifyServ              notifyserv.NotifyService
	transactionOrchestrator *msgqueue.OrderOrchestratorPub
	cacheEngine             *redisCache.CacheEngine
}

func NewOrderService(transRepos repos.TransactionRepository, orchestrator *msgqueue.OrderOrchestratorPub,
	notifyService notifyserv.NotifyService,
	cacheEngine *redisCache.CacheEngine) OrderService {
	return &orderService{
		transactionRepo:         transRepos,
		transactionOrchestrator: orchestrator,
		notifyServ:              notifyService,
		cacheEngine:             cacheEngine,
	}
}

func (o orderService) StartPurchaseTransaction(ctx context.Context, msg *message.OrderPendingMessage) error {
	fmt.Println()
	log.Infof("starting purchase transaction for order [%v]", msg.CheckoutMessage.CheckoutID)

	for _, i := range msg.OrderDetail {
		dao := entities.TransactionLog{
			OrderID:           i.OrderID,
			TransactionStatus: 0,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
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
			},
		}

		if i.Vouchers != "" {
			txs = append(txs, entities.Commits{
				ServiceName: PROMOTION_SERVICE_NAME,
				TxStatus:    0,
			})
		}

		dao.Commits = txs

		if err := o.transactionRepo.CreateTransactionData(ctx, &dao); err != nil {
			return err
		}

		if err := o.handlePurchaseTransaction(msg, &i); err != nil {
			continue
		}

	}
	log.Infof("purchase transaction is finished [%v]", msg.CheckoutMessage.CheckoutID)

	return nil
}

func (o orderService) handlePurchaseTransaction(msg *message.OrderPendingMessage, orderDetail *message.OrderDetail) error {

	// Create channels for error handling
	errCh := make(chan error, 5) // number of messages to send

	// create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(7) // number of goroutines is equal to the number of message types

	// define a function to send messages and handle errors
	handlerMessageGoroutine := func(fn func() error) {
		defer wg.Done() // decrement the wait group counter when the goroutine finishes
		if err := fn(); err != nil {
			errCh <- err
		}
	}

	// start goroutines for each message type and set cache message
	go handlerMessageGoroutine(func() error {
		return o.cacheEngine.Set(orderDetail.OrderID, msg, TTL_ORDER_MESSAGE)
	})

	//publish product message
	go handlerMessageGoroutine(func() error {
		// Create and send product message
		productMsg := message.OrderProductMessage{
			OrderId: orderDetail.OrderID,
			StoreId: orderDetail.StoreID,
		}

		var items []message.ProductMessage
		for _, i := range orderDetail.OrderItems {
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

	//publish payment message
	go handlerMessageGoroutine(func() error {
		// Create and send payment message
		paymentMsg := message.PaymentMessage{
			OrderID:       orderDetail.OrderID,
			PaymentMethod: orderDetail.PaymentMethod,
			SubTotal:      orderDetail.SubTotal,
			Amount:        orderDetail.Amount,
			UserId:        msg.UserRequest.UserId,
			Status:        orderDetail.Status,
		}
		return o.transactionOrchestrator.PublishPurchasePaymentMessage(&paymentMsg)
	})

	//publish voucher message
	go handlerMessageGoroutine(func() error {
		codes := strings.Split(orderDetail.Vouchers, "-")
		if codes[0] == "" {
			codes = nil
		}

		promotionReq := message.ApplyVoucherMessage{
			CheckoutID:   msg.CheckoutMessage.CheckoutID,
			UserID:       msg.UserRequest.UserId,
			OrderID:      orderDetail.OrderID,
			VoucherCodes: codes,
		}

		return o.transactionOrchestrator.PublishPurchasePromotionMessage(&promotionReq)

	})

	//publish email message
	go handlerMessageGoroutine(func() error {
		// Create and send email message
		return o.transactionOrchestrator.PublishPurchaseEmailMessage(&message.EmailPurchaseMessage{
			EmailRecipient: msg.UserRequest.Username,
			Name:           msg.Address.Name,
			OrderId:        orderDetail.OrderID,
		})
	})

	//publish delivery message
	go handlerMessageGoroutine(func() error {
		// Create and send delivery message
		deliveryMsg := message.DeliveryMessage{
			OrderID:       orderDetail.OrderID,
			DeliveryID:    orderDetail.Delivery.DeliveryId,
			PaymentMethod: orderDetail.PaymentMethod,
			Total:         orderDetail.Amount,
			ShippingCost:  orderDetail.ShippingCost,
			ReceiveDate:   orderDetail.Delivery.ReceivingDate,
		}
		deliveryMsg.Address.AddressId = msg.Address.AddressId
		deliveryMsg.Address.Name = msg.Address.Name
		deliveryMsg.Address.AddressDetail = msg.Address.AddressDetail
		deliveryMsg.Address.Phone = msg.Address.Phone
		return o.transactionOrchestrator.PublishPurchaseDeliveryMessage(&deliveryMsg)
	})

	//publish cart id message
	go handlerMessageGoroutine(func() error {
		if len(orderDetail.CartIds) > 0 {
			// Create and send delivery message
			deliveryMsg := message.CartMessage{
				CartIdVmList: orderDetail.CartIds,
			}

			return o.transactionOrchestrator.PublishPurchaseCartMessage(&deliveryMsg)
		}
		return nil
	})

	// wait for all goroutines to finish
	wg.Wait()

	// close the error channel after all goroutines have finished
	close(errCh)

	// handle any remaining errors in the error channel
	for err := range errCh {
		log.Errorf("Publish message failed: %v", err)
	}
	log.Infof("purchase transaction is finished [%v]", orderDetail.OrderID)
	return nil
}

func (o orderService) HandleTransactionPurchaseReply(ctx context.Context, msg *message.CreateOrderReplyMessage, serviceType int) error {
	trans, err := o.transactionRepo.FindByOrderID(ctx, msg.OrderID)
	if err != nil {
		return err
	}

	// Update transaction status if the transaction is pending
	// and a reply message has a status of success.

	// Rollback the transaction where the service has sent message has a status of success
	// and transaction status is failed
	if msg.Status == message.COMMIT_SUCCESS {

		switch trans.TransactionStatus {
		case entities.TX_PENDING:
			trans = MappingTxStatus(trans, serviceType, entities.TX_SUCCESS)
			commit := entities.Commits{
				ServiceName: MappingServiceName(serviceType),
				TxStatus:    entities.TX_SUCCESS,
			}
			if err := o.transactionRepo.UpdateTransactionCommit(ctx, trans, &commit); err != nil {
				return err
			}

			count, err := o.transactionRepo.CountTxSuccess(ctx, msg.OrderID)
			if err != nil {
				log.Error(err)
			}

			// Update order transaction status and reply purchase message to order
			// then delete the cache message
			if count == len(trans.Commits)-1 {
				errChan := make(chan error, 3)

				var wg sync.WaitGroup
				wg.Add(3)
				processPool := func(fn func() error) {
					defer wg.Done()
					if err := fn(); err != nil {
						errChan <- err
					}
				}

				go processPool(func() error {
					replyMsg := message.CreateOrderReplyMessage{
						Status:  entities.TX_SUCCESS,
						OrderID: msg.OrderID,
					}
					return o.transactionOrchestrator.ReplyPurchaseMessage(&replyMsg)
				})

				go processPool(func() error {
					trans.TransactionStatus = entities.TX_SUCCESS
					return o.transactionRepo.UpdateTransactionStatus(ctx, trans)
				})

				go processPool(func() error {
					return o.cacheEngine.Delete(msg.OrderID)
				})

				wg.Wait()
				close(errChan)

				// handle any remaining errors in the error channel
				for err := range errChan {
					log.Errorf("committing transaction is failed: %v", err)
				}
				log.Infof(" purchase transaction is committed [%v]", msg.OrderID)

				err := o.notifyServ.SendMessageToTelegram(trans, "success")
				if err != nil {
					log.Error(err)
				}

				return nil
			}

		case entities.TX_FAIL:
			commit := entities.Commits{
				ServiceName: MappingServiceName(serviceType),
				TxStatus:    entities.TX_FAIL,
			}

			if err := o.rollbackPurchaseToService(&commit, trans.OrderID); err != nil {
				return err
			}
		}

	}

	// rollback transaction if reply message has a status of failed
	if msg.Status == message.COMMIT_FAIL && trans.TransactionStatus == entities.TX_PENDING {
		trans = MappingTxStatus(trans, serviceType, entities.TX_FAIL)
		commit := entities.Commits{
			ServiceName: MappingServiceName(serviceType),
			TxStatus:    entities.TX_FAIL,
		}
		if err := o.transactionRepo.UpdateTransactionCommit(ctx, trans, &commit); err != nil {
			return err
		}

		trans.TransactionStatus = entities.TX_FAIL

		reason := fmt.Sprintf("Service failed:%s", MappingServiceName(serviceType))
		err := o.notifyServ.SendMessageToTelegram(trans, reason)
		if err != nil {
			log.Error(err)
		}

		if err := o.transactionRepo.UpdateTransactionStatus(ctx, trans); err != nil {
			return err
		}

		if err := o.RollbackTransactionPub(trans); err != nil {
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

func (o orderService) CancelOrder(ctx context.Context, req *message.CancelOrderMessage) error {
	msg := message.RollbackPurchaseMessage{
		Status:  req.CancelStatus,
		OrderID: req.OrderID,
	}

	// Create channels for error handling
	errCh := make(chan error, 4) // number of messages to send

	// create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(4) // number of goroutines is equal to the number of message types

	// define a function to send messages and handle errors
	handlerMessageGoroutine := func(fn func() error) {
		defer wg.Done() // decrement the wait group counter when the goroutine finishes
		if err := fn(); err != nil {
			errCh <- err
		}
	}

	// start goroutines for each message type and set cache message
	go handlerMessageGoroutine(func() error {
		return o.transactionOrchestrator.RollbackPurchaseDeliveryMessage(&msg)
	})

	go handlerMessageGoroutine(func() error {
		return o.transactionOrchestrator.RollbackPurchaseProductMessage(&msg)
	})

	go handlerMessageGoroutine(func() error {
		return o.transactionOrchestrator.RollbackPurchasePromotionMessage(&msg)
	})

	go handlerMessageGoroutine(func() error {
		return o.transactionOrchestrator.RollbackPurchasePaymentMessage(&msg)
	})

	// wait for all goroutines to finish
	wg.Wait()

	// close the error channel after all goroutines have finished
	close(errCh)

	// handle any remaining errors in the error channel
	for err := range errCh {
		log.Errorf("Publish message failed: %v", err)
	}

	return nil
}
