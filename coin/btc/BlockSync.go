// Process new block and unconfirmed transactions (via RPC).
package main

import (
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
	//return true
}
