package model

import (
	"evm-blockchain-parser/lib/slice"
	"log"
)

type AddressInfo struct {
	LastBlock  int
	TxHashList []string
}

var AddressMap map[string]*AddressInfo

func Subcribe(address string) error {
	_, ok := AddressMap[address]
	if ok {

	}
	addressInfo := AddressInfo{
		LastBlock:  0,
		TxHashList: []string{},
	}
	AddressMap[address] = &addressInfo
	return nil
}

func AddTransactionHash(address string, hash string) {
	if slice.Contains(AddressMap[address].TxHashList, hash) {
		log.Fatalf("[AddTransactionHash] hash %s already exist in address info map (address: %s) hash list", hash, address)
		return
	}
	(*AddressMap[address]).TxHashList = append((*AddressMap[address]).TxHashList, hash)
}

func UpdateLastBlock(UpdatedBlock int, addressList []string) {
	for _, address := range addressList {
		addressInfo, ok := AddressMap[address]
		if !ok {
			log.Fatalf("[UpdateLastBlock] address %s not found in address info map", address)
			continue
		}
		addressInfo.LastBlock = UpdatedBlock
	}
}
