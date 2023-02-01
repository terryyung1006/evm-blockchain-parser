package parser

import (
	"evm-blockchain-parser/cronjob"
	"evm-blockchain-parser/lib/http"
	"evm-blockchain-parser/lib/slice"
	"evm-blockchain-parser/model"
	"log"
	"strconv"
	"sync"
	"time"
)

type OldBlockParser struct {
	BlockParser
	cronjob.CronJob
}

func (obp OldBlockParser) Run(interval int) {
	defer time.Sleep(time.Duration(interval) * time.Millisecond)

	if len(model.AddressMap) == 0 {
		log.Printf("[OldBlockParser] no address added, skip worker")
		return
	}

	latestScannedBlockNum, err := strconv.Atoi(model.LatestScannedBlock.Number)
	if err != nil {

	}
	latestBlockNum, _ := http.GetLatestBlockNumber()

	addressLastBlockList := make([]int, 0, 1000)
	addressLastBlockMap := map[string]int{}
	for address, addressInfo := range model.AddressMap {
		if addressInfo.LastBlock == latestScannedBlockNum {
			continue
		}
		addressLastBlockList = append(addressLastBlockList, addressInfo.LastBlock)
		addressLastBlockMap[address] = addressInfo.LastBlock
	}
	earliestBlockNum := slice.MinInt(addressLastBlockList)

	jobMap := map[int]*[]string{}
	blockRange := int(latestBlockNum) - earliestBlockNum
	if blockRange > obp.jobCapacity {
		blockRange = obp.jobCapacity
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
				AddRetryJob(blockNumber, addressList)
				return
			}
			obp.ScanTransactionInBlock(block, addressList)
			model.UpdateLastBlock(int(blockNumber), addressList)
		}(&wg, blockNumber, *addressList)
	}
	wg.Wait()
}
