package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/GGBTC/explorer/pkgs"

	getdata "github.com/GGBTC/explorer/api/getData"

	"github.com/gorilla/mux"
)

// 根据id获取block
//Auto Handle Hash or Height Done
func GetBlockInfoByKeyValue(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	params := mux.Vars(req)
	Id := params["id"]
	coinName := params["coinName"]
	if coinName == "btc" {
		key := btcJudgeBlockKey(Id)
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "blocks")
		result := mongo.GetBlockInfoByKeyValue(key, Id)
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		key := btcJudgeBlockKey(Id)
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "blocks")
		result := mongo.GetBlockInfoByKeyValue(key, Id)
		json.NewEncoder(w).Encode(result)
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

// 根据id获取transcation Done
func GetTranscationById(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	params := mux.Vars(req)
	Id := params["id"]
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "txs")
		if len(Id) == 64 {
			result := mongo.GetTransCationInfoByHash(Id)
			json.NewEncoder(w).Encode(result)
		}
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "txs")
		if len(Id) == 64 {
			result := mongo.GetTransCationInfoByHash(Id)
			json.NewEncoder(w).Encode(result)
		}
	}
}

//Done
func GetAccountInfoByAddress(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	params := mux.Vars(req)
	address := params["id"]
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "address")
		result := mongo.GetAccountInfoByAddress(address)
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "address")
		result := mongo.GetAccountInfoByAddress(address)
		json.NewEncoder(w).Encode(result)
	}
}

//Done
func GetLastBlock(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//params := mux.Vars(req)
	params := mux.Vars(req)
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "blocks")
		result := mongo.GetCountNumber()
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "blocks")
		result := mongo.GetCountNumber()
		json.NewEncoder(w).Encode(result)
	}
}
func GetLastTx(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//params := mux.Vars(req)
	params := mux.Vars(req)
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "txs")
		result := mongo.GetCountNumber()
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "txs")
		result := mongo.GetCountNumber()
		json.NewEncoder(w).Encode(result)
	}
}
func GetInformation(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//params := mux.Vars(req)
	params := mux.Vars(req)
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "info")
		result := mongo.Getinfo()
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "info")
		result := mongo.Getinfo()
		json.NewEncoder(w).Encode(result)
	}
}
func GetUnSpnetNumber(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//params := mux.Vars(req)
	params := mux.Vars(req)
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "utxo")
		result := mongo.GetCountNumber()
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "utxo")
		result := mongo.GetCountNumber()
		json.NewEncoder(w).Encode(result)
	}
}

//Done
func GetAddressUnSpent(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	params := mux.Vars(req)
	coinName := params["coinName"]
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		address := params["id"]
		mongo.GetConnection("BTC", "utxo")
		result := mongo.GetUnSpent(address)
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		address := params["id"]
		mongo.GetConnection("LTC", "utxo")
		result := mongo.GetUnSpent(address)
		json.NewEncoder(w).Encode(result)
	}
}

//Done
func GetRecentTranscation(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	params := mux.Vars(req)
	coinName := params["coinName"]
	page := params["page"]
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println(err)
	}
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "txs")
		result := mongo.GetRecentTransCation(pageNum)
		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "txs")
		result := mongo.GetRecentTransCation(pageNum)
		json.NewEncoder(w).Encode(result)
	}
}

//Done
func GetRecentBlocks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	params := mux.Vars(req)
	coinName := params["coinName"]
	page := params["page"]
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println(err)
	}
	if coinName == "btc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("BTC", "blocks")
		result := mongo.GetRecentBlock(pageNum)

		json.NewEncoder(w).Encode(result)
	} else if coinName == "ltc" {
		mongo := getdata.Mongo{}
		mongo.GetConnection("LTC", "blocks")
		result := mongo.GetRecentBlock(pageNum)
		json.NewEncoder(w).Encode(result)
	}
}
func GetDataBaseInfo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	result := pkgs.DataName()
	json.NewEncoder(w).Encode(result)
}

func main() {
	//	Get handle function:
	router := mux.NewRouter()
	//Done
	router.HandleFunc("/block/{coinName}/{id}", GetBlockInfoByKeyValue).Methods("GET")
	//Done
	router.HandleFunc("/tx/{coinName}/{id}", GetTranscationById).Methods("GET")
	//Done
	router.HandleFunc("/address/{coinName}/{id}", GetAccountInfoByAddress).Methods("GET")
	//Done
	router.HandleFunc("/unspent/{coinName}/{id}", GetAddressUnSpent).Methods("GET")
	//Done
	router.HandleFunc("/latestblock/{coinName}", GetLastBlock).Methods("GET")
	router.HandleFunc("/latesttx/{coinName}", GetLastTx).Methods("GET")
	router.HandleFunc("/latestutxo/{coinName}", GetUnSpnetNumber).Methods("GET")
	router.HandleFunc("/blockinfo/{coinName}", GetInformation).Methods("GET")
	router.HandleFunc("/database", GetDataBaseInfo).Methods("GET")
	//Done
	router.HandleFunc("/recent/tx/{coinName}/{page}", GetRecentTranscation).Methods("GET")
	//Done
	router.HandleFunc("/recent/blocks/{coinName}/{page}", GetRecentBlocks).Methods("GET")
	log.Fatal(http.ListenAndServe(":9899", router))
}
