package uniswap_data

const (
	BASE_URL         = `https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2`
	GET_TRANSACTIONS = `
{
	"query":
	"{transactions (orderBy: timestamp, orderDirection: desc, skip: %d, first: %d){id swaps{id}}}",
	"variables":{}
}`
)
