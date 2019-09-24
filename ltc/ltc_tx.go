package ltc

import (
	"encoding/json"
	"fmt"
	"log"

	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CatchUpTx(txidArray []string, Database *mgo.Database) bool {
	TxCollection := Database.C("txs")
	var Time uint64
	for _, k := range txidArray {
		var q bson.M
		res := TxCollection.Find(bson.M{"txid": k}).One(&q)
		if res != nil {
			result := GetClearTx(k)
			err := TxCollection.Insert(result)
			if err != nil {
				log.Println("What could it be ?")
				return false
			}
			Time = result.BlockTime
			Vin, Vout := LTCUnspent(k, Database)
			GetAddress(Time, Vin, Vout, Database)
		}
	}

	return true
}

func GetClearTx(txid string) s.Tx {
	var Tx s.Tx
	if txid == GenesisTx {
		return Tx
	}
	ress := GetTxRPC(txid)
	var Txinfo s.TxOld
	data, _ := json.Marshal(ress)
	json.Unmarshal(data, &Txinfo)
	json.Unmarshal(data, &Tx)
	var VV []*s.VoutNew
	var WW []*s.Vin
	for v, k := range Tx.Vout {
		tar := Txinfo.Vout[v].ScriptPubKey
		k.Value = Txinfo.Vout[v].Value
		//fmt.Println("========================")
		//	fmt.Println(k.Value)
		k.Addr = tar.Addresses[0]
		//k.Type = tar.Type
		k.Currency = "LTC"
		VV = append(VV, k)
	}
	Tx.Vout = VV
	for _, in := range Tx.Vin {
		in.Currency = "LTC"
		inTxid := in.Hash
		inIndex := in.Index
		in.Value, in.Address = GetVinValue(inTxid, inIndex)
		WW = append(WW, in)
	}
	Tx.Vin = WW
	return Tx
}
func GetVinValue(txid string, index uint32) (uint64, string) {
	ress := GetTxRPC(txid)
	var Txinfo s.TxOld
	data, _ := json.Marshal(ress)
	json.Unmarshal(data, &Txinfo)
	var Value uint64
	var Address string
	for _, k := range Txinfo.Vout {
		if k.N == uint64(index) {
			Value = k.Value
			Address = k.ScriptPubKey.Addresses[0]
			return Value, Address
		}
	}
	return 0, ""
}
func GetTxRPC(txid string) map[string]interface{} {
	//var Tx Tx
	if txid == GenesisTx {
		return nil
	}
	res_tx, err := CallLTCRPC(URL, "getrawtransaction", 1, []interface{}{txid, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	pre := res_tx["result"]
	if pre == nil {
		return nil
	}
	ress := res_tx["result"].(map[string]interface{})
	if ress == nil {
		fmt.Println("Intersting")
		return nil
	}
	return ress
}
