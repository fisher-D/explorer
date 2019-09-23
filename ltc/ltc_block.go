package ltc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	URL       = s.LTCURL
	GenesisTx = s.LTCGenesisTx
	mongourl  = s.Mongourl
)

func GetltcCountRPC() int64 {
	res, err := CallLTCRPC(URL, "getblockcount", 1, []interface{}{})
	if err != nil {
		fmt.Println("Error")
	}
	count, _ := res["result"].(json.Number).Int64()
	return count
}

func GetltcHashRPC(height int64) string {
	res, err := CallLTCRPC(URL, "getblockhash", 1, []interface{}{height})
	if err != nil {
		fmt.Println("error")
	}
	tar := res["result"].(string)
	return tar
}

func GetBlocks(hash string) s.Blocks {
	res, err := s.CallBitcoinRPC(URL, "getblock", 1, []interface{}{hash})
	if err != nil {
		fmt.Println("Error")
	}
	rawInfo := res["result"].(map[string]interface{})
	//	fmt.Println(rawInfo)
	data, _ := json.Marshal(rawInfo)
	var ltcb s.Blocks
	json.Unmarshal(data, &ltcb)
	ltcb.NTx = len(ltcb.Tx)
	return ltcb

}

func CallLTCRPC(address string, method string, id interface{}, params []interface{}) (map[string]interface{}, error) {
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
	//var result map[string]interface{}
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
	//Optional Approach
	//startheight, _ := blockCollection.Count()
	blockCollection.Remove(target)
	startheight := int64(target.Height)
	endheight := GetltcCountRPC()
	return startheight, endheight
}

func CatchUpBlockss() string {
	s.GetMongo(mongourl)
	Database := s.GlobalS.DB("LTC")
	blockCollection := Database.C("blocks")
	//txCollection := Database.C("txs")
	//utxoCollection := Database.C("utxo")
	start, end := CalaulateTime(blockCollection)
	for i := start; i <= end; i++ {
		hash := GetltcHashRPC(i)
		blocks := GetBlocks(hash)
		log.Println("Process Block With Height: ", i, "; And Blocks Hash :", hash)
		blockCollection.Insert(blocks)
		txArray := blocks.Tx
		result := CatchUpTx(txArray, Database)
		if result != true {
			return "Failed"
		}
	}
	return "Success"

}
