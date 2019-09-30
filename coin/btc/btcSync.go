// Process new block and unconfirmed transactions (via RPC).
package main

import (
	"time"

	"github.com/GGBTC/explorer/btc"
)

func main() {
	//pkgs.CatchUpBlocks()
	//	service.GetMongo(mongourl)
	go UpdateInformation()
	KeepRunning()
}
func KeepRunning() {
	for {
		res := btc.CatchUpBlockss()
		if res == "Success" {
			time.Sleep(time.Duration(60) * time.Second)
		}
	}

	//return true
}

func UpdateInformation() {
	for {
		resu := btc.GetLastBitCoinPrice()
		if resu == "Success" {
			time.Sleep(time.Duration(360) * time.Second)
		}
	}

}
