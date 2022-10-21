package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var api Worker = NewAPIRateLimited()

	for i := 0; i < 10; i++ {
		err := api.Work(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type Worker interface {
	Work(ctx context.Context) error
}

type API struct{}

func (a API) Work(ctx context.Context) error {
	fmt.Println("not supposed to call often")

	return nil
}

type APIRateLimited struct {
	rateLimiter *rate.Limiter
}

func NewAPIRateLimited() *APIRateLimited {
	return &APIRateLimited{rateLimiter: rate.NewLimiter(rate.Limit(1), 3)}
}

func (a *APIRateLimited) Work(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}

	fmt.Println("not supposed to call often")
	return nil
}
