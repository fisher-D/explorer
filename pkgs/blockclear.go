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

	blockjson.Difficulty, _ = tar["difficulty"].(json.Number).Float64() //float64

	Tx := tar["tx"].([]interface{})
	blockjson.NTx = uint32(len(Tx))

	var str []string
	Txs, err := json.Marshal(Tx)
	if err != nil {
		fmt.Println("ERROR")
	}
	json.Unmarshal(Txs, &str)
	blockjson.Tx = str //[]string
	Version, _ := tar["version"].(json.Number).Int64()
	blockjson.Version = uint32(Version) //uint32
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
