package etherscan

type Transaction struct {
	Result struct {
		BlockHash string `json:"blockHash"`
		From      string `json:"from"`
	} `json:"result"`
}
