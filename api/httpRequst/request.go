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
	GetTransCationInfoAddress(value string) (transCationInfo []s.Tx)
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
	client.Find(bson.M{"hash": value}).One(&res)
	fmt.Println(res)
	return res
}

func (mongo Mongo) GetTransCationInfoAddress(value string) (transCationInfo []s.Tx) {
	var res []s.Tx
	client.Find(bson.M{"$or": []bson.M{bson.M{"in.prev_out.address": value}, bson.M{"out.prev_out.address": value}}}).All(&res)
	return res
}
