package main

import (
	"data_collector/client"
	"data_collector/etherscan_data"
	"data_collector/models"
	"data_collector/models/etherscan"
	"data_collector/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	utils.InitLogsCores()
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		logs.Error(err)
	}
	etherscanApiKey := os.Getenv("ETHERSCAN_API_KEY")

	swapsFile, err := os.OpenFile("swaps.csv", os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		logs.Error(err)
		os.Exit(-1)
	}

	var swaps []models.Swap
	err = gocsv.UnmarshalFile(swapsFile, &swaps)
	if err != nil {
		logs.Error(err)
		os.Exit(-1)
	}
	err = swapsFile.Close()
	if err != nil {
		logs.Error(err)
	}

	etherscanClient := client.NewHttpClient(15*time.Second, 200*time.Millisecond)
	err = etherscanClient.StartDelayer()
	if err != nil {
		logs.Error(err)
		os.Exit(-1)
	}
	logs.Info("Data collection started; Total data size: %d.", len(swaps))
	queueCh := make(chan struct{}, 10)
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	for i := 0; i < len(swaps); i++ {
		wg.Add(1)
		go func(i int) {
			queueCh <- struct{}{}
			defer func() { <-queueCh }()
			defer wg.Done()
			mu.RLock()
			swap := swaps[i]
			mu.RUnlock()
			url := fmt.Sprintf(etherscan_data.GET_TRANSACTION_BY_HASH, swap.Transaction, etherscanApiKey)
			var transaction etherscan.TransactionResponse
			err := etherscanClient.DoRequest(http.MethodGet, url, nil, nil, nil, &transaction)
			for k := 0; err != nil && k < 10; k++ {
				err = etherscanClient.DoRequest(http.MethodGet, url, nil, nil, nil, &transaction)
			}
			if err != nil {
				logs.Error(err)
				return
			}
			var block etherscan.BlockResponse
			url = fmt.Sprintf(etherscan_data.GET_BLOCK_BY_NUMBER, transaction.Result.BlockNumber, etherscanApiKey)
			err = etherscanClient.DoRequest(http.MethodGet, url, nil, nil, nil, &block)
			for k := 0; err != nil && k < 10; k++ {
				err = etherscanClient.DoRequest(http.MethodGet, url, nil, nil, nil, &block)
			}
			if err != nil {
				logs.Error(err)
				return
			}
			if len(block.Result.Timestamp) > 3 {
				timestamp, err := strconv.ParseInt(block.Result.Timestamp[2:], 16, 64)
				if err != nil {
					logs.Error(err)
					return
				} else {
					mu.Lock()
					swaps[i].Timestamp = timestamp
					mu.Unlock()
				}
			} else {
				logs.Error("Wrong timestamp %v. Index: %d, transaction: %+v, block: %+v.", block.Result.Timestamp, i, transaction, block)
			}
		}(i)

		if i%1000 == 0 {
			wg.Wait()
			mu.RLock()
			swapsFile2, err := os.OpenFile("swaps2.csv", os.O_WRONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				logs.Error(err)
				os.Exit(-1)
			}
			err = gocsv.Marshal(swaps, swapsFile2)
			if err != nil {
				logs.Error(err)
			}
			err = swapsFile2.Close()
			if err != nil {
				logs.Error(err)
			}
			logs.Info(i)
			mu.RUnlock()
		}
	}
	wg.Wait()
	mu.RLock()
	swapsFile2, err := os.OpenFile("swaps2.csv", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		logs.Error(err)
		os.Exit(-1)
	}

	err = gocsv.Marshal(swaps, swapsFile2)
	if err != nil {
		logs.Error(err)
	}
	err = swapsFile2.Close()
	if err != nil {
		logs.Error(err)
	}
	mu.RUnlock()
	logs.Info("Data collection ended")
}
