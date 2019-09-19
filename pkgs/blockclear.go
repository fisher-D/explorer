package pkgs

import (
	"encoding/json"
	"fmt"

	"github.com/GGBTC/explorer/service"
)

func GetClearBlock(hash string) ClearBlock {
	//hash := "000000000000e5e0da97bcbda0cdec1ac1d15761afb76e42e13ea53edd65cba4"
	res, err := service.CallBitcoinRPC(URL, "getblock", 1, []interface{}{hash})
	if err != nil {
		fmt.Println("Error")
	}
	//fmt.Println(res)
	tar := res["result"].(map[string]interface{})
	//fmt.Println(tar)
	//jsoninfo := res["result"].(map[string]interface{})
	var blockjson ClearBlock
	blockjson.Hash = tar["hash"].(string)
	Height, _ := tar["height"].(json.Number).Int64()
	blockjson.Height = uint32(Height)
	//	blockjson.Bits = tar["bits"].(string)           //string
	//	blockjson.Chainwork = tar["chainwork"].(string) //string
	//Confirmations, _ := tar["confirmations"].(json.Number).Int64()
	//blockjson.Confirmations = uint32(Confirmations)                     //uint32
	blockjson.Difficulty, _ = tar["difficulty"].(json.Number).Float64() //float64
	//Mediantime, _ := tar["mediatime"].(int64)
	//blockjson.Mediantime = uint32(Mediantime)         //uint32
	//	blockjson.Merkleroot = tar["merkleroot"].(string) //string
	Tx := tar["tx"].([]interface{})
	blockjson.NTx = uint32(len(Tx))
	// if tar["nextblockhash"] != nil { //uint32
	// 	blockjson.Nextblockhash = tar["nextblockhash"].(string)
	// } else {

	//	} //string
	//	Nonce, _ := tar["nonce"].(json.Number).Int64()
	// blockjson.Nonce = uint64(Nonce)
	// if tar["previousblockhash"] != nil { //uint3264
	// 	blockjson.Previousblockhash = tar["previousblockhash"].(string)
	// } else {
	// 	blockjson.Previousblockhash = "N/A"
	// } //string
	//	Size, _ := tar["size"].(json.Number).Int64()
	//blockjson.Size = uint32(Size) //uint32
	//Strippedsize, _ := tar["strippedsize"].(json.Number).Int64()
	//	blockjson.Strippedsize = uint32(Strippedsize) //uint32
	//	Time, _ := tar["time"].(json.Number).Int64()
	//	blockjson.Time = uint32(Time) //uint32

	var str []string
	txs, err := json.Marshal(Tx)
	if err != nil {
		fmt.Println("ERROR")
	}
	json.Unmarshal(txs, &str)
	blockjson.Tx = str //[]string
	Version, _ := tar["version"].(json.Number).Int64()
	blockjson.Version = uint32(Version) //uint32
	//	blockjson.VersionHex = tar["versionHex"].(string) //string
	//Weight, _ := tar["weight"].(json.Number).Int64()
	//	blockjson.Weight = uint32(Weight) //uint32
	// data, _ := json.Marshal(blockjson)
	// fmt.Println(string(data))
	return blockjson
}

type ClearBlock struct {
	Hash string `bson:"hash"`
	//	Confirmations     uint32   `bson:"confirmations"`
	//	Strippedsize      uint32   `bson:"strippedsize"`
	//	Size              uint32   `bson:"size"`
	//	Weight            uint32   `bson:"weight"`
	Height  uint32 `bson:"height"`
	Version uint32 `bson:"version"`
	//	VersionHex        string   `bson:"versionHex"`
	//	Merkleroot string   `bson:"merkleroot"`

	//Time              uint32   `bson:"time"`
	//Mediantime        uint32   `bson:"mediantime"`
	//Nonce             uint64   `bson:"nonce"`
	//Bits              string   `bson:"bits"`
	Difficulty float64 `bson:"difficulty"`
	//	Chainwork  string  `bson:"chainwork"`
	NTx uint32   `bson:"nTx"`
	Tx  []string `bson:"tx"`
	//	Previousblockhash string   `bson:"previousblockhash"`
	//	Nextblockhash     string   `bson:"nextblockhash"`
}

// type BlockSample struct {
// 	Hash       string   `json:"Hash"`
// 	Height     int      `json:"Height"`
// 	Version    int      `json:"Version"`
// 	Difficulty float64  `json:"Difficulty"`
// 	NTx        int      `json:"NTx"`
// 	Tx         []string `json:"transaction"`
// }
