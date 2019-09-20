package getdata

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/GGBTC/explorer/pkgs"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GetData interface {
	GetConnection(dataName string, collectionName string)
	GetBlockInfoById(key string, value string) (blockInfo blockInfo)
	GetTransCationInfoByHash(value string) (transCationInfo TransCationInfo)
	GetAllTransCationData() (transCationInfo []TransCationInfo)
	FindOutByAddress(address string) (result OutAccount)
	//GetTransCationInfoAddress(value string) (transCationInfo []TransCationInfo)
	InsertOut(info OutAccount)
	UpdateOutAccount(str string, info OutAccount)
	InsertErr(err string, reason interface{})
	GetAccountInfoByAddress(address string) (outAccount OutAccount)
	GetRecentTransCation() (transCationInfo []TransCationInfo)
	GetRecentBlock() (Block []pkgs.ClearBlock)
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
	Vout          []Vout `bson:"vout"`
	Confirmations int    `bson:"confirmations"`
	Blocktime     int    `bson:"blocktime"`
	Hex           string `bson:"hex"`
	Locktime      int    `bson:"locktime"`
	Size          int    `bson:"size"`
	Vsize         int    `bson:"vsize"`
	Version       int    `bson:"version"`
	Vin           []Vin  `bson:"vin"`
	Blockhash     string `bson:"blockhash"`
	Time          int    `bson:"time"`
	Txid          string `bson:"txid"`
	Hash          string `bson:"hash"`
}
type ScriptPubKey struct {
	Asm       string   `bson:"asm"`
	Hex       string   `bson:"hex"`
	ReqSigs   int      `bson:"reqSigs"`
	Type      string   `bson:"type"`
	Addresses []string `bson:"addresses"`
}
type Vout struct {
	Value        float64      `bson:"value"`
	N            int          `bson:"n"`
	ScriptPubKey ScriptPubKey `bson:"scriptPubKey"`
}
type Vin struct {
	Coinbase string `bson:"coinbase"`
	Sequence int64  `bson:"sequence"`
	Txid     string `bson:"txid"`
}

type OutAccount struct {
	Address        string         `bson:"Address"`
	TotalOutAmount float64        `bson:"TotalOutAmount"`
	TotalOutTimes  int64          `bson:"TotalOutTimes"`
	FirstSeen      int64          `bson:"FirstSeen"`
	LastSeen       int64          `bson:"LastSeen"`
	TransInfo      []TransOutInfo `bson:"TransOutInfo"`
}
type TransOutInfo struct {
	TransHash string  `bson:"TransHash"`
	Value     float64 `bson:"Value"`
}

type AccountInfo struct {
	Result    Result `bson:"result"`
	Detail    Detail `bson:"detail"`
	Firstseen int    `bson:"firstseen"`
	Lastseen  int    `bson:"lastseen"`
}
type Txdetails struct {
	Index int    `bson:"index"`
	Spent bool   `bson:"spent"`
	Txid  string `bson:"txid"`
	Value int64  `bson:"value"`
}
type Result struct {
	Address   string      `bson:"address"`
	Txdetails []Txdetails `bson:"txdetails"`
}
type Detail struct {
	Totalreccount  int   `bson:"totalreccount"`
	Totalsentcount int   `bson:"totalsentcount"`
	Totalreceived  int64 `bson:"totalreceived"`
	Totalsent      int   `bson:"totalsent"`
	Balance        int64 `bson:"balance"`
}

const (
	URL = "192.168.3.16:27017" //连接mongoDB启动服务的端口号 你得先启动mongoDB服务
)

var client *mgo.Collection
var session *mgo.Session

type err struct {
	ErrTransCationId string      `bson:"errTransCationId"`
	ErrReason        interface{} `bson:"errReason"`
}
type Mongo struct {
}

func (mongo Mongo) GetConnection(dataName string, collectionName string) {
	session, _ = mgo.Dial(URL) //连接数据库
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dataName)    //数据库名称
	client = db.C(collectionName) //如果该集合已经存在的话，则直接返回
}
func (mongo Mongo) GetBlockInfoByKeyValue(key string, value string) (block pkgs.ClearBlock) {
	blockResult := pkgs.ClearBlock{}
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
func (mongo Mongo) GetTransCationInfoByHash(value string) (transCationInfo pkgs.Tx) {
	var res pkgs.Tx
	client.Find(bson.M{"txid": value}).One(&res)
	fmt.Println(res)
	return res
}

// func (mongo Mongo) GetTransCationInfoAddress(value string) (transCationInfo []TransCationInfo) {
// 	var res []TransCationInfo
// 	client.Find(bson.M{"$or": []bson.M{bson.M{"in.prev_out.address": value}, bson.M{"out.prev_out.address": value}}}).All(&res)
// 	return res
// }
func (mongo Mongo) GetAllTransCationData() (transCationInfo []pkgs.Tx) {
	var res []pkgs.Tx
	client.Find(bson.M{}).All(&res)
	return res
}

func (mongo Mongo) FindOutByAddress(address string) (result pkgs.Address) {
	var res pkgs.Address
	client.Find(bson.M{"address": address}).One(&res)
	return res
}
func (mongo Mongo) InsertOut(info pkgs.Address) {
	client.Insert(info)
}

func (mongo Mongo) InsertErr(arr string, reason interface{}) {
	var errObj err
	errObj.ErrReason = reason
	errObj.ErrTransCationId = arr
	client.Insert(errObj)
}

func (mongo Mongo) UpdateOutAccount(str string, info OutAccount) {
	client.Update(bson.M{"address": str}, bson.M{"$set": bson.M{"TransOutInfo": info.TransInfo, "TotalOutTimes": info.TotalOutTimes, "TotalOutAmount": info.TotalOutAmount}})
}

func (mongo Mongo) GetAccountInfoByAddress(address string) (outAccount AccountInfo) {
	var res AccountInfo
	client.Find(bson.M{"result.address": address}).One(&res)
	return res
}
func (mongo Mongo) GetBlockHeight() int {
	//var res AccountInfo
	Number, _ := client.Count()
	res := Number - 1
	return res
}
func (mongo Mongo) GetUnSpent(address string) []pkgs.Unspent {
	// var res pkgs.Address
	// var Unspent []pkgs.TxOut
	// client.Find(bson.M{"result.address": address}).One(&res)
	// asd := res.Result.Txdetails
	// for _, k := range asd {
	// 	if k.Spent == false {
	// 		Unspent = append(Unspent, k)
	// 	}
	// }
	// return Unspent
	//var res pkgs.Unspent
	var Unspent []pkgs.Unspent
	client.Find(bson.M{"address": address}).All(&Unspent)
	return Unspent
}
func (mongo Mongo) GetRecentTransCation() (transCationInfo []pkgs.Tx) {
	var res []pkgs.Tx
	client.Find(bson.M{}).Sort("-blocktime").Limit(1000).All(&res)
	fmt.Println(len(res))
	return res
}

func (mongo Mongo) GetRecentBlock() (Block []pkgs.ClearBlock) {
	var res []pkgs.ClearBlock
	client.Find(bson.M{}).Sort("-height").Limit(1000).All(&res)
	fmt.Println(len(res))
	return res
}
