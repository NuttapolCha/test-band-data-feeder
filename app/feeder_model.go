package app

// Data Source

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

type PricingResultResp struct {
	PricingResults []PricingResult `json:"price_results"`
}

// Destination

type UpdatePricingParams struct {
	Symbols   []string  `json:"symbols"`
	Prices    []float64 `json:"prices"`
	Timestamp int64     `json:"timestamp"`
}

type UpdatedPricingResp struct {
	LastUpdate int64   `json:"last_update"`
	Price      float64 `json:"price"`
}
