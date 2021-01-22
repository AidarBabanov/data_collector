package etherscan_data

const (
	GET_TRANSACTION_BY_HASH     = `https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=%s`
	GET_TRANSACTIONS_BY_ADDRESS = `https://api.etherscan.io/api?module=account&action=txlist&address=%s&page=%d&offset=%dsort=asc&apikey=%s` // address, page, offset, api key
	UNISWAP_CONTRACT_ADDRESS    = `0x7a250d5630b4cf539739df2c5dacb4c659f2488d`
	GET_BLOCK_BY_NUMBER         = `https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s` // block number, api key
)
