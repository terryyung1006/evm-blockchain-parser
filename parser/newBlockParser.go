package parser

import (
	"evm-blockchain-parser/cronjob"
	"evm-blockchain-parser/lib/http"
	"evm-blockchain-parser/model"
	"log"
	"strconv"
	"time"
)

type NewBlockParser struct {
	BlockParser
	cronjob.CronJob
}

func (nbp NewBlockParser) Run(interval int) {
	defer time.Sleep(time.Duration(nbp.Interval) * time.Millisecond)
	if len(model.AddressMap) == 0 {
		log.Printf("[newBlockParser] no address added, skip worker")
		return
	}
	latestScannedBlockNum, err := strconv.Atoi(model.LatestScannedBlock.Number)
	if err != nil {

	}
	latestBlockNum, _ := http.GetLatestBlockNumber()
	if int(latestBlockNum) == latestScannedBlockNum {
		return
	}
	addressList := make([]string, 0, 1000)
	for address, addressInfo := range model.AddressMap {
		if addressInfo.LastBlock <= latestScannedBlockNum {
			continue
		}
		addressList = append(addressList, address)
	}
	latestBlock, err := http.GetBlockByNumber(int(latestBlockNum))
	if err != nil {
		log.Printf("[GetBlockByNumber] failed in new block parser with err: [%s] for block: %v, skip this cronjob", err.Error(), latestBlockNum)
		return
	}
	nbp.ScanTransactionInBlock(latestBlock, addressList)
	model.UpdateLastBlock(int(latestBlockNum), addressList)
	model.LatestScannedBlock = latestBlock
}
