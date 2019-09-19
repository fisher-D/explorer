package service

import (
	"log"

	"gopkg.in/mgo.v2"
)

var GlobalS *mgo.Session

//GetMongo 。。。。
func GetMongo(URL string) {
	s, err := mgo.Dial(URL)
	if err != nil {
		log.Fatalf("Create Session: %s\n", err)
	}
	GlobalS = s
}

//Connect ...
func Connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	ms := GlobalS.Copy()
	c := ms.DB(db).C(collection)
	ms.SetMode(mgo.Monotonic, true)
	return ms, c
}

// Insert ....
func Insert(db, collection string, doc interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()

	return c.Insert(doc)
}

// FindOne ...
// selector ： bson.M{"_data":"parameters"}
func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).One(result)
}

//FindAll 。。。
func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).All(result)
}

// Update ...
func Update(db, collection string, selector, update interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()

	return c.Update(selector, update)
}

//Upsert 更新，如果不存在就插入一个新的数据 `upsert:true`
func Upsert(db, collection string, selector, update interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

// Count 查询表内数量
func Count(db, collection string, query interface{}) (int, error) {
	ms, c := Connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}
