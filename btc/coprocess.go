package btc

import (
	"log"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CoPString(array []string, Database *mgo.Database) bool {
	var wg sync.WaitGroup
	wg.Add(4)
	jobchstr := genJobstr(array)
	workerPoolstr(4, jobchstr, &wg, Database)
	wg.Wait()
	return true
}

func genJobstr(array []string) <-chan string {
	jobChstr := make(chan string, 200)
	go func() {
		for _, k := range array {
			jobChstr <- k
		}
		close(jobChstr)
	}()
	return jobChstr
}

func workerPoolstr(n int, jobChstr <-chan string, wg *sync.WaitGroup, Database *mgo.Database) {
	for i := 0; i < n; i++ {
		go workerstr(i, jobChstr, wg, Database)
	}

}

func workerstr(id int, jobChstr <-chan string, wg *sync.WaitGroup, Database *mgo.Database) {
	//count := 0
	defer wg.Done()
	for k := range jobChstr {
		//count++
		//res := op
		//	pkgs.CatchUpTx()
		//log.Println(res)
		CatchUpTx1(k, Database)
		//log.Println("txid:", k, "   Process:", id, "   Times:", count)
		//return res
	}
	//return true
}

func CatchUpTx1(txidArray string, Database *mgo.Database) {
	TxCollection := Database.C("txs")
	var Time uint64
	//for _, k := range txidArray {
	var q bson.M
	res := TxCollection.Find(bson.M{"txid": txidArray}).One(&q)
	if res != nil {
		result, _ := GetClearTx(txidArray)
		TxCollection.Insert(result)
		if result == nil {
			Time = 0
		} else {
			Time = result.BlockTime
		}
		//TODO
		//将拆分的方法并行
		Vin, Vout := BTCUnspent(txidArray, Database)
		GetAddress(Time, Vin, Vout, Database)
	}
}

/////////////////////////////////////////////////////////////
func CoPNumber(startime, endtime int, op func()) {
	var wg sync.WaitGroup
	wg.Add(4)
	jobCh := genJob(startime, endtime)
	workerPool(4, jobCh, &wg, op)
	wg.Wait()
}
func genJob(star int, n int) <-chan int {
	jobCh := make(chan int, 5)
	go func() {
		for i := star; i < n; i++ {
			jobCh <- i
		}
		close(jobCh)
	}()
	return jobCh
}

func workerPool(n int, jobCh <-chan int, wg *sync.WaitGroup, op func()) {
	for i := 0; i < n; i++ {
		go worker(i, jobCh, wg, op)
	}
}

func worker(id int, jobCh <-chan int, wg *sync.WaitGroup, op func()) {
	defer wg.Done()
	for job := range jobCh {
		log.Println("Processer:", id, "  Result", job)
	}
}
