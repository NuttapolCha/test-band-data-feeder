package cache

import (
	"fmt"
	"time"
)

func GetLatestPriceFromDataSource(symbol string) (*Pricing, error) {
	latestDataSource.mu.Lock()
	defer latestDataSource.mu.Unlock()

	pricing, ok := latestDataSource.m[symbol]
	if !ok {
		return nil, fmt.Errorf("data source cache has not been inited")
	}

	return pricing, nil
}

func GetLatestPriceFromDestination(symbol string) (*Pricing, error) {
	if time.Now().After(destinationCacheUpdatedAt.Add(liveTime)) {
		return nil, fmt.Errorf("destination cache is old")
	}

	latestDestination.mu.Lock()
	defer latestDestination.mu.Unlock()

	pricing, ok := latestDestination.m[symbol]
	if !ok {
		return nil, fmt.Errorf("destination cache has not been inited")
	}

	return pricing, nil
}

func UpdatePriceToDataSource(symbol string, price float64, timestamp int64) error {
	latestDataSource.mu.Lock()
	defer latestDataSource.mu.Unlock()

	latestDataSource.m[symbol] = &Pricing{
		Price:     price,
		Timestamp: timestamp,
	}

	return nil
}

func UpdatePriceToDestination(symbol string, price float64, timestamp int64) error {
	latestDestination.mu.Lock()
	defer latestDestination.mu.Unlock()

	latestDestination.m[symbol] = &Pricing{
		Price:     price,
		Timestamp: timestamp,
	}

	return nil
}

func UpdateDestinationCacheTime() {
	destinationCacheUpdatedAt = time.Now()
}
