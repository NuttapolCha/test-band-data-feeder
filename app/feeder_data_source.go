package app

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/NuttapolCha/test-band-data-feeder/cache"
)

func (app *App) requestPricingFromSource(config *FeederConfig) (int, error) {
	bs, err := json.Marshal(&RequestPricingDataSourceParams{
		Symbols: config.symbols,
	})
	if err != nil {
		app.logger.Errorf("could not marshal RequestPricingDataSourceParams to JSON payload because: %v", err)
		return -1, err
	}

	respBody, err := app.httpClient.PostJSON(config.requestPricingDataEndpoint, bs, config.dataSourceRetryCount)
	if err != nil {
		app.logger.Errorf("could not PostJSON because: %v", err)
		return -1, err
	}
	ref := &RequestPricingDataSourceResp{}
	if err := json.Unmarshal(respBody, ref); err != nil {
		app.logger.Errorf("could not unmarshal request pricing data source response into GO struct because: %v", err)
		return -1, err
	}

	return ref.ID, nil
}

func (app *App) getRequestedPricingFromSource(reqId int, config *FeederConfig) ([]PricingResult, error) {
	pricingEndpoint := fmt.Sprintf("%s/%d", config.getPricingDataEndpoint, reqId)
	respBody, err := app.httpClient.Get(pricingEndpoint, nil, config.dataSourceRetryCount)
	if err != nil {
		app.logger.Errorf("could not get the requested pricing data from source because: %v", err)
		return nil, err
	}

	pricingResp := &PricingResultResp{}
	if err := json.Unmarshal(respBody, pricingResp); err != nil {
		app.logger.Errorf("could not unmarshal pricing result into GO struct because: %v", err)
		return nil, err
	}

	return pricingResp.PricingResults, nil
}

func (app *App) cachePricingResults(results []PricingResult) error {
	for _, pricing := range results {
		multipliedPrice, err := strconv.ParseFloat(pricing.Px, 64)
		if err != nil {
			app.logger.Errorf("could not parse Px: %s to float64 because: %v", pricing.Px, err)
			return err
		}
		multiplier, err := strconv.ParseFloat(pricing.Multiplier, 64)
		if err != nil {
			app.logger.Errorf("could not parse Multiplier: %s to float64 because: %v", pricing.Px, err)
			return err
		}
		updatedTime, err := strconv.ParseInt(pricing.ResolveTime, 10, 64)
		if err != nil {
			app.logger.Errorf("could not parse ResolveTime: %s to int64 because: %v", pricing.ResolveTime, err)
			return err
		}

		price := multipliedPrice / multiplier
		app.logger.Infof("price of %s = %.4f USD", pricing.Symbol, price)

		err = cache.UpdatePriceInfo(pricing.Symbol, price, updatedTime)
		if err != nil {
			app.logger.Errorf("could not update price info of %s at %d because: %v", pricing.Symbol, updatedTime, err)
			return err
		}
	}

	return nil
}