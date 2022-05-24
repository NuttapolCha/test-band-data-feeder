package cache

import (
	"fmt"
	"sync"
)

// we might call pricingWithTimestamp a virtual database of this service
type pricingWithTimestamp struct {
	symbol string
	price  float64

	updateDstTime int64
	dstTime       int64
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

func UpdatePricing(symbol string, price float64) error {
	ltsp.mu.Lock()
	defer ltsp.mu.Unlock()

	curr, ok := ltsp.m[symbol]
	if ok {
		ltsp.m[symbol] = &pricingWithTimestamp{
			symbol:        symbol,
			price:         price,
			updateDstTime: curr.updateDstTime,
			dstTime:       curr.dstTime,
		}
	} else {
		ltsp.m[symbol] = &pricingWithTimestamp{
			symbol: symbol,
			price:  price,
		}
	}

	return nil
}

// UpodateDstTime updates 2 timestamps.
// 1. updateDstTime is a time at which we called for update destination service
// 2. dstTime is a timestamp in 1.'s request body (i.e. time of symbol price)
func UpdateDstTime(symbol string, updateDstTime, dstTime int64) error {
	ltsp.mu.Lock()
	defer ltsp.mu.Unlock()

	curr, ok := ltsp.m[symbol]
	if !ok {
		return fmt.Errorf("not found %s in cache", symbol)
	}

	if symbol != curr.symbol {
		panic("invalid symbol caching and this should not be occurred!")
	}

	ltsp.m[symbol] = &pricingWithTimestamp{
		symbol:        curr.symbol,
		price:         curr.price,
		updateDstTime: updateDstTime,
		dstTime:       dstTime,
	}

	return nil
}
