package main

import (
	"log"
	"time"

	zec "github.com/GGBTC/explorer/zec"
)

func main() {
	//var wg sync.WaitGroup
	Tx()
	//Block()
	//wg.Add(1)

	//wg.Wait()
	//UTXO()

}

func Block() {
	for {
		res := zec.CatchUpBlocks()
		if res == "Success" {
			time.Sleep(5 * time.Second)
			log.Println("============Block==============")
			log.Println("============Finish=============")
		}
	}
}
func Tx() {
	for {
		//	defer wg.Done()
		res := zec.CatchUpTx()
		if res == "Success" {
			time.Sleep(2 * time.Second)
			log.Println("==========Transaction==========")
			log.Println("============Finish=============")
		}
	}
}
