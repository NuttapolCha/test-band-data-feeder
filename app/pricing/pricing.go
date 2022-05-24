package pricing

type Information interface {
	GetSymbol() string
	GetPrice() float64
	GetTimestamp() int64
}

func Equal(x, y Information) bool {
	return x.GetSymbol() == y.GetSymbol() &&
		x.GetPrice() == y.GetPrice() &&
		x.GetTimestamp() == y.GetTimestamp()
}
