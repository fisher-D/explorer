package pkgs

import (
	"fmt"

	"github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2/bson"
)

//TODO
//Test the VinRemove And VoutInsert
type Unspent struct {
	Txid     string
	Address  string
	Value    uint64
	Currency string
	Index    uint32
}

func GetUnSpentTransaction(Txid string) {
	session := service.GlobalS.DB("GGBTC").C("transaction")
	var q Tx
	session.Find(bson.M{"txid": Txid}).One(&q)
	//fmt.Println(q)
	var Vin []*TxIn
	var Vout []*TxOut
	Vin = q.TxIns
	Vout = q.TxOuts
	UTXOI := drawVinOut(Vin)
	UTXOT := drawVoutOut(Vout, Txid)
	//fmt.Println(UTXOI)
	RemoveVin(UTXOI)
	InsertVout(UTXOT)
	//fmt.Println("========================")
	//fmt.Println(UTXOT)

}

// func mongoUpdateTest() {
// 	session := service.GlobalS.DB("GGBTC").C("Test")
// 	var q Unspent
// 	target := bson.M{"height": "asd"}
// 	session.Update()
// }
func drawVinOut(Vin []*TxIn) *Unspent {
	UTXO := new(Unspent)
	if Vin == nil {
		return nil
	}

	if Vin[0].Coinbase != "" {
		return nil
	}
	for _, k := range Vin {
		UTXO.Address = k.Address
		UTXO.Txid = k.Hash
		UTXO.Value = k.Value
		UTXO.Currency = k.Currency
		UTXO.Index = k.Index
		return UTXO
	}
	return nil
}
func drawVoutOut(Vout []*TxOut, Txid string) *Unspent {
	UTXO := new(Unspent)
	for _, k := range Vout {
		UTXO.Address = k.Addr
		UTXO.Txid = Txid
		UTXO.Value = k.Value
		UTXO.Currency = k.Currency
		UTXO.Index = k.Index
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
		fmt.Println("Coinbase Tx does not need to be instert")
	}
}
func RemoveVin(UTXOI *Unspent) {
	if UTXOI != nil {
		session := service.GlobalS.DB("GGBTC").C("unspent")
		session.Remove(UTXOI)
		fmt.Println("Vin Remove Success")
	} else {
		fmt.Println("Coinbase Tx does not need to be instert")
	}
}
