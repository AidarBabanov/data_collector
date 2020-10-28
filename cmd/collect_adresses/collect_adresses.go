package main

import (
	"data_collector/client"
	"data_collector/etherscan_data"
	"data_collector/models/etherscan"
	"data_collector/models/uniswap"
	"data_collector/uniswap_data"
	"encoding/csv"
	"fmt"
	"github.com/Workiva/go-datastructures/set"
	"github.com/astaxie/beego/logs"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	addresses = set.New()
)

func main() {
	file, err := os.Create("addresses.csv")
	if err != nil {
		logs.Error(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	uniswapClient := client.NewGraphQLClient(uniswap_data.BASE_URL, 15*time.Second, 200*time.Millisecond)
	etherscanClient := client.NewHttpClient(15*time.Second, 200*time.Millisecond)
	err = uniswapClient.StartDelayer()
	if err != nil {
		logs.Error(err)
	}
	err = etherscanClient.StartDelayer()
	if err != nil {
		logs.Error(err)
	}
	defer uniswapClient.Close()
	defer etherscanClient.Close()
	skip := 0
	first := 1000
	logs.Info("Data collection started")
	for addresses.Len() < 1e6 {
		//get 1k transactions
		getTransactionsQuery := fmt.Sprintf(uniswap_data.GET_TRANSACTIONS, skip, first)
		var transactions uniswap.TransactionsData
		err = uniswapClient.DoGraphqlRequest(getTransactionsQuery, &transactions)
		if err != nil {
			logs.Error(err)
		}
		addressesTemp := set.New()
		wg := new(sync.WaitGroup)
		for _, transaction := range transactions.Data.Transactions {
			// check if it has swaps
			if len(transaction.Swaps) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					// get transaction info from etherscan
					url := fmt.Sprintf(etherscan_data.GET_TRANSACTION_BY_HASH, transaction.Id, etherscan_data.API_KEY)
					var transaction etherscan.Transaction
					for i := 0; i < 10; i++ {
						err = etherscanClient.DoRequest(http.MethodGet, url, nil, nil, nil, &transaction)
						if err == nil {
							break
						}
					}
					if err != nil {
						logs.Error(err)
						return
					}

					// add address
					addr := transaction.Result.From
					if !addresses.Exists(addr) {
						addresses.Add(addr)
						addressesTemp.Add(addr)
					}
				}()
			}
		}
		wg.Wait()
		skip += 1000
		for _, addr := range addressesTemp.Flatten() {
			err = writer.Write([]string{addr.(string)})
			if err != nil {
				logs.Error(err)
			}
		}
		logs.Info("Passed %d transactions, saved: %d addresses", skip, addresses.Len())
	}
	logs.Info("Data collection ended")
	logs.Info(addresses.Len())
}
