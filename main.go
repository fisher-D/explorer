package main

import (
	"fmt"
	"log"

	"github.com/GGBTC/explorer/btc"
	s "github.com/GGBTC/explorer/service"
)

const (
	starttime = 0
	endtime   = 32
)
const (
	URL       = s.BTCURL
	GenesisTx = s.BTCGenesisTx
	mongourl  = s.Mongourl
	apikey    = s.APIKEY
)

func main() {
	jobCh := genJob(starttime, endtime)
	//wg.Add(1)
	workerPool(4, jobCh)
	//time.Sleep(time.Duration(1) * time.Second)
	//Use this Select to hold the main process
	fmt.Println("good")

}

func genJob(star int, n int) <-chan int {
	jobCh := make(chan int, 200)
	go func() {
		for i := star; i < n; i++ {
			jobCh <- i
		}
		close(jobCh)
	}()
	//	fmt.Println("len--------------", len(jobCh))
	return jobCh
}

func workerPool(n int, jobCh <-chan int) {
	for i := 0; i < n; i++ {
		go worker(i, jobCh)
	}
	//return true
}

func worker(id int, jobCh <-chan int) {
	for job := range jobCh {
		CatchUpBlock(job)
		//	fmt.Println("Processer:", id, "  Result", job)
	}
}
func CatchUpBlock(start int) string {
	s.GetMongo(mongourl)
	Database := s.GlobalS.DB("BTC")
	blockCollection := Database.C("blocks")
	hash := btc.GetbtcHashRPC(int64(start))
	blocks := btc.GetBlocks(hash)
	log.Println("Process Block With Height: ", start, "; And Blocks Hash :", hash)
	blockCollection.Insert(blocks)
	txArray := blocks.Tx
	result := btc.CatchUpTx(txArray, Database)
	if result != true {
		return "Failed"
	}
	return "Success"
}
