package parser

import (
	"evm-blockchain-parser/cronjob"
	"evm-blockchain-parser/lib/http"
	"evm-blockchain-parser/model"
	"log"
	"time"
)

type NewBlockParser struct {
	BlockParser BlockParser
	Cronjob     cronjob.CronJob
}

func (nbp NewBlockParser) Run() {
	defer time.Sleep(time.Duration(nbp.Cronjob.Interval) * time.Millisecond)
	if len(model.AddressMap) == 0 {
		// log.Printf("[newBlockParser] no address added, skip worker")
		return
	}
	latestScannedBlockNum := model.LatestScannedBlock
	latestBlockNum, err := http.GetLatestBlockNumber()
	if err != nil {
		log.Printf("[NewBlockParser] failed in GetLatestBlockNumber with err: [%s], skip this cronjob", err.Error())
		return
	}
	if latestScannedBlockNum != latestBlockNum-1 {
		return
	}
	addressList := make([]string, 0, 1000)
	for address, addressInfo := range model.AddressMap {
		if addressInfo.LastBlock != int(latestBlockNum)-1 {
			continue
		}
		addressList = append(addressList, address)
	}
	latestBlock, err := http.GetBlockByNumber(int(latestBlockNum))
	if err != nil {
		log.Printf("[NewBlockParser] failed in new block parser with err: [%s] for block: %v, skip this cronjob", err.Error(), latestBlockNum)
		return
	}
	nbp.BlockParser.ScanTransactionInBlock(latestBlock, addressList)
	model.UpdateLastBlock(latestBlockNum, addressList, nbp.BlockParser.Mu)
	model.LatestScannedBlock = latestBlockNum
}
