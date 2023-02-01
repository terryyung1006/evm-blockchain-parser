package model

import (
	"evm-blockchain-parser/lib/slice"
	"fmt"
	"log"
	"sync"
)

// memory storage for subcribed addresses and transaction hash
type AddressInfo struct {
	LastBlock  int
	TxHashList []string
}

var AddressMap = map[string]*AddressInfo{}

func Subcribe(address string) error {
	_, ok := AddressMap[address]
	if ok {
		return fmt.Errorf("[Subcribe] address %s already subceibed", address)
	}
	addressInfo := AddressInfo{
		LastBlock:  StartFromBlockNum,
		TxHashList: []string{},
	}
	AddressMap[address] = &addressInfo
	return nil
}

func AddTransactionHash(address string, hash string, mu *sync.Mutex) {
	if _, ok := AddressMap[address]; !ok {
		log.Fatalf("[AddTransactionHash] AddressMap doesnt contain address: %s", address)
		return
	}

	if slice.Contains(AddressMap[address].TxHashList, hash) {
		log.Fatalf("[AddTransactionHash] hash %s already exist in address info map (address: %s) hash list", hash, address)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	log.Printf("[AddTransactionHash] transaction found! hash: %s, address: %s", hash, address)
	(*AddressMap[address]).TxHashList = append((*AddressMap[address]).TxHashList, hash)
}

func UpdateLastBlock(UpdatedBlock int64, addressList []string, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	for _, address := range addressList {
		addressInfo, ok := AddressMap[address]
		if !ok {
			log.Fatalf("[UpdateLastBlock] address %s not found in address info map", address)
			continue
		}
		if int(UpdatedBlock) > addressInfo.LastBlock {
			addressInfo.LastBlock = int(UpdatedBlock)
			if UpdatedBlock > LatestScannedBlock {
				LatestScannedBlock = UpdatedBlock
			}
		}
	}
}
