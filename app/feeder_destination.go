package app

type symbolPrice struct {
	symbol []string
	price  []float64
}

type TimeStampPricing map[int64]*symbolPrice

// toDestParams constructs payload model in which we use for updating to destination service.
// It also filters out non neccessary for update symbols; i.e. DOEST NOT meet the following 2 conditions
//
// 1. latest price is delay more than configured (1 hour)
//
// 2. price difference ratio is grater than threshold (0.1)
func (m TimeStampPricing) toDstParams() []*UpdatePricingParams {
	ret := make([]*UpdatePricingParams, 0, len(m))

	for timeStamp, pair := range m {
		ret = append(ret, &UpdatePricingParams{
			Symbols:   pair.symbol,
			Prices:    pair.price,
			Timestamp: timeStamp,
		})
	}

	return ret
}
