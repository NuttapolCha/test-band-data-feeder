package cache

import (
	"fmt"
	"sync"
)

type Pricing struct {
	Price     float64
	Timestamp int64
}

type symbolMapPricing map[string]*Pricing

type LatestPricing struct {
	mu sync.Mutex
	m  symbolMapPricing
}

var ltsp LatestPricing

func init() {
	ltsp = LatestPricing{
		mu: sync.Mutex{},
		m:  make(symbolMapPricing),
	}
}

func GetLatestPrice(symbol string) (*Pricing, error) {
	ltsp.mu.Lock()
	defer ltsp.mu.Unlock()

	pricing, ok := ltsp.m[symbol]
	if !ok {
		return nil, fmt.Errorf("not found pricing information for symbol %s", symbol)
	}

	return pricing, nil
}

func UpdatePriceInfo(symbol string, price float64, timestamp int64) error {
	ltsp.mu.Lock()
	defer ltsp.mu.Unlock()

	ltsp.m[symbol] = &Pricing{
		Price:     price,
		Timestamp: timestamp,
	}

	return nil
}
