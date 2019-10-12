package main

import (
	"log"
	"sync"
	"time"

	ltc "github.com/GGBTC/explorer/ltc"
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
		res := ltc.CatchUpBlocks()
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
		res := ltc.CatchUpTx()
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
		res := ltc.GetLastBitCoinPrice()
		if res == "Success" {
			time.Sleep(6 * time.Minute)
			log.Println("==========LastPrice==========")
			log.Println("============Finish=============")
		}
	}
}
