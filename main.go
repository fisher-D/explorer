package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GGBTC/explorer/service"

	zec "github.com/GGBTC/explorer/zcash"
)

func main() {
	//3
	//txid := "9a4adaf3953818eb1634407032db0e00ef2441c49c1364161411d0743ec1a939"
	//2
	//txid := "8974d08d1c5f9c860d8b629d582a56659a4a1dcb2b5f98a25a5afcc2a784b0f4"
	//txid := "ae61ee40c37fdc05468786cf6e83a98741a214c41fc5d62057f768d9b0af769e"
	txid := "a311f80d7e54f4d1d1c3e0ea475ad7ed709818a91301939f474cbde72ba0f38b"
	res_tx, err := CallZECRPC1(service.ZECURL, "getrawtransaction", 1, []interface{}{txid, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	fmt.Println(res_tx)
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
