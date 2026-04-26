package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
)

const productCount = 1000_000
const validatorWorkerCount = 10

type Product struct {
	ID    string  `faker:"uuid_hyphenated"`
	Name  string  `faker:"username"`
	Price float64 `faker:"amount"`
}

func (product *Product) Enrich() error {
	err := faker.FakeData(product)
	if err != nil {
		return err
	}
	return nil
}

func EnrichedProducts() <-chan Product {
	products := make(chan Product, productCount)
	var wg sync.WaitGroup
	for range productCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := Product{}
			p.Enrich()
			products <- p
		}()
	}
	go func() {
		wg.Wait()
		close(products)
	}()
	return products
}

type ProcessedProduct struct {
	Product
	Score       int
	IsPromoted  bool
	ProcessedBy int
}

func validate(workerId int, product Product) ProcessedProduct {
	time.Sleep(1 * time.Second)
	score := rand.N(100) + 1
	var isPromoted bool
	if score > 65 {
		isPromoted = true
	} else {
		isPromoted = false
	}
	return ProcessedProduct{Product: product, Score: score, IsPromoted: isPromoted, ProcessedBy: workerId}

}

func TrustValidate(products <-chan Product) []<-chan ProcessedProduct {
	workerChannels := make([]chan ProcessedProduct, validatorWorkerCount)
	result := make([]<-chan ProcessedProduct, validatorWorkerCount)
	for i := range validatorWorkerCount {
		ch := make(chan ProcessedProduct, productCount)
		workerChannels[i] = ch
		result[i] = ch
	}

	for workerId := range validatorWorkerCount {
		go func(workerId int) {
			defer close(workerChannels[workerId])
			for product := range products {
				workerChannels[workerId] <- validate(workerId, product)
			}
		}(workerId)
	}

	return result
}

func MergeChannels(channels []<-chan ProcessedProduct) chan ProcessedProduct {
	result := make(chan ProcessedProduct)
	var wg sync.WaitGroup
	for _, c := range channels {
		wg.Add(1)
		go func(c <-chan ProcessedProduct) {
			defer wg.Done()
			for v := range c {
				result <- v
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(result)
	}()
	return result

}

func main() {

	products := EnrichedProducts()
	processedProducts := TrustValidate(products)
	merged := MergeChannels(processedProducts)
	fmt.Println("Processes products:")

	for range productCount {
		fmt.Println(<-merged)
	}

}
