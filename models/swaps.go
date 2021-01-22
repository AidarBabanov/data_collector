package models

type Swap struct {
	Transaction       string  `csv:"Transaction"`
	Address           string  `json:"Address"`
	FromTokenId       string  `csv:"FromTokenId"`
	FromTokenSymbol   string  `json:"FromTokenSymbol"`
	FromTokenName     string  `json:"FromTokenName"`
	ToTokenId         string  `json:"ToTokenId"`
	ToTokenSymbol     string  `json:"ToTokenSymbol"`
	ToTokenName       string  `json:"ToTokenName"`
	FromAmount        float64 `json:"FromAmount"`
	ToAmount          float64 `json:"ToAmount"`
	Gas               int64   `json:"Gas"`
	CumulativeGasUsed int64   `json:"CumulativeGasUsed"`
	GasUsed           int64   `json:"GasUsed"`
	BlockNumber       string  `json:"BlockNumber"`
	Timestamp         int64   `json:"Timestamp"`
}
