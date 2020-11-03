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

var (
	GET_TRANSACTION = "{\"query\":\"{\\n  transaction(id:\\\"%s\\\"){\\n    id\\n    swaps{\\n      pair{\\n        token0{\\n          id\\n          symbol\\n          name\\n        }\\n        token1{\\n          id\\n          symbol\\n          name\\n        }\\n      }\\n      \\n      amount0In\\n      amount1In\\n      \\n      amount0Out\\n      amount1Out\\n    }\\n  }\\n}\\n\",\"variables\":{}}"
)
