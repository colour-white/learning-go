package main

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"math/rand/v2"
	"sync"
	"time"
)

const productCount = 1_000_000
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

func EnrichedProducts(ctx context.Context) <-chan Product {
	products := make(chan Product, productCount)
	var wg sync.WaitGroup
	for range productCount {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			p := Product{}
			p.Enrich()
			select {
			case <-ctx.Done():
				return
			case products <- p:
			}
		}(ctx)
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

func TrustValidate(ctx context.Context, products <-chan Product) []<-chan ProcessedProduct {
	workerChannels := make([]chan ProcessedProduct, validatorWorkerCount)
	result := make([]<-chan ProcessedProduct, validatorWorkerCount)
	for i := range validatorWorkerCount {
		ch := make(chan ProcessedProduct, productCount)
		workerChannels[i] = ch
		result[i] = ch
	}

	for workerId := range validatorWorkerCount {
		go func(ctx context.Context, workerId int) {
			defer close(workerChannels[workerId])

			for product := range products {

				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{
						workerChannels[workerId] <- validate(workerId, product)
					}

				}
			}
		}(ctx, workerId)
	}

	return result
}

func MergeChannels(ctx context.Context, channels []<-chan ProcessedProduct) chan ProcessedProduct {
	result := make(chan ProcessedProduct)
	var wg sync.WaitGroup
	for _, c := range channels {
		wg.Add(1)
		go func(ctx context.Context, c <-chan ProcessedProduct) {
			defer wg.Done()

			for v := range c {

				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{

						result <- v
					}
				}
			}

		}(ctx, c)
	}
	go func() {
		wg.Wait()
		close(result)
	}()
	return result

}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	products := EnrichedProducts(ctx)
	processedProducts := TrustValidate(ctx, products)
	merged := MergeChannels(ctx, processedProducts)
	fmt.Println("Processes products:")

	for item := range merged {
		fmt.Println(item)
	}

}
