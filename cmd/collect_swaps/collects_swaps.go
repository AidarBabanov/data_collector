/*
   	read all addresses
      for all addresses:
      	get 1 address
      	until the end of transactions:
      		get 1k addresses
      		for all addresses:
      			get 1 transaction
      			check if it has swaps
      			get input info
      			get output info
      			get transaction info from etherscan
      			get gas info
      			save swap info
      		save addresses

*/

package main

import (
	"data_collector/client"
	"data_collector/etherscan_data"
	"data_collector/models"
	"data_collector/models/etherscan"
	"data_collector/models/uniswap"
	"data_collector/uniswap_data"
	"data_collector/utils"
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/joho/godotenv"
	"github.com/mohae/struct2csv"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	addresses []string
)

func main() {
	utils.InitLogsCores()
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		logs.Error(err)
	}
	etherscanApiKey := os.Getenv("ETHERSCAN_API_KEY")

	if err != nil {
		logs.Error(err)
	}

	// create output outputFile with addresses
	outputFile, err := os.Create("swaps.csv")
	if err != nil {
		logs.Error(err)
	}
	defer outputFile.Close()
	writer := struct2csv.NewWriter(outputFile)
	defer writer.Flush()

	err = writer.WriteColNames(models.Swap{})
	if err != nil {
		logs.Error(err)
	}

	// create clients for etherscan and uniswap
	uniswapClient := client.NewGraphQLClient(uniswap_data.BASE_URL, 15*time.Second, 75*time.Millisecond)
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

	inputFile, err := os.Open("addresses.csv")
	if err != nil {
		logs.Error(err)
	}
	reader := csv.NewReader(inputFile)
	defer inputFile.Close()

	addressesMatrix, err := reader.ReadAll()
	if err != nil {
		logs.Error(err)
	}

	for _, address := range addressesMatrix {
		addresses = append(addresses, address[0])
	}

	logs.Info("Data collection started")
	page := 0
	offset := 10000
	totalSwaps := 0

	for addrIndex, address := range addresses {
		url := fmt.Sprintf(etherscan_data.GET_TRANSACTIONS_BY_ADDRESS, address, page, offset, etherscanApiKey)
		addressTxs := etherscan.AddressTransactionsResponse{}
		for i := 0; i < 10; i++ {
			err = etherscanClient.DoRequest(http.MethodGet, url, nil, nil, nil, &addressTxs)
			if err == nil {
				break
			}
		}
		if err != nil {
			logs.Error(err)
		}

		if addressTxs.Status != "1" || addressTxs.Message != "OK" {
			continue
		}
		buffer := make(chan struct{}, runtime.NumCPU())
		mu := new(sync.Mutex)
		var swaps []models.Swap
		wg := new(sync.WaitGroup)
		for _, tx := range addressTxs.Result {
			buffer <- struct{}{}
			wg.Add(1)
			go func(tx etherscan.Transaction) {
				defer func() {
					<-buffer
					wg.Done()
				}()
				if tx.IsError != "0" || tx.To != etherscan_data.UNISWAP_CONTRACT_ADDRESS {
					return
				}
				uniswapTx := uniswap.TransactionResponse{}
				for i := 0; i < 10; i++ {
					url := fmt.Sprintf(uniswap_data.GET_TRANSACTION, tx.Hash)
					err = uniswapClient.DoGraphqlRequest(url, &uniswapTx)
					if err == nil {
						break
					}
				}
				if err != nil {
					logs.Error(err)
					return
				}
				if len(uniswapTx.Data.Transaction.Swaps) == 0 {
					return
				}

				from := uniswapTx.Data.Transaction.Swaps[0]
				fromAmount0, _ := from.Amount0In.Float64()
				fromAmount1, _ := from.Amount1In.Float64()
				fromAmount := 0.0
				fromToken := uniswap.Token{}
				if fromAmount0 > 0 && fromAmount1 == 0 {
					fromAmount = fromAmount0
					fromToken = from.Pair.Token0
				} else if fromAmount0 == 0 && fromAmount1 > 0 {
					fromAmount = fromAmount1
					fromToken = from.Pair.Token1
				} else {
					logs.Error("Skipped; %+v %+v", tx, uniswapTx)
					return
				}

				to := uniswapTx.Data.Transaction.Swaps[len(uniswapTx.Data.Transaction.Swaps)-1]
				toAmount0, _ := to.Amount0Out.Float64()
				toAmount1, _ := to.Amount1Out.Float64()
				toAmount := 0.0
				toToken := uniswap.Token{}
				if toAmount0 > 0 && toAmount1 == 0 {
					toAmount = toAmount0
					toToken = to.Pair.Token0
				} else if toAmount0 == 0 && toAmount1 > 0 {
					toAmount = toAmount1
					toToken = to.Pair.Token1
				} else {
					logs.Error("Skipped; %+v %+v", tx, uniswapTx)
					return
				}

				gas, _ := tx.Gas.Int64()
				cumulativeGasUsed, _ := tx.Gas.Int64()
				gasUsed, _ := tx.Gas.Int64()
				swap := models.Swap{
					Transaction:       tx.Hash,
					Address:           address,
					FromTokenId:       fromToken.Id,
					FromTokenSymbol:   fromToken.Symbol,
					FromTokenName:     fromToken.Name,
					FromAmount:        fromAmount,
					ToTokenId:         toToken.Id,
					ToTokenSymbol:     toToken.Symbol,
					ToTokenName:       toToken.Name,
					ToAmount:          toAmount,
					Gas:               gas,
					CumulativeGasUsed: cumulativeGasUsed,
					GasUsed:           gasUsed,
				}
				mu.Lock()
				swaps = append(swaps, swap)
				mu.Unlock()

				totalSwaps++
			}(tx)
		}
		wg.Wait()
		if len(swaps) > 0 {
			err = writer.WriteStructs(swaps)
			if err != nil {
				logs.Error(err)
				return
			}
		}
		logs.Info("Addresses managed: %d; swaps saved: %d", addrIndex+1, totalSwaps)
	}

	logs.Info("Data collection ended")
}
