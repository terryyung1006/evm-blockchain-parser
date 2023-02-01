package main

import (
	"evm-blockchain-parser/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/get_current_block", controller.GetCurrentBlock)
	router.POST("/api/subscribe", controller.Subcribe)
	router.POST("/api/get_transactions", controller.GetTransactions)

	router.Run(":8080")

	fmt.Println("Exit")
}
