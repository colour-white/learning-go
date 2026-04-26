package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const userCount = 100_000
const maxItemCount = 100

type Inventory struct {
	itemsInStock chan int
	itemsBought  atomic.Int32
	done         chan struct{}
	closeOnce    sync.Once
}

func (inv *Inventory) Finish() {
	inv.closeOnce.Do(func() { close(inv.done) })
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
	select {

	case itemId := <-inventory.itemsInStock:

		if paymentGateRoll() {
			if inventory.itemsBought.Add(1) == maxItemCount {
				inventory.Finish()
			}
			fmt.Printf("User %d managed to buy an item №%d!\n", userId, itemId)

			return true, nil
		}
		select {
		case inventory.itemsInStock <- itemId:
			time.Sleep(1 * time.Second)
		case <-inventory.done:
		}
		fmt.Printf("User %d payment failed!\n", userId)
		return false, nil
	case <-inventory.done:
		return false, nil
	}
}

func main() {

	inventory := Inventory{itemsInStock: make(chan int, maxItemCount), done: make(chan struct{})}

	for itemId := range maxItemCount {
		inventory.itemsInStock <- itemId
	}
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
