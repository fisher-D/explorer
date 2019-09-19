package service

import (
	"io/ioutil"

	"gopkg.in/mgo.v2/bson"
)

// Struct holding our configuration
type Config struct {
	BitcoindBlocksPath string `bson:"bitcoind_blocks_path"`
	BitcoindRpcUrl     string `bson:"bitcoind_rpc_url"`
	SsdbHost           string `bson:"ssdb_host"`
	RedisHost          string `bson:"redis_host"`
	LevelDbPath        string `bson:"leveldb_path"`
	AppUrl             string `bson:"app_url"`
	AppPort            uint32 `bson:"app_port"`
	AppApiRateLimited  bool   `bson:"app_api_rate_limited"`
	AppTemplatesPath   string `bson:"app_templates_path"`
	AppGoogleAnalytics string `bson:"app_google_analytics"`
}

// Load configuration from bson file
func LoadConfig(path string) (conf *Config, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	conf = new(Config)
	bson.Unmarshal(file, conf)
	return
}

///
// type AutoTx struct {
// 	ID            string `json:"_id"`
// 	Blockhash     string `json:"blockhash"`
// 	Blocktime     int64  `json:"blocktime"`
// 	Confirmations int    `json:"confirmations"`
// 	Hash          string `json:"hash"`
// 	Hex           string `json:"hex"`
// 	Locktime      int    `json:"locktime"`
// 	Size          uint32 `json:"size"`
// 	Time          int    `json:"time"`
// 	Txid          string `json:"txid"`
// 	Version       uint32 `json:"version"`
// 	Vin           []Vin  `json:"vin"`
// 	Vout          []Vout `json:"vout"`
// 	Vsize         int    `json:"vsize"`
// }
// type ScriptSig struct {
// 	Asm string `json:"asm"`
// 	Hex string `json:"hex"`
// }
// type Vin struct {
// 	ScriptSig ScriptSig `json:"scriptSig"`
// 	Sequence  int64     `json:"sequence"`
// 	Txid      string    `json:"txid"`
// 	Vout      int       `json:"vout"`
// }
// type ScriptPubKey struct {
// 	Addresses []string `json:"addresses"`
// 	Asm       string   `json:"asm"`
// 	Hex       string   `json:"hex"`
// 	ReqSigs   int      `json:"reqSigs"`
// 	Type      string   `json:"type"`
// }
// type Vout struct {
// 	N            int          `json:"n"`
// 	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
// 	Value        float64      `json:"value"`
// }

// type Address struct {
// 	Address string `json:"Address"`
// 	Txs     []Txs  `json:"Txs"`
// }
// type Txs struct {
// 	TxID  string  `json:"TxID"`
// 	Value float64 `json:"value"`
// }
