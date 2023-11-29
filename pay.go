package main

import (
	"context"
	"fmt"
	"golang.org/x/exp/constraints"
)

func test() {
	fmt.Println(pay(10, []denomination{twenty}))
}

func pay(price float32, inserted []denomination) (map[denomination]int, error) {
	m := map[denomination]int{
		cent:        10,
		twentyCents: 10,
		fiftyCents:  10,
		one:         10,
		two:         10,
		five:        10,
		ten:         0,
		twenty:      10,
		fifty:       10,
	}

	reg := NewRegister(m, orderedDenominations)

	// processPaymentInCents is concurrent function. Can be called safely from concurrent goroutines.
	ds, err := processPaymentInCents(context.Background(), reg, int(price*100), inserted)
	if err != nil {
		return nil, fmt.Errorf("error processing payment: %w", err)
	}
	returnDenominations := make(map[denomination]int)
	for _, d := range ds {
		returnDenominations[d]++
	}
	return returnDenominations, nil
}

func processPaymentInCents(ctx context.Context, reg registerer[denomination], price int, inserted []denomination) ([]denomination, error) {
	insertedAmount := 0
	for _, d := range inserted {
		insertedAmount += int(d)
	}
	if insertedAmount < price {
		return inserted, ErrInsufficientFunds
	}
	err := reg.Put(ctx, inserted)
	if err != nil {
		return inserted, err
	}

	returnAmount := insertedAmount - price
	ds, err := reg.Withdraw(ctx, returnAmount)
	if err != nil {
		return nil, err
	}
	if len(ds) == 0 {
		for _, i := range inserted {
			d, _ := reg.Withdraw(ctx, int(i))
			ds = append(ds, d...)
		}
		return ds, ErrRegisterEmpty
	}
	return ds, nil
}
