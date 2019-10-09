package main

import (
	"fmt"
	"log"
	"time"

	"github.com/GGBTC/explorer/btc"
	s "github.com/GGBTC/explorer/service"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

const (
	URL       = s.BTCURL
	GenesisTx = s.BTCGenesisTx
	mongourl  = s.Mongourl
)

func run(task_id, sleeptime int, ch chan string) {

	time.Sleep(time.Duration(sleeptime) * time.Second)
	ch <- fmt.Sprintf("task id %d , sleep %d second", task_id, sleeptime)
	return
}
func main() {
	//	Number()
	GetTxArray()
}
func Number() {
	input := []int{3, 2, 1}
	//ch := make(chan string)
	chs := make([]chan string, len(input))
	startTime := time.Now()
	fmt.Println("Multirun start")
	for i, sleeptime := range input {
		//go run(i, sleeptime, ch)
		chs[i] = make(chan string)
		go run(i, sleeptime, chs[i])
	}
	//for range input {
	//	fmt.Println(<-ch)
	//}
	for _, ch := range chs {
		fmt.Println(<-ch)
	}

	endTime := time.Now()
	fmt.Printf("Multissh finished. Process time %s. Number of tasks is %d", endTime.Sub(startTime), len(input))
}
func runstr(txid string, ch chan *s.Tx, Database *mgo.Database) {
	//fmt.Println("fUCK")
	//time.Sleep(time.Duration(100) * time.Millisecond)
	res := CatchTx(txid, Database)

	//ch <- fmt.Sprintf("task id %d , sleep %d second", task_id, sleeptime)
	ch <- res
	return
}
func Character(txArray []string, Database *mgo.Database) {
	//input := []string{"a", "2", "1"}
	//ch := make(chan string)
	// var num int
	// if len(txArray) >= 5 {
	// 	num = 5
	// } else {
	// 	num = len(txArray)
	// }
	TxCollection := Database.C("txs")
	chs := make([]chan *s.Tx, len(txArray))
	//	startTime := time.Now()
	//fmt.Println("Multirun start")
	for i, txid := range txArray {
		//go run(i, sleeptime, ch)
		//fmt.Println("asd")

		chs[i] = make(chan *s.Tx)
		go runstr(txid, chs[i], Database)

	}
	//for range input {
	//	fmt.Println(<-ch)
	//}
	for _, ch := range chs {
		TxCollection.Insert(<-ch)
	}

	//endTime := time.Now()
	//fmt.Printf("Multissh finished. Process time %s. Number of tasks is %d", endTime.Sub(startTime), len(txArray))
}
func GetTxArray() string {
	s.GetMongo(mongourl)
	Database := s.GlobalS.DB("BTC")
	blockCollection := Database.C("blocks")
	//start, end := btc.CalaulateTime(blockCollection)
	for i := 520001; i <= 520002; i++ {
		//for i := start; i <= end; i++ {
		hash := btc.GetbtcHashRPC(int64(i))
		blocks := btc.GetBlocks(hash)
		log.Println("Process Block With Height: ", i, "; And Blocks Hash :", hash)
		blockCollection.Insert(blocks)
		txArray := blocks.Tx
		//	fmt.Println(txArray[0])
		Character(txArray, Database)
		// result := btc.CatchUpTx(txArray, Database)
		// if result != true {
		// 	return "Failed"
		// }
	}
	return "Success"
}
func CatchTx(txid string, Database *mgo.Database) *s.Tx {
	TxCollection := Database.C("txs")
	//var Time uint64
	//for _, k := range txidArray {
	var q bson.M
	res := TxCollection.Find(bson.M{"txid": txid}).One(&q)
	//fmt.Println(txid)
	if res != nil {
		result, _ := btc.GetClearTx(txid)

		//TxCollection.Insert(result)
		//	}
		return result
	}

	return nil
}
