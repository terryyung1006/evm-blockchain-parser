EVM Blockchain Parser

steps:
1, go mod tidy
2, go run main.go

API interface:
    a. GET localhost:8080/api/get_current_block
        -get latest parsed block
    b. POST localhost:8080/api/subscribe?address="0x45060b5cee190661fa27d1e189f431f7b2b52275"
        -subscribe address for transaction scanning
    c. POST localhost:8080/api/get_transactions?address="0x45060b5cee190661fa27d1e189f431f7b2b52275"
        -get scanned transaction for address