package getdata

import (
	"fmt"
	"regexp"
	"strconv"

	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GetData interface {
	GetConnection(dataName string, collectionName string)
	GetBlockInfoById(key string, value string) (blockInfo s.Blocks)
	GetTransCationInfoByHash(value string) (transCationInfo s.Tx)
	GetAllTransCationData() (transCationInfo []s.Tx)
	FindOutByAddress(address string) (result s.Address)
	//GetTransCationInfoAddress(value string) (transCationInfo []TransCationInfo)
	InsertOut(info s.Address)
	UpdateOutAccount(str string, info s.Address)
	InsertErr(err string, reason interface{})
	GetAccountInfoByAddress(address string) (outAccount s.Address)
	GetRecentTransCation(pageNum int) (transCationInfo []s.Tx)
	GetRecentBlock(pageNum int) (Block []s.Blocks)
}

const (
	URL = s.Mongourl //连接mongoDB启动服务的端口号 你得先启动mongoDB服务
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

func (mongo Mongo) GetBlockInfoByKeyValue(key string, value string) (block s.Blocks) {
	blockResult := s.Blocks{}
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
func (mongo Mongo) GetTransCationInfoByHash(value string) (transCationInfo s.Tx) {
	var res s.Tx
	client.Find(bson.M{"txid": value}).One(&res)
	fmt.Println(res)
	return res
}

func (mongo Mongo) GetAllTransCationData() (transCationInfo []s.Tx) {
	var res []s.Tx
	client.Find(bson.M{}).All(&res)
	return res
}

func (mongo Mongo) FindOutByAddress(address string) (result s.Address) {
	var res s.Address
	client.Find(bson.M{"address": address}).One(&res)
	return res
}
func (mongo Mongo) InsertOut(info s.Address) {
	client.Insert(info)
}

func (mongo Mongo) InsertErr(arr string, reason interface{}) {
	var errObj err
	errObj.ErrReason = reason
	errObj.ErrTransCationId = arr
	client.Insert(errObj)
}

func (mongo Mongo) GetAccountInfoByAddress(address string) (outAccount s.Address) {
	var res s.Address
	client.Find(bson.M{"address": address}).One(&res)
	return res
}

func (mongo Mongo) GetCountNumber() int {
	//var res AccountInfo
	Number, _ := client.Count()
	res := Number - 1
	return res
}

func (mongo Mongo) GetUnSpent(address string) []s.UTXO {
	var Unspent []s.UTXO
	client.Find(bson.M{"address": address}).All(&Unspent)
	return Unspent
}

func (mongo Mongo) GetRecentTransCation(pageNum int) (transCationInfo []s.Tx) {
	var res []s.Tx
	limitNum := 10
	skipNum := (pageNum) * 10
	client.Find(bson.M{}).Sort("-blocktime").Limit(limitNum).Skip(skipNum).All(&res)
	return res
}

func (mongo Mongo) GetRecentBlock(pageNum int) (Block []s.Blocks) {
	var res []s.Blocks
	limitNum := 10
	skipNum := (pageNum) * 10
	client.Find(bson.M{}).Sort("-height").Limit(limitNum).Skip(skipNum).All(&res)
	return res
}
