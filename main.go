package main

import (
	"log"
	"sync"
	"time"

	zec "github.com/GGBTC/explorer/zec"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go Tx()
	wg.Add(1)
	go Price()
	Block()
	wg.Wait()
	//Cop()

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

func Price() {
	for {
		//	defer wg.Done()
		res := zec.GetLastBitCoinPrice()
		if res == "Success" {
			time.Sleep(6 * time.Minute)
			log.Println("==========LastPrice==========")
			log.Println("============Finish=============")
		}
	}
}
