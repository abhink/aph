package main

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds inserted")
	ErrRegisterEmpty     = errors.New("register does not have adequate denominations")
)

func NewRegister(contents map[denomination]int, order []denomination) *Register {
	return &Register{
		contents:             contents,
		orderedDenominations: order,
		m:                    sync.Mutex{},
	}
}

// Register is concurent access safe type that implements the `registerer` interface for `denomination` type.
type Register struct {
	contents             map[denomination]int
	orderedDenominations []denomination
	m                    sync.Mutex
}

func (r *Register) Put(_ context.Context, ds []denomination) error {
	r.m.Lock()
	defer r.m.Unlock()

	for _, d := range ds {
		r.contents[d]++
	}
	return nil
}

func (r *Register) Withdraw(_ context.Context, amount int) ([]denomination, error) {
	r.m.Lock()
	defer r.m.Unlock()
	var j = 0
	var ds []denomination
	for amount > 0 {
		currDenom := r.orderedDenominations[j]
		if j == len(r.orderedDenominations) {
			return nil, nil
		} else if amount-int(currDenom) >= 0 && r.contents[currDenom] > 0 {
			amount -= int(currDenom)
			ds = append(ds, currDenom)
			r.contents[currDenom]--
		} else {
			j++
		}
	}
	return ds, nil
}

func (r *Register) OrderedDenominations(_ context.Context) ([]denomination, error) {
	return r.orderedDenominations, nil
}
