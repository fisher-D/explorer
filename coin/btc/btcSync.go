// Process new block and unconfirmed transactions (via RPC).
package main

import (
	"log"
	"time"

	"github.com/GGBTC/explorer/btc"
)

func main() {
	//pkgs.CatchUpBlocks()
	//	service.GetMongo(mongourl)
	running := true
	for running {
		KeepRunning()
	}
}
func KeepRunning() {
	res := btc.CatchUpBlockss()
	if res == "Success" {
		time.Sleep(time.Duration(60) * time.Second)
	}
	resu := btc.GetLastBitCoinPrice()
	if resu == "Success" {
		log.Println("Succecss And Waiting for next rand")
		time.Sleep(time.Duration(60) * time.Second)
	}
	//return true
}
