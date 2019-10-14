package ltc

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/GGBTC/explorer/service"
	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CatchUpTx() string {
	s.GetMongo(mongourl)
	Database := s.GlobalS.DB("LTC")
	star, end := GenerateTime(Database)
	if star != end {
		log.Println("Start")
		for i := star; i <= end; i++ {
			list := GenerateTxArray(i, Database)
			ProcessTx(list, Database)
		}
	} else {
		return "Success"
	}
	return "Success"
}

//10:55:37
func GenerateTime(Database *mgo.Database) (uint64, uint64) {
	blockCollection := Database.C("blocks")
	txIndex1 := mgo.Index{
		Key:    []string{"txid"},
		Unique: false,
	}
	txIndex2 := mgo.Index{
		Key:    []string{"-blocktime"},
		Unique: false,
	}

	TxCollection := Database.C("txs")
	TxCollection.EnsureIndex(txIndex1)
	TxCollection.EnsureIndex(txIndex2)
	var target s.Blocks
	var targettx s.Tx
	var starttime uint64
	blockCollection.Find(bson.M{}).Sort("-height").Limit(1).One(&target)
	endtime := uint64(target.Height)
	TxCollection.Find(bson.M{}).Sort("-blocktime").Limit(1).One(&targettx)
	time := targettx.BlockTime
	if time != 0 {
		targetbl := new(s.Blocks)
		blockCollection.Find(bson.M{"time": time}).One(&targetbl)
		starttime = uint64(targetbl.Height)
	} else {
		starttime = 0
	}

	return starttime, endtime

}
func GenerateTxArray(time uint64, Database *mgo.Database) []string {
	blockCollection := Database.C("blocks")
	var target s.Blocks
	blockCollection.Find(bson.M{"height": time}).One(&target)
	list := target.Tx
	return list
}
func BuildUTXO(tx *s.Tx, Database *mgo.Database) string {

	if tx == nil {
		return "Success"
	}
	UTXOIndex := mgo.Index{
		Key:    []string{"utxo"},
		Unique: false,
	}
	UtxoCollection := Database.C("utxos")
	UtxoCollection.EnsureIndex(UTXOIndex)
	vin := tx.Vin
	vout := tx.Vout
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, k := range vin {
			if k.Address != "" {
				Remove := new(s.UTXO)
				Remove.Address = k.Address
				Remove.Currency = k.Currency
				Remove.Index = k.Index
				Remove.Utxo = k.Hash
				Remove.Value = k.Value
				UtxoCollection.Remove(Remove)
			}
		}
	}()
	wg.Add(1)

	go func() {
		defer wg.Done()
		for _, k := range vout {
			if k.Addr != "" {
				Insert := new(s.UTXO)
				Insert.Address = k.Addr
				Insert.Currency = k.Currency
				Insert.Index = k.Index
				Insert.Utxo = tx.Txid
				Insert.Value = k.Value
				UtxoCollection.Insert(Insert)
			}
		}
	}()
	wg.Wait()
	log.Println("UXTO Process Finish")
	return "Success"

}
func ProcessTx(txidArray []string, Database *mgo.Database) bool {
	var wg sync.WaitGroup
	TxCollection := Database.C("txs")
	for _, k := range txidArray {
		var q bson.M
		res := TxCollection.Find(bson.M{"txid": k}).One(&q)
		if res != nil {
			result, _ := GetClearTx(k, TxCollection)
			TxCollection.Insert(result)
			wg.Add(1)
			go func() {
				defer wg.Done()
				resul := BuildUTXO(result, Database)
				if resul == "Success" {
					log.Println("Inserting Tx :", k)
				}
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				resull := GetAddress(result, Database)
				if resull == "Success" {
					log.Println("Address Build Finish")
				}
			}()
			wg.Wait()
		}
	}
	return true
}

