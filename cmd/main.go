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

	var wg sync.WaitGroup

	startSubscribers(serv, &wg)
	startCronJobs(serv, &wg)
	startAPIHandler(serv, &wg)

	wg.Wait()
}

func startSubscribers(serv *server.Server, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		serv.PurchaseCreateSub().ListenProductPurchaseCreate(wg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenProductPurchaseReply(wg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenPromotionPurchaseReply(wg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenPaymentPurchaseReply(wg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenDeliveryPurchaseReply(wg)
	}()
}

func startCronJobs(serv *server.Server, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := serv.CheckTxStatusCron().StartJob(wg)
		if err != nil {
			log.Error(err)
		}
	}()
}

func startAPIHandler(serv *server.Server, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
			fmt.Printf("%s", err)
		}
	}()
}
