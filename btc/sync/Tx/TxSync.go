// Process new block and unconfirmed transactions (via RPC).
package main

import (
	"time"

	"github.com/GGBTC/explorer/pkgs"
)

func main() {
	//pkgs.CatchUpBlocks()
	//service.GetMongo("loaclhost:27017")
	running := true
	for running {
		KeepRunning()
	}
}
func KeepRunning() {
	res := pkgs.CatchUpTx()
	if res == "Success" {
		time.Sleep(time.Duration(60) * time.Second)
	}
	//return true
}