func GetClearTx(txid string, TxCollection *mgo.Collection) (tx *service.Tx, err error) {
	if txid == GenesisTx {
		return
	}
	res_tx, err := CallLTCRPC(URL, "getrawtransaction", 1, []interface{}{txid, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	tx = new(service.Tx)
	txjson := res_tx["result"].(map[string]interface{})
	//fmt.Println(txjson)
	blocktime, _ := txjson["blocktime"].(json.Number).Int64()
	tx.BlockTime = uint64(blocktime)
	tx.BlockHash = txjson["blockhash"].(string)
	Version, _ := txjson["version"].(json.Number).Int64()
	tx.Version = uint32(Version)
	total_tx_out := uint64(0)
	total_tx_in := uint64(0)
	//var baseinfor string
	for _, txijson := range txjson["vin"].([]interface{}) {
		_, coinbase := txijson.(map[string]interface{})["coinbase"]
		if !coinbase {
			txi := new(service.Vin)
			txi.Hash = txijson.(map[string]interface{})["txid"].(string)
			tmpvout, _ := txijson.(map[string]interface{})["vout"].(json.Number).Int64()
			txi.Index = uint32(tmpvout)

			// Check if bitcoind is patched to fetch value/address without additional RPC call
			// cf. README
			_, bitcoindPatched := txijson.(map[string]interface{})["value"]
			if bitcoindPatched {
				pval, _ := txijson.(map[string]interface{})["value"].(json.Number).Float64()
				txi.Address = txijson.(map[string]interface{})["address"].(string)
				txi.Value = service.FloatToUint(pval)
				txi.Currency = "LTC"
				txi.Spent = "true"
			} else {
				prevout, _ := GetVoutNewRPC(txi.Hash, txi.Index, TxCollection)
				txi.Address = prevout.Addr
				txi.Value = prevout.Value
				txi.Currency = "LTC"
				txi.Spent = "true"

			}
			total_tx_in += uint64(txi.Value)
			tx.Vin = append(tx.Vin, txi)
		} else {
			txi := new(service.Vin)
			txi.Coinbase = txijson.(map[string]interface{})["coinbase"].(string)
			//baseinfor = txi.Coinbase
			txi.Sequence, _ = txijson.(map[string]interface{})["sequence"].(json.Number).Int64()
			tx.Vin = append(tx.Vin, txi)
			txi.Currency = "LTC"
			txi.Spent = "true"
			total_tx_in += uint64(txi.Value)
		}

	}
	for _, txojson := range txjson["vout"].([]interface{}) {
		txo := new(service.VoutNew)
		txoval, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
		txo.Value = uint64(txoval * 1e8)
		n, _ := txojson.(map[string]interface{})["n"].(json.Number).Float64()
		txo.Index = uint32(n)
		if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nulldata" {
			txodata, txoisinterface := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})
			if txoisinterface {
				txo.Addr = txodata[0].(string)
				txo.Currency = "LTC"
				tx.Type = "LTC"
				txo.Spent = "false"

			} else {
				txo.Addr = ""
			}
		}
		tx.Vout = append(tx.Vout, txo)
		//if txo.Currency == "LTC" {
		total_tx_out += uint64(txo.Value)
		//}
	}

	tx.Totalin = uint64(total_tx_in)
	tx.Totalout = uint64(total_tx_out)
	tx.Txid = txid
	return tx, nil
}

func GetVoutNewRPC(tx_id string, txo_vout uint32, TxCollection *mgo.Collection) (txo *service.VoutNew, err error) {
	// Hard coded genesis tx since it's not included in bitcoind RPC API
	if tx_id == GenesisTx {
		return
		//return TxData{GenesisTx, []Vin{}, []VoutNew{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
	}
	var res_tx1 *s.Tx
	TxCollection.Find(bson.M{"txid": tx_id}).One(&res_tx1)
	txo = new(s.VoutNew)
	for v, k := range res_tx1.Vout {
		if uint32(v) == txo_vout {
			txo = k
			return txo, nil
		}
		return
	}

	return
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
