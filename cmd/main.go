package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	server "latipe-transaction-service/internal"
	"sync"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[Info] Application has closed")

	serv, err := server.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	//subscriber
	var wg sync.WaitGroup
	//publish transaction
	{
		wg.Add(1)
		go serv.PurchaseCreateSub().ListenProductPurchaseCreate(&wg)
	}
	//waiting reply
	{
		wg.Add(1)
		go serv.PurchaseReplySub().ListenProductPurchaseReply(&wg)

		wg.Add(1)
		go serv.PurchaseReplySub().ListenPromotionPurchaseReply(&wg)

		wg.Add(1)
		go serv.PurchaseReplySub().ListenPaymentPurchaseReply(&wg)

		wg.Add(1)
		go serv.PurchaseReplySub().ListenDeliveryPurchaseReply(&wg)
	}

	//api handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
			fmt.Printf("%s", err)
		}
	}()

	wg.Wait()
}
