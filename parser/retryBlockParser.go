package parser

import (
	"evm-blockchain-parser/cronjob"
	"evm-blockchain-parser/lib/http"
	"log"
	"sync"
	"time"
)

var retryJob map[int][]string

func AddRetryJob(blockNum int, addressList []string) {
	retryJob[blockNum] = addressList
}

func DeleteRetryJob(blockNum int) {
	if _, ok := retryJob[blockNum]; !ok {
		log.Printf("[DeleteRetryJob]retry job for block: %v already exists", blockNum)
		return
	}
	delete(retryJob, blockNum)
}

type RetryBlockParser struct {
	BlockParser
	cronjob.CronJob
}

func (rbp *RetryBlockParser) Run() {
	defer time.Sleep(time.Duration(rbp.Interval) * time.Millisecond)

	if len(retryJob) == 0 {
		log.Printf("[RetryBlockParser] no retry job, skip worker")
		return
	}

	var wg sync.WaitGroup
	count := 0
	for blockNumber, addressList := range retryJob {
		if count > rbp.jobCapacity {
			break
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, blockNumber int, addressList []string) {
			defer wg.Done()
			block, err := http.GetBlockByNumber(int(blockNumber))
			if err != nil {
				log.Printf("Retry job failed with error [%s], block num: %v", err.Error(), blockNumber)
				return
			}
			rbp.ScanTransactionInBlock(block, addressList)
			DeleteRetryJob(blockNumber)
		}(&wg, blockNumber, addressList)
		count++
	}
	wg.Wait()
}
