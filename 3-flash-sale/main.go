package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const userCount = 100_000

type Inventory struct {
	itemsInStock atomic.Int32
	itemsBought  atomic.Int32
	
}

func paymentGateRoll() bool {
	time.Sleep(1 * time.Second)
	roll := rand.Intn(100)
	if roll < 20 {
		return false
	}
	return true
}

func (inventory *Inventory) Purchase(userId int) (bool, error) {
	for {
		itemCount := inventory.itemsInStock.Load()
		if itemCount <= 0 {
			return false, nil
		}

		if inventory.itemsInStock.CompareAndSwap(itemCount, itemCount-1) {
			if paymentGateRoll() {
				fmt.Printf("User %d managed to buy an item №%d!\n", userId, itemCount)
				inventory.itemsBought.Add(1)
				return true, nil
			} else {
				inventory.itemsInStock.Add(1)
				fmt.Printf("Users %d payment failed!\n", userId)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (Inventory *Inventory) Reserve(userId, item int) {

}

func main() {

	var inventory Inventory
	inventory.itemsInStock.Store(100)
	inventory.itemsBought.Store(0)
	var wg sync.WaitGroup

	for user := range userCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inventory.Purchase(user)
		}()
	}
	wg.Wait()
	fmt.Println("Items bought: ", inventory.itemsBought.Load())

}
