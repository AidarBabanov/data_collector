package etherscan_data

const (
	GET_TRANSACTION_BY_HASH = `https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=%s`
)
