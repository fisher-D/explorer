package main

import (
	"fmt"

	"github.com/GGBTC/explorer/pkgs"
	"github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2/bson"
)

const (
	URL       = "EXAMPLE"
	GenesisTx = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"
	mongourl  = "localhost:27017"
)

type BTCconfig struct {
	URL       string
	GenesisTx string
	mongourl  string
}
type Test1 struct {
	Name  string
	Input Input
}
type Input struct {
	Age    int
	Gender string
}

func main() {
	// service.GetMongo(mongourl)
	// Txid := "0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9"
	// GetUnSpentTransaction(Txid)
	c := new(Input)
	c.Age = 123
	c.Gender = "boy"
	d := new(Input)
	d.Age = 1123
	d.Gender = "boy1"
	selector := bson.M{"Name": "jack"}
	pretx := bson.M{"Input": d}
	data := bson.M{"$addToSet": pretx}
	service.GetMongo(mongourl)
	session := service.GlobalS.DB("test").C("asd")
	a := new(Test1)
	a.Name = "jack"

	a.Input = *c
	session.Upsert(selector, data)
}

func GetUnSpentTransaction(Txid string) {
	session := service.GlobalS.DB("GGBTC").C("transaction")
	var q pkgs.Tx
	session.Find(bson.M{"txid": Txid}).One(&q)
	var Vin []*pkgs.TxIn
	var Vout []*pkgs.TxOut
	Vin = q.TxIns
	Vout = q.TxOuts
	UTXOI := drawVinOut(Vin)
	UTXOT := drawVoutOut(Vout, Txid)
	//fmt.Println(UTXOI)
	InsertVout(UTXOT)
	//fmt.Println("========================")
	//fmt.Println(UTXOT)
	RemoveVin(UTXOI)

}

// func mongoUpdateTest() {
// 	session := service.GlobalS.DB("GGBTC").C("Test")
// 	var q Unspent
// 	target := bson.M{"height": "asd"}
// 	session.Update()
// }
func drawVinOut(Vin []*pkgs.TxIn) *Unspent {
	UTXO := new(Unspent)
	if Vin[0].Coinbase != "" {
		return nil
	}
	for _, k := range Vin {
		UTXO.address = k.Address
		UTXO.txid = k.Hash
		UTXO.value = k.Value
		UTXO.currency = k.Currency
		UTXO.index = k.Index
		return UTXO
	}
	return nil
}
func drawVoutOut(Vout []*pkgs.TxOut, Txid string) *Unspent {
	UTXO := new(Unspent)
	for _, k := range Vout {
		UTXO.address = k.Addr
		UTXO.txid = Txid
		UTXO.value = k.Value
		UTXO.currency = k.Currency
		UTXO.index = k.Index
		return UTXO
	}
	return nil
}
func InsertVout(UTXOT *Unspent) {
	if UTXOT != nil {

		session := service.GlobalS.DB("GGBTC").C("unspent")
		session.Insert(UTXOT)
		fmt.Println("Vout Update Success")
	} else {
		fmt.Println("Coinbase Tx does not need to be inster")
	}
}
func RemoveVin(UTXOI *Unspent) {
	session := service.GlobalS.DB("GGBTC").C("unspent")
	session.Remove(UTXOI)
	fmt.Println("Vin Remove Success")
}

//TODO
//Test the VinRemove And VoutInsert
type Unspent struct {
	txid     string
	address  string
	value    uint64
	currency string
	index    uint32
}

func ClearTx() {
	// GetUnSpentTransaction(Txid)
	service.GetMongo(mongourl)
	service.GlobalS.DB("GGBTC").C("unspent").DropCollection()
	service.GlobalS.DB("GGBTC").C("transaction").DropCollection()
	service.GlobalS.DB("GGBTC").C("status").DropCollection()
	service.GlobalS.DB("GGBTC").C("address").DropCollection()
	service.GlobalS.DB("GGBTC").C("addresrelatetx").DropCollection()
}
