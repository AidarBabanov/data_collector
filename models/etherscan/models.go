package etherscan

import jsoniter "github.com/json-iterator/go"

type Transaction struct {
	Hash              string          `json:"hash"`
	From              string          `json:"from"`
	To                string          `json:"to"`
	Timestamp         jsoniter.Number `json:"timeStamp"`
	Gas               jsoniter.Number `json:"gas"`
	CumulativeGasUsed jsoniter.Number `json:"cumulativeGasUsed"`
	GasUsed           jsoniter.Number `json:"GasUsed"`
	IsError           jsoniter.Number `json:"isError"` // should be "0"
}

type TransactionResponse struct {
	Result Transaction `json:"result"`
}

type AddressTransactionsResponse struct {
	Status  jsoniter.Number `json:"status"`
	Message string          `json:"message"`
	Result  []Transaction   `json:"result"`
}
