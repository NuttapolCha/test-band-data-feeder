package app

import (
	"encoding/json"
)

type UpdatePricingParams struct {
	Symbols   []string  `json:"symbols"`
	Prices    []float64 `json:"prices"`
	Timestamp int64     `json:"timestamp"`
}

type DestinationPricingResp struct {
	Price      float64 `json:"price"`
	LastUpdate int64   `json:"last_update"`
}

func (p *DestinationPricingResp) GetPrice() float64 {
	return p.Price
}

func (p *DestinationPricingResp) GetTimestamp() int64 {
	return p.LastUpdate
}

func (app *App) getLatestPriceFromDestination(symbol string, config *FeederConfig) (*DestinationPricingResp, error) {
	logger := app.logger

	body, err := app.httpClient.Get(config.getUpdatedPricingData, map[string]string{
		"symbol": symbol,
	}, config.destinationRetryCount)
	if err != nil {
		logger.Errorf("could not http GET because: %v", err)
		return nil, err
	}

	ret := &DestinationPricingResp{}
	if err := json.Unmarshal(body, ret); err != nil {
		logger.Errorf("could not unmarshal destination pricing body into GO struct because: %v", err)
		return nil, err
	}

	return ret, nil
}

type pricingWithTimestamp interface {
	GetPrice() float64
	GetTimestamp() int64
}

// need update pricing criterias are
// 1. no new update than 1 hour (configurable)
// 2. pricing difference is more than threshold 0.1 (configurable)
func (app *App) isNeedUpdatePricingToDestination(
	symbol string,
	latestDataFromCache,
	latestDataFromDst pricingWithTimestamp,
	config *FeederConfig,
) (is, immediatly bool) {
	logger := app.logger

	latestCacheTime := latestDataFromCache.GetTimestamp()
	currentTime := latestDataFromDst.GetTimestamp()
	timeDiff := latestCacheTime - currentTime
	logger.Debugf("latest cache time of %s = %d", symbol, latestCacheTime)
	logger.Debugf("latest destination time of %s = %d", symbol, currentTime)
	logger.Debugf("time diff of %s = %d", symbol, timeDiff)

	if timeDiff >= config.maximumDelay {
		logger.Infof("symbol %s has not send update longer than %vs and need to be updated", symbol, config.maximumDelay)
		return true, false
	}

	latestCachePrice := latestDataFromCache.GetPrice()
	latestDstPrice := latestDataFromDst.GetPrice()
	logger.Debugf("latest cache price of %s = %.4f", symbol, latestCachePrice)
	logger.Debugf("latest destination price of %s = %.4f", symbol, latestDstPrice)

	// use absolute value
	priceDiffRatio := (latestCachePrice - latestDstPrice) / latestDstPrice
	if priceDiffRatio < 0 {
		priceDiffRatio *= -1
	}
	logger.Debugf("price diff ratio of %s = %.4f", symbol, priceDiffRatio)

	if priceDiffRatio > config.diffThreshold {
		logger.Infof("symbol %s has difference grater than threshold %v and need to be updated imediatly !!!", symbol, config.diffThreshold)
		return true, true
	}

	logger.Infof("symbol %s no need to send update at destination because delay = %ds and diff ratio = %.4f", symbol, timeDiff, priceDiffRatio)
	return false, false
}

func (app *App) updatePricingToDestination(updatePricingParamsList []*UpdatePricingParams, config *FeederConfig) ([]string, error) {
	logger := app.logger

	updatedSymbols := make([]string, 0, len(updatePricingParamsList))

	var err error
	for _, params := range updatePricingParamsList {
		var reqBody []byte

		reqBody, err = json.Marshal(params)
		if err != nil {
			logger.Errorf("could not marshal UpdatePricingParams to GO struct because: %v", err)
			continue // current params error, try next
		}

		_, err = app.httpClient.PostJSON(config.updatePricingDataEndpoint, reqBody, config.destinationRetryCount)
		if err != nil {
			logger.Errorf("could not PostJSON because: %v", err)
			continue // current params error, try next
		}
		updatedSymbols = append(updatedSymbols, params.Symbols...)
		logger.Infof("successfully update pricing information of %+v prices %+v at timestamp %v", params.Symbols, params.Prices, params.Timestamp)
	}

	return updatedSymbols, err
}
