/*

	This code searches for unique addresses which made swaps on uniswap.org.

	1. Loads transactions from uniswap sorted by time.
	2. Selects transactions which have swaps.
	3. Finds initiator of the transaction from etherscan.
	4. Saves address of initiator into csv file.

*/
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
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	addresses = set.New()
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		logs.Error(err)
	}
	etherscan_api_key := os.Getenv("ETHERSCAN_API_KEY")

	// create output file with addresses
	// file is just one column of addresses
	file, err := os.Create("addresses.csv")
	if err != nil {
		logs.Error(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// create clients for etherscan and uniswap
	uniswapClient := client.NewGraphQLClient(uniswap_data.BASE_URL, 15*time.Second, 275*time.Millisecond)
	etherscanClient := client.NewHttpClient(15*time.Second, 275*time.Millisecond)
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
		// 'skip' parameter skips certain amount of transactions (in our we skip 1000 transactions each time)
		// 'first' parameter takes certain amount of transactions for response (in our case we take 1000 transactions)
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
					// transaction.Id is hash of transaction
					// each request in etherscan needs a key, etherscan_api_key is key parameter
					url := fmt.Sprintf(etherscan_data.GET_TRANSACTION_BY_HASH, transaction.Id, etherscan_api_key)
					var transaction etherscan.Transaction
					// use for loop to be sure that data is loaded
					for i := 0; i < 10; i++ {
						err = etherscanClient.DoRequest(http.MethodGet, url, nil, nil, nil, &transaction)
						if err == nil {
							break
						}
						if err != nil {
							logs.Error(err)
						}
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
