package btcd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"

	s "github.com/GGBTC/explorer/service"
)

const URL = s.BTCURL
const GenesisTx = s.BTCGenesisTx
const mongourl = s.Mongourl

//mongourl = "192.168.72.250:27017"
const apikey = s.ZECAPIKEY

func GetBTCDCountRPC() int64 {
	res, err := CallBTCDRPC(URL, "getblockcount", 1, []interface{}{})
	if err != nil {
		fmt.Println("Error")
	}
	count, _ := res["result"].(json.Number).Int64()
	return count
}

func GetBTCDHashRPC(height int64) string {
	res, err := CallBTCDRPC(URL, "getblockhash", 1, []interface{}{height})
	if err != nil {
		fmt.Println("error")
	}
	tar := res["result"].(string)
	return tar
}

func GetBlocks(hash string) s.Blocks {
	res, err := CallBTCDRPC(URL, "getblock", 1, []interface{}{hash})
	if err != nil {
		fmt.Println("Error")
	}
	rawInfo := res["result"].(map[string]interface{})
	data, _ := json.Marshal(rawInfo)
	var BTCDd s.Blocks
	err =json.Unmarshal(data, &BTCDd)
	if err!=nil{
		return BTCDd
	}
	BTCDd.NTx = len(BTCDd.Tx)
	var Difficulty, _ = rawInfo["difficulty"].(json.Number).Float64()
	BTCDd.Difficulty = uint64(Difficulty)
	return BTCDd

}

func CallBTCDRPC(address string, method string, id interface{}, params []interface{}) (map[string]interface{}, error) {
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
	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}
	return result, nil
}

func CalaulateTime(blockCollection *mgo.Collection) (int64, int64) {
	var target s.Blocks
	blockCollection.Find(bson.M{}).Sort("-height").Limit(1).One(&target)
	startheight := int64(target.Height)
	endheight := GetBTCDCountRPC()
	return startheight, endheight
}

func CatchUpBlocks() string {
	s.GetMongo(mongourl)
	Database := s.GlobalS.DB("BTCD")
	blockCollection := Database.C("blocks")
	//使用单一检索，以简便查询
	blockIndex1 := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}
	blockIndex2 := mgo.Index{
		Key:    []string{"-height"},
		Unique: true,
	}
	blockIndex3 := mgo.Index{
		Key:    []string{"-time"},
		Unique: false,
	}

	blockCollection.EnsureIndex(blockIndex1)
	blockCollection.EnsureIndex(blockIndex2)
	blockCollection.EnsureIndex(blockIndex3)
	start, end := CalaulateTime(blockCollection)
	//for i := 520000; i <= 520001; i++ {
	chlimint := make(chan bool, 5)
	//orderlimint := make(chan int64, 5)
	if start != end {
		blockCollection.Remove(bson.M{"height": start})
		for i := start; i <= end; i++ {
			chlimint <- true
			//orderlimint <- i
			go func(i int64) {
				hash := GetBTCDHashRPC(i)
				blocks := GetBlocks(hash)
				log.Println("Process Block With Height: ", i, "; And Blocks Hash :", hash)
				blockCollection.Insert(blocks)
				<-chlimint
			}(i)
		}
	} else {
		return "Success"
	}
	return "Success"
}
