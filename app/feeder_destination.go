package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NuttapolCha/test-band-data-feeder/app/pricing"
	"github.com/NuttapolCha/test-band-data-feeder/cache"
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

func (app *App) getPricingFromDst(symbol string, config *FeederConfig) (*DestinationPricingResp, error) {
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

// need update pricing criterias are
// 1. no new update than 1 hour (configurable)
// 2. price difference is more than threshold 0.1 (configurable)
func (app *App) isNeedUpdatePricingToDestination(
	prevPricing,
	currPricing pricing.Information,
	config *FeederConfig,
) (is, immediatly bool) {
	logger := app.logger
	symbol := prevPricing.GetSymbol()

	// optional: checking if prevPricing and currPricing are the same symbol
	if symbol != currPricing.GetSymbol() {
		errMsg := fmt.Sprintf("cannot compare pricing information of difference symbols (%s and %s)", symbol, currPricing.GetSymbol())
		panic(errMsg)
	}

	prevTime := prevPricing.GetTimestamp()
	currTime := time.Now().Unix()
	timeDiff := currTime - prevTime
	logger.Debugf("previous cache time of %s = %v", symbol, prevTime)
	logger.Debugf("current time = %v", currTime)
	logger.Debugf("time diff of %s = %v", symbol, time.Duration(timeDiff)*time.Second)

	if timeDiff >= config.maximumDelay {
		logger.Infof("RELAY: symbol %s has not send update longer than %vs and need to be updated", symbol, config.maximumDelay)
		return true, false
	}

	prevPrice := prevPricing.GetPrice()
	currPrice := currPricing.GetPrice()
	logger.Debugf("preious price of %s = %f", symbol, prevPrice)
	logger.Debugf("current price of %s = %f", symbol, currPrice)

	// use absolute value
	priceDiffRatio := (currPrice - prevPrice) / currPrice
	if priceDiffRatio < 0 {
		priceDiffRatio *= -1
	}
	logger.Debugf("price diff ratio of %s = %.4f", symbol, priceDiffRatio)

	if priceDiffRatio > config.diffThreshold {
		logger.Infof("URGENT: symbol %s has difference grater than threshold %v and need to be updated immediatly", symbol, config.diffThreshold)
		return true, true
	}

	logger.Infof("symbol %s no need to send update at destination because delay = %v and diff ratio = %f", symbol, time.Duration(timeDiff)*time.Second, priceDiffRatio)
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
		logger.Infof("successfully updated pricing information of %+v prices %+v at timestamp %v", params.Symbols, params.Prices, params.Timestamp)
	}

	dstTime := time.Now().Unix()
	for _, symbol := range updatedSymbols {
		if err := cache.UpdateDstTime(symbol, dstTime); err != nil {
			logger.Errorf("could not update destination timestamp to cache of %s because: %v", symbol, err)
		}
	}

	return updatedSymbols, err
}
