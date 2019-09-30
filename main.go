package main

import (
	"fmt"
	"log"
	"time"

	"github.com/GGBTC/explorer/service"

	"github.com/GGBTC/explorer/btc"
)

func main() {
	//txid := "hahah"
	res_tx, err := btc.CallBTCRPC(service.BTCURL, "gettxoutsetinfo", 1, []interface{}{})
	//(URL, "getblockcount", 1, []interface{}{})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	tar := res_tx["result"].(map[string]interface{})
	fmt.Println(tar)
	// go newTask() //新建一个goroutine
	// //	AnoutherTask()
	// AnoutherTask1()

}

func newTask() {
	for {
		fmt.Println("this is a new Task.")
		time.Sleep(time.Duration(5) * time.Second) //延时1s
	}
}

func AnoutherTask() {
	for {
		fmt.Println("Fuck")
		time.Sleep(time.Duration(1) * time.Second) //延时1s
	}
}
func AnoutherTask1() {
	for {
		fmt.Println("HAHAHAHAHAHAH")
		time.Sleep(time.Duration(1) * time.Second) //延时1s
	}
}
