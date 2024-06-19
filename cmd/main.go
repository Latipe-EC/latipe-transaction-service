package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	server "latipe-transaction-service/internal"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("Init application")
	// Get the number of CPU cores
	numCPU := runtime.NumCPU()
	fmt.Printf("Number of CPU cores: %d\n", numCPU)

	serv, err := server.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	var wg sync.WaitGroup

	startSubscribers(serv, &wg)
	startAPIHandler(serv, &wg)

	wg.Wait()
	log.Fatal("Application has closed")
}

func startSubscribers(serv *server.Server, wg *sync.WaitGroup) {
	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.PurchaseCreateSub().ListenPurchaseCreate(wg)
	})

	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.PurchaseCancelSub().ListenPurchaseCancel(wg)
	})

	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenProductPurchaseReply(wg)
	})

	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenPromotionPurchaseReply(wg)
	})

	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenPaymentPurchaseReply(wg)
	})

	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.PurchaseReplySub().ListenDeliveryPurchaseReply(wg)
	})
}

func startAPIHandler(serv *server.Server, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.App().Listen(serv.Config().Server.RestPort); err != nil {
			fmt.Printf("%s", err)
		}
	}()
}

func runWithRecovery(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Debugf("Recovered from panic: %v", r)
			}
		}()
		fn()
	}()
}
