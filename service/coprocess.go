package service

import (
	"log"
	"sync"
)

func CoPString(array []string, op func()) {
	var wg sync.WaitGroup
	wg.Add(4)
	jobchstr := genJobstr(array)
	workerPoolstr(4, jobchstr, &wg, op)
	wg.Wait()
}

func CoPNumber(startime, endtime int, op func()) {
	var wg sync.WaitGroup
	wg.Add(4)
	jobCh := genJob(startime, endtime)
	workerPool(4, jobCh, &wg, op)
	wg.Wait()
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
func workerPoolstr(n int, jobChstr <-chan string, wg *sync.WaitGroup, op func()) {
	for i := 0; i < n; i++ {
		go workerstr(i, jobChstr, wg, op)
	}

}

func worker(id int, jobCh <-chan int, wg *sync.WaitGroup, op func()) {
	defer wg.Done()
	for job := range jobCh {
		log.Println("Processer:", id, "  Result", job)
	}
}
func workerstr(id int, jobChstr <-chan string, wg *sync.WaitGroup, op func()) {
	count := 0
	defer wg.Done()
	for k := range jobChstr {
		count++
		res := op
		log.Println(res)
		log.Println("txid:", k, "   Process:", id, "   Times:", count)
	}

}
