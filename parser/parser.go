package parser

import (
	"evm-blockchain-parser/lib/slice"
	"evm-blockchain-parser/model"
	"sync"
)

type BlockParser struct {
	JobCapacity int
	Mu          *sync.Mutex
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
			model.AddTransaction(transaction, bp.Mu)
			model.AddTransactionHash(matchedAddress, transaction.Hash, bp.Mu)
		}
	}
}
