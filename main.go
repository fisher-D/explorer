package main

import (
	"encoding/json"
	"fmt"

	zec "github.com/GGBTC/explorer/zcash"
)

func main() {
	//ltctxid := "52c7a210b9b93fde2a4ca44317e8f1ad83a32cc2bb4078d36dcd649f44e95619"
	//	btctxid := "a91103494a6bea71bea569706638297b5efe2ac33581725240acdb32ba828e2f"
	//omnitxid := "ae79d9120a6f5cea1e6d2beaf497a1388611c18fdbd142f804dcf0fbe48c7dc8"
	height := zec.GetzecCountRPC()
	hash := zec.GetzecHashRPC(height)
	block := zec.GetBlocks(hash)
	data, _ := json.Marshal(block)

	fmt.Println(string(data))
	fmt.Println("===================================")
	//fmt.Println(string(data1))

}
