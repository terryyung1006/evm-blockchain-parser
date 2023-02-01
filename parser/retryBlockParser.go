package parser

import (
	"evm-blockchain-parser/cronjob"
	"evm-blockchain-parser/lib/http"
	"log"
	"sync"
	"time"
)

// memory storage for failed job
var retryJob = map[int][]string{}

func AddRetryJob(blockNum int, addressList []string, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	retryJob[blockNum] = addressList
}

func DeleteRetryJob(blockNum int, mu *sync.Mutex) {
	if _, ok := retryJob[blockNum]; !ok {
		// log.Printf("[DeleteRetryJob]retry job for block: %v already exists", blockNum)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	delete(retryJob, blockNum)
}

type RetryBlockParser struct {
	BlockParser BlockParser
	Cronjob     cronjob.CronJob
}

func (rbp RetryBlockParser) Run() {
	defer time.Sleep(time.Duration(rbp.Cronjob.Interval) * time.Millisecond)

	if len(retryJob) == 0 {
		log.Printf("[RetryBlockParser] no retry job, skip worker")
		return
	}

	var wg sync.WaitGroup
	count := 0
	for blockNumber, addressList := range retryJob {
		if count > rbp.BlockParser.JobCapacity {
			break
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, blockNumber int, addressList []string) {
			defer wg.Done()
			log.Printf("[RetryBlockParser]Rescanning job %v", blockNumber)
			block, err := http.GetBlockByNumber(int(blockNumber))
			if err != nil {
				log.Printf("[RetryBlockParser]Retry job failed with error [%s], block num: %v", err.Error(), blockNumber)
				return
			}
			rbp.BlockParser.ScanTransactionInBlock(block, addressList)
			DeleteRetryJob(blockNumber, rbp.BlockParser.Mu)
		}(&wg, blockNumber, addressList)
		count++
	}
	wg.Wait()
}
