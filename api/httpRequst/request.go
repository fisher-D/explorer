package getdata

import (
	"fmt"
	"regexp"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GetData interface {
	GetConnection(dataName string, collectionName string)
	GetBlockInfoById(key string, value string) (blockInfo blockInfo)
	GetTransCationInfoByHash(value string) (transCationInfo TransCationInfo)
	GetTransCationInfoAddress(value string) (transCationInfo []TransCationInfo)
}

type blockInfo struct {
	Hash          string   `bson:"hash"`
	Confirmations int      `bson:"confirmations"`
	Strippedsize  int      `bson:"strippedsize"`
	Size          int      `bson:"size"`
	Weight        int      `bson:"weight"`
	Height        int      `bson:"height"`
	Version       int      `bson:"version"`
	VersionHex    string   `bson:"versionHex"`
	Merkleroot    string   `bson:"merkleroot"`
	Tx            []string `bson:"tx"`
	Time          int      `bson:"time"`
	Mediantime    int      `bson:"mediantime"`
	Nonce         struct {
		NumberLong string `bson:"$numberLong"`
	} `bson:"nonce"`
	Bits              string  `bson:"bits"`
	Difficulty        float64 `bson:"difficulty"`
	Chainwork         string  `bson:"chainwork"`
	NTx               int     `bson:"nTx"`
	Previousblockhash string  `bson:"previousblockhash"`
	Nextblockhash     string  `bson:"nextblockhash"`
}

type TransCationInfo struct {
	Hash            string    `bson:"hash"`
	Size            int       `bson:"size"`
	LockTime        int       `bson:"lock_time"`
	Ver             int       `bson:"ver"`
	VinSz           int       `bson:"vin_sz"`
	VoutSz          int       `bson:"vout_sz"`
	In              []In      `bson:"in"`
	Out             []Out     `bson:"out"`
	VoutTotal       VoutTotal `bson:"vout_total"`
	VinTotal        VinTotal  `bson:"vin_total"`
	BlockHash       string    `bson:"block_hash"`
	BlockHeight     int       `bson:"block_height"`
	BlockTime       int       `bson:"block_time"`
	FirstSeenTime   int       `bson:"first_seen_time"`
	FirstSeenHeight int       `bson:"first_seen_height"`
}
type Value struct {
	NumberLong string `bson:"$numberLong"`
}
type PrevOut struct {
	Hash    string `bson:"hash"`
	N       int    `bson:"n"`
	Address string `bson:"address"`
	Value   Value  `bson:"value"`
}
type In struct {
	TxHash    string  `bson:"txHash"`
	BlockHash string  `bson:"blockHash"`
	BlockTime int     `bson:"blockTime"`
	PrevOut   PrevOut `bson:"prev_out"`
	N         int     `bson:"n"`
}
type Spent struct {
	Spent bool `bson:"spent"`
}
type Out struct {
	Hash  string `bson:"hash"`
	Value Value  `bson:"value"`
	N     int    `bson:"n"`
	Spent Spent  `bson:"spent"`
}
type VoutTotal struct {
	NumberLong string `bson:"$numberLong"`
}
type VinTotal struct {
	NumberLong string `bson:"$numberLong"`
}

const (
	URL = "localhost:27017" //连接mongoDB启动服务的端口号 你得先启动mongoDB服务
)

var client *mgo.Collection
var session *mgo.Session

type Mongo struct {
}

func (mongo Mongo) GetConnection(dataName string, collectionName string) {
	session, _ = mgo.Dial(URL) //连接数据库
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dataName)    //数据库名称
	client = db.C(collectionName) //如果该集合已经存在的话，则直接返回
	fmt.Println("clinet =", client)
}
func (mongo Mongo) GetBlockInfoByKeyValue(key string, value string) (block blockInfo) {
	blockResult := blockInfo{}
	pattern := "^[0-9]*$"
	boolRes, _ := regexp.MatchString(pattern, value)
	if boolRes == true {
		val, _ := strconv.Atoi(value)
		client.Find(bson.M{key: val}).One(&blockResult)
		return blockResult
	} else {
		client.Find(bson.M{key: value}).One(&blockResult)
		return blockResult
	}
}

func (mongo Mongo) GetTransCationInfoByHash(value string) (transCationInfo TransCationInfo) {
	var res TransCationInfo
	client.Find(bson.M{"hash": value}).One(&res)
	fmt.Println(res)
	return res
}

func (mongo Mongo) GetTransCationInfoAddress(value string) (transCationInfo []TransCationInfo) {
	var res []TransCationInfo
	client.Find(bson.M{"$or": []bson.M{bson.M{"in.prev_out.address": value}, bson.M{"out.prev_out.address": value}}}).All(&res)
	return res
}
