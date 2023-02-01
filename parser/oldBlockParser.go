package parser

import (
	"evm-blockchain-parser/cronjob"
	"evm-blockchain-parser/lib/http"
	"evm-blockchain-parser/lib/slice"
	"evm-blockchain-parser/model"
	"log"
	"sync"
	"time"
)

type OldBlockParser struct {
	BlockParser BlockParser
	Cronjob     cronjob.CronJob
}

func (obp OldBlockParser) Run() {
	defer time.Sleep(time.Duration(obp.Cronjob.Interval) * time.Millisecond)

	if len(model.AddressMap) == 0 {
		// log.Printf("[OldBlockParser] no address added, skip worker")
		return
	}

	latestBlockNum, _ := http.GetLatestBlockNumber()

	addressLastBlockList := make([]int, 0, 1000)
	addressLastBlockMap := map[string]int{}
	for address, addressInfo := range model.AddressMap {
		if addressInfo.LastBlock == int(latestBlockNum)-1 {
			continue
		}
		addressLastBlockList = append(addressLastBlockList, addressInfo.LastBlock)
		addressLastBlockMap[address] = addressInfo.LastBlock
	}
	var earliestBlockNum int
	if len(addressLastBlockList) == 0 {
		earliestBlockNum = 0
	} else {
		earliestBlockNum = slice.MinInt(addressLastBlockList)
	}

	jobMap := map[int]*[]string{}
	blockRange := int(latestBlockNum) - earliestBlockNum
	if blockRange > obp.BlockParser.JobCapacity {
		blockRange = obp.BlockParser.JobCapacity
	}
	for address, lastBlockNum := range addressLastBlockMap {
		for i := 1; i <= blockRange; i++ {
			if _, ok := jobMap[earliestBlockNum+i]; !ok {
				var addressList []string
				jobMap[earliestBlockNum+i] = &addressList
			}
			if lastBlockNum < (earliestBlockNum + i) {
				*jobMap[earliestBlockNum+i] = append(*jobMap[earliestBlockNum+i], address)
			}
		}
	}

	var wg sync.WaitGroup
	for blockNumber, addressList := range jobMap {
		wg.Add(1)
		go func(wg *sync.WaitGroup, blockNumber int, addressList []string) {
			defer wg.Done()
			block, err := http.GetBlockByNumber(int(blockNumber))
			if err != nil {
				//add retry job
				log.Printf("[GetBlockByNumber] failed with err: [%s], for block: %v, skip this cronjob", err.Error(), blockNumber)
				AddRetryJob(blockNumber, addressList, obp.BlockParser.Mu)
				return
			}
			obp.BlockParser.ScanTransactionInBlock(block, addressList)
			model.UpdateLastBlock(int64(blockNumber), addressList, obp.BlockParser.Mu)
		}(&wg, blockNumber, *addressList)
	}
	wg.Wait()
}
