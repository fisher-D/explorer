package main

import (
	"fmt"
	"time"

	"github.com/GGBTC/explorer/pkgs"
	"github.com/GGBTC/explorer/service"
)

const (
	//URL       = "http://admin:admin@47.244.98.227:13143"
	//GenesisTx = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"
	mongourl = "localhost:27017"
)

func main() {
	//pkgs.CatchUpBlocks()
	service.GetMongo(mongourl)
	running := true
	for running {
		KeepRunning()
	}
}
func KeepRunning() {
	res := pkgs.AutoDrawTx()
	if res == "Success" {
		fmt.Println("Waiting 1 minute to start next Process")
		time.Sleep(time.Duration(60) * time.Second)
	}
	//return true
}
