package model

import (
	"evm-blockchain-parser/lib/slice"
	"sync"
)

// memory storage for transaction
type Transaction struct {
	BlockHash        string `json:"blockHAsh"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

var TransactionHashMap = map[string]Transaction{}

func AddTransaction(transaction Transaction, mu *sync.Mutex) {
	existingHashList := slice.GetSliceOfKeys(TransactionHashMap)
	if slice.Contains(existingHashList, transaction.Hash) {
		// log.Fatalf("[AddTransaction]transaction: %s already exists", transaction.Hash)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	TransactionHashMap[transaction.Hash] = transaction
}
