package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GGBTC/explorer/btc"
	zec "github.com/GGBTC/explorer/zcash"
)

func main() {
	//Omni
	//txid := "6fa242458251959ce83ccc5e1e55527eae93ba9fa40410cf05cc507dbccfa014"
	//BTC
	//txid := "2ad3bb53f75ddc5de92be3550e4759d17c5c8851f6b750f685d881d90faf8851"
	//CoinBase
	txid := "2e9de4a044dd584527c46100c206ec709e810a79993bd7c6807ff94107b2636d"
	//Complex
	//txid := "55fc58508ec64d5e907bc0b3eb9918fc4bf6b31fb87ca6e830f9e0f00ab62779"
	// res_tx, err := CallZECRPC1(service.ZECURL, "getrawtransaction", 1, []interface{}{txid, 1})
	// if err != nil {
	// 	log.Fatalf("Err: %v", err)
	// }
	// fmt.Println(res_tx)
	res, _ := btc.GetClearTx(txid)
	data, _ := json.Marshal(res)
	fmt.Println(string(data))
}

// func main() {
// 	//pkgs.CatchUpBlocks()
// 	//	service.GetMongo(mongourl)
// 	running := true
// 	for running {
// 		KeepRunning()
// 	}
// }
func KeepRunning() {
	res := zec.CatchUpBlockss()
	if res == "Success" {
		time.Sleep(time.Duration(60) * time.Second)
	}
	//return true
}
func CallZECRPC1(address string, method string, id interface{}, params []interface{}) (interface{}, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return nil, err
	}
	resp, err := http.Post(address, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	//var result map[string]interface{}
	var result interface{}
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
