package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	getdata "github.com/GGBTC/explorer/api/getData"

	"github.com/gorilla/mux"
)

// 根据id获取block
//Auto Handle Hash or Height
func GetBlockInfoByKeyValue(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	Id := params["id"]
	coinName := params["coinName"]
	if coinName == "btc" {
		key := btcJudgeBlockKey(Id)
		mongo := getdata.Mongo{}
		mongo.GetConnection("GGBTC", "blocks")
		result := mongo.GetBlockInfoByKeyValue(key, Id)
		json.NewEncoder(w).Encode(result)
	} else if coinName == "Eth" {
		fmt.Println("------------------------", coinName)
	}
}

func btcJudgeBlockKey(id string) (key string) {
	pattern := "^[0-9]*$"
	result, _ := regexp.MatchString(pattern, id)
	if result == true {
		return "height"
	} else {
		return "hash"
	}
}

// 根据id获取transcation
func GetTranscationById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	Id := params["id"]
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("GGBTC", "transaction")
		if len(Id) == 64 {
			result := mongo.GetTransCationInfoByHash(Id)
			json.NewEncoder(w).Encode(result)
		}
	} else if coinName == "Eth" {
		fmt.Println("------------------------", coinName)
	}
}

func GetAccountInfoByAddress(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	address := params["id"]
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("GGBTC", "address")
		result := mongo.GetAccountInfoByAddress(address)
		json.NewEncoder(w).Encode(result)
	} else if coinName == "eth" {
		fmt.Println("-----------------------------", coinName)
	}
}
func GetLastBlock(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	mongo := getdata.Mongo{}
	mongo.GetConnection("GGBTC", "blocks")
	result := mongo.GetBlockHeight()
	json.NewEncoder(w).Encode(result)
}
func GetAddressUnSpent(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	mongo := getdata.Mongo{}
	address := params["id"]
	mongo.GetConnection("GGBTC", "unspent")
	result := mongo.GetUnSpent(address)
	json.NewEncoder(w).Encode(result)
}

func GetRecentTranscation(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("hahahahahaha")
	params := mux.Vars(req)
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("GGBTC", "transaction")
		result := mongo.GetRecentTransCation()

		json.NewEncoder(w).Encode(result)
	} else if coinName == "eth" {
		fmt.Println("-----------------------------", coinName)
	}
}
func GetRecentBlocks(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("")
	params := mux.Vars(req)
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("GGBTC", "blocks")
		result := mongo.GetRecentBlock()

		json.NewEncoder(w).Encode(result)
	} else if coinName == "eth" {
		fmt.Println("-----------------------------", coinName)
	}
}
func main() {
	//	Get handle function:
	router := mux.NewRouter()
	router.HandleFunc("/block/{coinName}/{id}", GetBlockInfoByKeyValue).Methods("GET")
	router.HandleFunc("/tx/{coinName}/{id}", GetTranscationById).Methods("GET")
	router.HandleFunc("/address/{coinName}/{id}", GetAccountInfoByAddress).Methods("GET")
	router.HandleFunc("/unspent/{coinName}/{id}", GetAddressUnSpent).Methods("GET")
	router.HandleFunc("/latestblock/{coinName}", GetLastBlock).Methods("GET")
	router.HandleFunc("/recent/tx/{coinName}", GetRecentTranscation).Methods("GET")
	router.HandleFunc("/recent/blocks/{coinName}", GetRecentBlocks).Methods("GET")
	log.Fatal(http.ListenAndServe(":9899", router))
	fmt.Println("1")
}
