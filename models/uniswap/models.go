package uniswap

import jsoniter "github.com/json-iterator/go"

type Token struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Pair struct {
	Token0 Token `json:"token0"`
	Token1 Token `json:"token1"`
}

type Swap struct {
	Amount0In  jsoniter.Number `json:"amount0In"`
	Amount1In  jsoniter.Number `json:"amount1In"`
	Amount0Out jsoniter.Number `json:"amount0Out"`
	Amount1Out jsoniter.Number `json:"amount1Out"`
	Pair       Pair            `json:"pair"`
}

type Transaction struct {
	Id    string `json:"id"`
	Swaps []Swap `json:"swaps"`
}

type TransactionResponse struct {
	Data struct {
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}

type TransactionsResponse struct {
	Data struct {
		Transactions []Transaction `json:"transactions"`
	} `json:"data"`
}
