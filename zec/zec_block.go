package zec

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

const (
	URL       = s.LTCURL
	GenesisTx = s.LTCGenesisTx
	mongourl  = s.Mongourl
	apikey    = s.ZECAPIKEY
)

func GetZECCountRPC() int64 {
	res, err := CallZECRPC(URL, "getblockcount", 1, []interface{}{})
	if err != nil {
		fmt.Println("Error")
	}
	//fmt.Println(res)
	//	fmt.Println("======================4")
	count, _ := res["result"].(json.Number).Int64()
	return count
}

func GetZECHashRPC(height int64) string {
	res, err := CallZECRPC(URL, "getblockhash", 1, []interface{}{height})
	if err != nil {
		fmt.Println("error")
	}
	tar := res["result"].(string)
	//	fmt.Println(tar)
	return tar
}

func GetBlocks(hash string) s.Blocks {
	res, err := CallZECRPC(URL, "getblock", 1, []interface{}{hash})
	if err != nil {
		fmt.Println("Error")
	}
	rawInfo := res["result"].(map[string]interface{})
	//	fmt.Println(rawInfo)
	data, _ := json.Marshal(rawInfo)
	var ZECd s.Blocks
	json.Unmarshal(data, &ZECd)
	ZECd.NTx = len(ZECd.Tx)
	Difficulty, _ := rawInfo["difficulty"].(json.Number).Float64()
	ZECd.Difficulty = uint64(Difficulty)
	return ZECd

}

func CallZECRPC(address string, method string, id interface{}, params []interface{}) (map[string]interface{}, error) {
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
	blockCollection.Find(bson.M{}).Sort("height").Limit(1).One(&target)
	//blockCollection.Remove(target)
	startheight := int64(target.Height)
	endheight := GetZECCountRPC()
	return startheight, endheight
}

func CatchUpBlocks() string {
	s.GetMongo(mongourl)
	Database := s.GlobalS.DB("ZEC")
	blockIndex := mgo.Index{
		Key:    []string{"hash", "height", "time"},
		Unique: false,
	}

	blockCollection := Database.C("blocks")
	blockCollection.EnsureIndex(blockIndex)
	start, end := CalaulateTime(blockCollection)
	//for i := 520000; i <= 520001; i++ {
	if start != end {
		blockCollection.Remove(bson.M{"height": start})
		for i := start; i <= end; i++ {
			hash := GetZECHashRPC(int64(i))
			blocks := GetBlocks(hash)
			log.Println("Process Block With Height: ", i, "; And Blocks Hash :", hash)
			blockCollection.Insert(blocks)
		}
	} else {
		return "Success"
	}
	return "Success"
}
