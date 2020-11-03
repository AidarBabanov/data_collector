package models

type Swap struct {
	Transaction       string
	Address           string
	FromTokenId       string
	FromTokenSymbol   string
	FromTokenName     string
	ToTokenId         string
	ToTokenSymbol     string
	ToTokenName       string
	FromAmount        float64
	ToAmount          float64
	Gas               int64
	CumulativeGasUsed int64
	GasUsed           int64
}
