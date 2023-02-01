package parser

import (
	"evm-blockchain-parser/lib/slice"
	"evm-blockchain-parser/model"
)

type BlockParser struct {
	jobCapacity int
}

func (bp BlockParser) ScanTransactionInBlock(block model.Block, addressList []string) {
	for _, transaction := range block.Transactions {
		var matchedAddress string
		if slice.Contains(addressList, transaction.From) {
			matchedAddress = transaction.From
		}
		if slice.Contains(addressList, transaction.To) {
			matchedAddress = transaction.To
		}
		if matchedAddress != "" {
			model.AddTransaction(transaction)
			model.AddTransactionHash(matchedAddress, transaction.Hash)
		}
	}
}
