package etherplorer

type TransactionInfo struct {
	Hash        string `json:"hash"`
	Timestamp   int64  `json:"timestamp"`
	BlockNumber int64  `json:"blockNumber"`
	//Confirmations int64  `json:"confirmations"`
	//Success       bool   `json:"success"`
	//From          string `json:"from"`
	//To            string `json:"to"`
	//Value         int64  `json:"value"`
	//Input         string `json:"input"`
	//GasLimit      int64  `json:"gasLimit"`
	//GasUsed       int64  `json:"gasUsed"`
}
