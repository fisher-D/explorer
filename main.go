package main

import (
	"log"
	"sync"

	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	Cop()
}
func Cop() {
	url := "192.168.72.250:27017"
	//url := "localhost:27017"
	s.GetMongo(url)
	Collection := s.GlobalS.DB("LTC").C("blocks")
	//endtime, _ := Collection.Count()
	var target s.Blocks
	Collection.Find(bson.M{}).Sort("-height").Limit(1).One(&target)
	endtime := target.Height
	starttime := 0
	var blocks []s.Blocks
	//var l sync.Mutex
	var wg sync.WaitGroup
	ch := make(chan int, endtime)
	chlimit := make(chan bool, 5)
	dup := func(height int) {
		defer wg.Done()
		Collection.Find(bson.M{"height": height}).All(&blocks)
		number := len(blocks)
		log.Println(number)
		if number > 1 {
			target := blocks[0]
			Collection.RemoveAll(bson.M{"height": height})
			Collection.Insert(target)
			log.Println("Success", height)
		}
		log.Println("Good Record", height)
		Collection.Find(bson.M{"height": height}).All(&blocks)
		count := len(blocks)
		log.Println("Rest Record Number:", count)
		<-chlimit
	}

	for i := starttime; i < endtime; i++ {
		ch <- i
		chlimit <- true
		wg.Add(1)
		go dup(<-ch)
		wg.Wait()
	}
	log.Println("Finish")

}

//联合 需要对应位置
//没有 3.44 22345
//非联合 0.004
func Test() {
	url := "localhost:27017"
	s.GetMongo(url)
	Collection := s.GlobalS.DB("LTC").C("blocks")
	var target s.Blocks
	Collection.Find(bson.M{"height": 124}).One(&target)
	Collection.RemoveAll(bson.M{"height": 124})
	Collection.Insert(target)
}
