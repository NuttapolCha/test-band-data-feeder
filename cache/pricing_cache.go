package cache

import (
	"fmt"
	"sync"
)

// we might call pricingWithTimestamp a virtual database of this service
type pricingWithTimestamp struct {
	symbol string

	// price is the latest price we known at destination
	price float64

	// updateDstTime is a time at which we called for update destination service
	updateDstTime int64

	// dstTime is a timestamp appearred in update destination request body
	// (i.e. time of symbol price)
	dstTime int64
}

func (p *pricingWithTimestamp) GetSymbol() string {
	return p.symbol
}

func (p *pricingWithTimestamp) GetPrice() float64 {
	return p.price
}

func (p *pricingWithTimestamp) GetTimestamp() int64 {
	return p.dstTime
}

type symbolMapPricing map[string]*pricingWithTimestamp

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

// GetPrevPrice return pricing
func GetPricing(symbol string) (*pricingWithTimestamp, error) {
	ltsp.mu.Lock()
	defer ltsp.mu.Unlock()

	pricing, ok := ltsp.m[symbol]
	if !ok {
		return nil, fmt.Errorf("cache has not been inited")
	}

	return pricing, nil
}

func GetPrevUpdatedDstTime(symbol string) (int64, error) {
	ltsp.mu.Lock()
	defer ltsp.mu.Unlock()

	pricing, ok := ltsp.m[symbol]
	if !ok {
		return 0, fmt.Errorf("cache has not been inited")
	}

	return pricing.updateDstTime, nil
}

func UpdatePricing(
	symbol string,
	price float64,
	updateDstTime,
	dstTime int64,
) {
	if symbol == "" {
		panic("cannot update with empty string key")
	}
	ltsp.mu.Lock()
	defer ltsp.mu.Unlock()

	ltsp.m[symbol] = &pricingWithTimestamp{
		symbol:        symbol,
		price:         price,
		updateDstTime: updateDstTime,
		dstTime:       dstTime,
	}
}
