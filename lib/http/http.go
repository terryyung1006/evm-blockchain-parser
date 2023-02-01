package http

import (
	"bytes"
	"encoding/json"
	"evm-blockchain-parser/model"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpPost(url string, body interface{}, result interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("[HttpPost] req body marshal with error: %s", err.Error())
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("[HttpPost] post failed with error: %s", err.Error())
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[HttpPost] ioutil ReadAll resp body failed with error: %s", err.Error())
	}
	err = json.Unmarshal([]byte(respBody), result)
	if err != nil {
		if strings.Contains(string(respBody), "Rate limiting threshold exceeded") {
			defer time.Sleep(time.Duration(3) * time.Second)
			return fmt.Errorf("[HttpPost] Rate limiting threshold exceeded")

		}
		return fmt.Errorf("[HttpPost] response unmarshal failed with error: %s", err.Error())
	}
	return nil
}

type BlockchainRPCRequest struct {
	Id      int           `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func BlockchainHttpPost(method string, params []interface{}, result interface{}) error {
	body := BlockchainRPCRequest{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
	err := HttpPost("https://cloudflare-eth.com", body, result)
	if err != nil {
		return fmt.Errorf("[BlockchainHttpPost] method %s post request failed with err: %s", method, err.Error())
	}
	return nil
}

type GetLatestBlockNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func GetLatestBlockNumber() (int64, error) {
	params := []interface{}{}
	resp := GetLatestBlockNumberResponse{}
	err := BlockchainHttpPost("eth_blockNumber", params, &resp)
	if err != nil {
		return -1, fmt.Errorf("[GetLatestBlockNumber] failed with err: %s", err.Error())
	}
	blockNumber, err := strconv.ParseInt(resp.Result[2:], 16, 64)
	if err != nil {
		return -1, fmt.Errorf("[GetLatestBlockNumber] result block number parse from hex to decimal failed: %s", err.Error())
	}
	return blockNumber, nil
}

type GetBlockByNumberResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  model.Block `json:"result"`
}

func GetBlockByNumber(blockNumber int) (model.Block, error) {
	blockNumberInHex := fmt.Sprintf("0x%x", blockNumber)
	params := []interface{}{blockNumberInHex, true}
	resp := GetBlockByNumberResponse{}
	err := BlockchainHttpPost("eth_getBlockByNumber", params, &resp)
	if err != nil {
		return model.Block{}, fmt.Errorf("[GetBlockByNumber] failed with err: %s", err.Error())
	}
	return resp.Result, nil
}

func ResponseJson(ctx *gin.Context, data interface{}, err error) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": data,
		})
	}
}
