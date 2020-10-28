package uniswap

type Swap struct {
	Id string `json:"id"`
}

type Transaction struct {
	Id    string `json:"id"`
	Swaps []Swap `json:"swaps"`
}

type TransactionsData struct {
	Data struct {
		Transactions []Transaction `json:"transactions"`
	} `json:"data"`
}
