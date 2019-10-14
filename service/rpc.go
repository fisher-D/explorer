package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

// GenesisTx 。。。
const GenesisTx = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"

// Helper to make call to bitcoind RPC API
func CallBitcoinRPC(address string, method string, id interface{}, params []interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return nil, err
	}
	resp, err := http.Post(address,
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	result := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}
	return result, nil
}

// TODO
// 优化后续代码
//Remove Duplications
func RemoveDups() {
	//url := "192.168.72.250:27017"
	url := "localhost:27017"
	GetMongo(url)
	Collection := GlobalS.DB("LTC").C("blocks")
	//endtime, _ := Collection.Count()
	var target Blocks
	Collection.Find(bson.M{}).Sort("-height").Limit(1).One(&target)
	endtime := target.Height
	starttime := 0
	var blocks []Blocks
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
