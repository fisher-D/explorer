package main

import (
	"fmt"
	"log"
	"time"

	"github.com/GGBTC/explorer/service"

	zec "github.com/GGBTC/explorer/zcash"
)

func main() {
	//3
	txid := "9a4adaf3953818eb1634407032db0e00ef2441c49c1364161411d0743ec1a939"
	//2
	//txid := "8974d08d1c5f9c860d8b629d582a56659a4a1dcb2b5f98a25a5afcc2a784b0f4"
	//txid := "ae61ee40c37fdc05468786cf6e83a98741a214c41fc5d62057f768d9b0af769e"
	res_tx, err := zec.CallZECRPC(service.ZECURL, "getrawtransaction", 1, []interface{}{txid, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	fmt.Println(res_tx)
}

// func main() {
// 	//pkgs.CatchUpBlocks()
// 	//	service.GetMongo(mongourl)
// 	running := true
// 	for running {
// 		KeepRunning()
// 	}
// }
func KeepRunning() {
	res := zec.CatchUpBlockss()
	if res == "Success" {
		time.Sleep(time.Duration(60) * time.Second)
	}
	//return true
}
