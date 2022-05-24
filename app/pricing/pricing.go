package pricing

type Information interface {
	GetSymbol() string
	GetPrice() float64
	GetTimestamp() int64
}
