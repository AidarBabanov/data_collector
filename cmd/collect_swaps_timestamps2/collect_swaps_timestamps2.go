package main

import (
	"data_collector/client"
	"data_collector/etherplorer_data"
	"data_collector/models"
	"data_collector/models/etherplorer"
	"data_collector/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"net/http"
	"os"
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
	manager := client.NewManager(15*time.Second, 127*time.Millisecond,
		"EK-wxHtH-TsiS9mC-L7L9q",
		"EK-k3NuB-5nkpo73-fAWWh",
		"EK-gvBaW-kRGqqqd-wUqsS",
		"EK-4BgiR-SqrBSwq-WywCm",
		"EK-fJjCg-zpPvbNm-hdjYL",
		"EK-daRTQ-RZ9hqEq-YUqwQ",
		"EK-uwxUP-KNjKU1U-EfGds",
		"EK-rrZwL-esiyuuq-1bjSm",
	)
	err = manager.Serve()
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
			var transaction etherplorer.TransactionInfo
			var err error
			f := func(cl *client.Client) {
				url := fmt.Sprintf(etherplorer_data.GET_TRANSACTION_INFO, swap.Transaction, cl.Key)
				err = cl.DoRequest(http.MethodGet, url, nil, nil, nil, &transaction)
			}
			manager.UseClient(f)
			for k := 0; err != nil && k < 10; k++ {
				manager.UseClient(f)
			}
			if err != nil {
				logs.Error(err)
				return
			}
			mu.Lock()
			swaps[i].BlockNumber = transaction.BlockNumber
			swaps[i].Timestamp = transaction.Timestamp
			mu.Unlock()
		}(i)

		if i%1000 == 0 {
			wg.Wait()
			mu.RLock()
			swapsFile2, err := os.OpenFile("swaps2.csv", os.O_WRONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				logs.Error(err)
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
