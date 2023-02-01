package controller

import (
	"evm-blockchain-parser/lib/http"
	"evm-blockchain-parser/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetCurrentBlock(ctx *gin.Context) {
	http.ResponseJson(ctx, model.LatestScannedBlock.Number, nil)
}

func Subcribe(ctx *gin.Context) {
	address := ctx.Query("address")
	if len(address) != 42 && len(address) != 40 {
		http.ResponseJson(ctx, nil, fmt.Errorf("address [%s] length invalid", address))
		return
	}
	err := model.Subcribe(address)
	if err != nil {
		http.ResponseJson(ctx, nil, fmt.Errorf("subcribtion failed with error: %s", err.Error()))
		return
	}
	http.ResponseJson(ctx, nil, nil)
}

func GetTransactions(ctx *gin.Context) {
	address := ctx.Query("address")
	if len(address) != 42 && len(address) != 40 {
		http.ResponseJson(ctx, nil, fmt.Errorf("address [%s] length invalid", address))
		return
	}
	addressInfo, ok := model.AddressMap[address]
	if !ok {
		http.ResponseJson(ctx, nil, fmt.Errorf("address [%s] is not subcribed", address))
		return
	}
	result := make([]model.Transaction, 0, 1000)
	for _, value := range addressInfo.TxHashList {
		transaction, ok := model.TransactionHashMap[value]
		if !ok {
			http.ResponseJson(ctx, nil, fmt.Errorf("get transaction by hash [%s] failed", value))
			return
		}
		result = append(result, transaction)
	}
	http.ResponseJson(ctx, result, nil)
}
