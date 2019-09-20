package pkgs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	URL       = "http://admin:admin@47.244.98.227:13143"
	GenesisTx = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"
	mongourl  = "192.168.3.16:27017"
)

func GetAndStoreBlock(height int) {
	service.GetMongo(mongourl)
	session := service.GlobalS.DB("GGBTC").C("blocks")
	hash := GetBlockHashRPC(height)
	block := GetClearBlock(hash)
	//fmt.Println("Processing Block||Height :", height, "||Hash :", hash)
	result := SaveBlockRPC(block)
	DrawOutTxs(session, height)
	fmt.Println(time.Now(), result)

}

//GetBlock Done
// func GetBlock(hash string) service.OrigonalBlock {
// 	res, err := service.CallBitcoinRPC(URL, "getblock", 1, []interface{}{hash})
// 	if err != nil {
// 		fmt.Println("Error")
// 	}

// 	tar := res["result"].(map[string]interface{})
// 	//jsoninfo := res["result"].(map[string]interface{})
// 	var blockjson service.OrigonalBlock
// 	blockjson.Hash = tar["hash"].(string)
// 	Height, _ := tar["height"].(json.Number).Int64()
// 	blockjson.Height = uint32(Height)
// 	blockjson.Bits = tar["bits"].(string)           //string
// 	blockjson.Chainwork = tar["chainwork"].(string) //string
// 	Confirmations, _ := tar["confirmations"].(json.Number).Int64()
// 	blockjson.Confirmations = uint32(Confirmations)                     //uint32
// 	blockjson.Difficulty, _ = tar["difficulty"].(json.Number).Float64() //float64
// 	Mediantime, _ := tar["mediatime"].(int64)
// 	blockjson.Mediantime = uint32(Mediantime)         //uint32
// 	blockjson.Merkleroot = tar["merkleroot"].(string) //string
// 	NTx, _ := tar["ntx"].(int64)
// 	blockjson.NTx = uint32(NTx)
// 	if tar["nextblockhash"] != nil { //uint32
// 		blockjson.Nextblockhash = tar["nextblockhash"].(string)
// 	} else {

// 	} //string
// 	Nonce, _ := tar["nonce"].(json.Number).Int64()
// 	blockjson.Nonce = uint64(Nonce)
// 	if tar["previousblockhash"] != nil { //uint3264
// 		blockjson.Previousblockhash = tar["previousblockhash"].(string)
// 	} else {
// 		blockjson.Previousblockhash = "N/A"
// 	} //string
// 	Size, _ := tar["size"].(json.Number).Int64()
// 	blockjson.Size = uint32(Size) //uint32
// 	Strippedsize, _ := tar["strippedsize"].(json.Number).Int64()
// 	blockjson.Strippedsize = uint32(Strippedsize) //uint32
// 	Time, _ := tar["time"].(json.Number).Int64()
// 	blockjson.Time = uint32(Time) //uint32
// 	Tx := tar["tx"].([]interface{})
// 	var str []string
// 	txs, err := json.Marshal(Tx)
// 	if err != nil {
// 		fmt.Println("ERROR")
// 	}
// 	json.Unmarshal(txs, &str)
// 	blockjson.Tx = str //[]string
// 	Version, _ := tar["version"].(json.Number).Int64()
// 	blockjson.Version = uint32(Version)               //uint32
// 	blockjson.VersionHex = tar["versionHex"].(string) //string
// 	Weight, _ := tar["weight"].(json.Number).Int64()
// 	blockjson.Weight = uint32(Weight) //uint32

// 	return blockjson
// }

func SaveBlockRPC(block ClearBlock) string {

	err := service.Insert("GGBTC", "blocks", block)
	if err != nil {
		fmt.Println("ERROR")
	}
	return "Success"
}

func BlockQueryOptions(query interface{}) {

	var session *mgo.Session
	session, err := mgo.Dial("192.168.3.16:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	boll := session.DB("GGBTC").C("blocks")
	var q []bson.M
	boll.Find(query).All(&q)
	for a, info := range q {
		b, err := json.Marshal(info)
		if err != nil {
			panic(err)
		}
		fmt.Println("Result Count :", a+1)
		fmt.Println("Result:", string(b))

	}
}

func QueryBlockByHeight(height int) {
	query := bson.M{"height": height}
	BlockQueryOptions(query)
}

func QueryBlockByHash(hash string) {
	query := bson.M{"hash": hash}
	BlockQueryOptions(query)
}

//CatchUpBlocks Main Functions
func CatchUpBlocks() string {
	var session *mgo.Session
	session, err := mgo.Dial("192.168.3.16:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	BlockCount := GetBlockCountRPC()
	fmt.Println(BlockCount)
	coll := session.DB("GGBTC").C("blocks")
	countNum, err := coll.Count()
	if err != nil {
		panic(err)
	}

	//fmt.Println(countNum)
	fmt.Println("Start from Block Height", countNum)
	//Test Purpose
	for i := countNum; i <= BlockCount; i++ {
		GetAndStoreBlock(i)
	}
	return "Success"
}
