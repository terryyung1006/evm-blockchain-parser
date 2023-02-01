package main

import (
	"evm-blockchain-parser/controller"
	"evm-blockchain-parser/cronjob"
	"evm-blockchain-parser/model"
	"evm-blockchain-parser/parser"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	model.StartFromBlockNum = 600000 - 1

	var mu sync.Mutex
	newBlockParser := parser.NewBlockParser{
		BlockParser: parser.BlockParser{
			JobCapacity: 1,
			Mu:          &mu,
		},
		Cronjob: cronjob.CronJob{
			Interval: 10000,
		},
	}
	oldBlockParser := parser.OldBlockParser{
		BlockParser: parser.BlockParser{
			JobCapacity: 10,
			Mu:          &mu,
		},
		Cronjob: cronjob.CronJob{
			Interval: 100,
		},
	}
	retryJobParser := parser.RetryBlockParser{
		BlockParser: parser.BlockParser{
			JobCapacity: 10,
			Mu:          &mu,
		},
		Cronjob: cronjob.CronJob{
			Interval: 5000,
		},
	}
	cronjob.RunCronJobs([]cronjob.Icronjob{newBlockParser, oldBlockParser, retryJobParser})

	router := gin.Default()

	router.GET("/api/get_current_block", controller.GetCurrentBlock)
	router.POST("/api/subscribe", controller.Subcribe)
	router.POST("/api/get_transactions", controller.GetTransactions)

	router.Run(":8080")

	fmt.Println("Exit")
}
