package app

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type RequestPricingDataSourceParams struct {
	Symbols []string `json:"symbols"`
}

type RequestPricingDataSourceResp struct {
	ID int `json:"id"`
}

type PricingResult struct {
	Multiplier  string `json:"multiplier"`
	Px          string `json:"px"`
	RequestID   string `json:"request_id"`
	ResolveTime string `json:"resolve_time"`
	Symbol      string `json:"symbol"`
}

func (p *PricingResult) GetSymbol() string {
	return p.Symbol
}

func (p *PricingResult) GetPrice() float64 {
	multipliedPrice, err := strconv.ParseFloat(p.Px, 64)
	if err != nil {
		panic(err)
	}
	multiplier, err := strconv.ParseFloat(p.Multiplier, 64)
	if err != nil {
		panic(err)
	}
	return multipliedPrice / multiplier
}

func (p *PricingResult) GetTimestamp() int64 {
	t, err := strconv.ParseInt(p.ResolveTime, 10, 64)
	if err != nil {
		panic(err)
	}
	return t
}

type PricingResultResp struct {
	PricingResults []*PricingResult `json:"price_results"`
}

func (app *App) requestPricingFromSource(config *FeederConfig) (int, error) {
	logger := app.logger

	bs, err := json.Marshal(&RequestPricingDataSourceParams{
		Symbols: config.symbols,
	})
	if err != nil {
		logger.Errorf("could not marshal RequestPricingDataSourceParams to JSON payload because: %v", err)
		return -1, err
	}

	respBody, err := app.httpClient.PostJSON(config.requestPricingDataEndpoint, bs, config.dataSourceRetryCount)
	if err != nil {
		logger.Errorf("could not PostJSON because: %v", err)
		return -1, err
	}
	ref := &RequestPricingDataSourceResp{}
	if err := json.Unmarshal(respBody, ref); err != nil {
		logger.Errorf("could not unmarshal request pricing data source response into GO struct because: %v", err)
		return -1, err
	}

	return ref.ID, nil
}

func (app *App) getRequestedPricingFromSource(reqId int, config *FeederConfig) ([]*PricingResult, error) {
	logger := app.logger

	pricingEndpoint := fmt.Sprintf("%s/%d", config.getPricingDataEndpoint, reqId)
	respBody, err := app.httpClient.Get(pricingEndpoint, nil, config.dataSourceRetryCount)
	if err != nil {
		logger.Errorf("could not get the requested pricing data from source because: %v", err)
		return nil, err
	}

	pricingResp := &PricingResultResp{}
	if err := json.Unmarshal(respBody, pricingResp); err != nil {
		logger.Errorf("could not unmarshal pricing result into GO struct because: %v", err)
		return nil, err
	}

	return pricingResp.PricingResults, nil
}
