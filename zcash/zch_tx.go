package zec

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/GGBTC/explorer/service"
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
			result, _ := GetClearTx(k)
			err := TxCollection.Insert(result)
			if err != nil {
				log.Println("What could it be ?")
				return false
			}
			Time = result.BlockTime
			Vin, Vout := ZECUnspent(k, Database)
			GetAddress(Time, Vin, Vout, Database)
		}
	}

	return true
}

func GetClearTx(txid string) (tx *service.Tx, err error) {
	if txid == GenesisTx {
		return
	}
	res_tx, err := CallZECRPC(URL, "getrawtransaction", 1, []interface{}{txid, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	tx = new(service.Tx)
	txjson := res_tx["result"].(map[string]interface{})
	blocktime, _ := txjson["blocktime"].(json.Number).Int64()
	tx.BlockTime = uint64(blocktime)
	tx.BlockHash = txjson["blockhash"].(string)

	total_tx_out := uint64(0)
	total_tx_in := uint64(0)

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
				txi.Currency = "ZEC"
				txi.Spent = true
			} else {
				prevout, _ := GetVoutNewRPC(txi.Hash, txi.Index)
				txi.Address = prevout.Addr
				txi.Value = prevout.Value
				txi.Currency = "ZEC"
				txi.Spent = true
			}

			total_tx_in += uint64(txi.Value)
			tx.Vin = append(tx.Vin, txi)
		} else {
			txi := new(service.Vin)
			txi.Coinbase = txijson.(map[string]interface{})["coinbase"].(string)
			txi.Sequence, _ = txijson.(map[string]interface{})["sequence"].(json.Number).Int64()
			tx.Vin = append(tx.Vin, txi)
			txi.Currency = "ZEC"
			txi.Spent = true
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
				txo.Currency = "ZEC"
				//txo.Currency = "ZEC"
				txo.Spent = false
			} else {
				txo.Addr = ""
			}
			// } else {
			// 	res := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["asm"].(string)
			// 	//Omni, err := OmniProcesser(res)
			// 	if err != nil {
			// 		txo.Addr = "Unknown"
			// 		txo.Currency = "Not strandard Omni"
			// 	} else {
			// 		txo.Addr = Omni.OP_RETURN
			// 		txo.Currency = Omni.TokenName
			// 		//txo.Index =
			// 		txo.Value = Omni.Value
			// 	}
		}

		tx.Vout = append(tx.Vout, txo)
		total_tx_out += uint64(txo.Value)
	}
	tx.Txid = txid
	return tx, nil
}
func GetVoutNewRPC(tx_id string, txo_vout uint32) (txo *service.VoutNew, err error) {
	// Hard coded genesis tx since it's not included in bitcoind RPC API
	if tx_id == GenesisTx {
		return
		//return TxData{GenesisTx, []Vin{}, []VoutNew{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
	}
	res_tx, err := CallZECRPC(URL, "getrawtransaction", 1, []interface{}{tx_id, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	txjson := res_tx["result"].(map[string]interface{})

	txojson := txjson["vout"].([]interface{})[txo_vout]
	txo = new(service.VoutNew)
	valtmp, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
	txo.Value = service.FloatToUint(valtmp)
	if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nulldata" {
		txodata, txoisinterface := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})
		if txoisinterface {
			txo.Addr = txodata[0].(string)

		} else {
			txo.Addr = ""
		}
		// } else {
		// 	res := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["asm"].(string)
		// 	Omni, err := OmniProcesser(res)
		// 	if err != nil {
		// 		txo.Addr = "Unknown"
		// 		txo.Currency = "Not strandard Omni"
		// 	} else {
		// 		txo.Addr = Omni.OP_RETURN
		// 		txo.Currency = Omni.TokenName
		// 		//txo.Index =
		// 		txo.Value = Omni.Value
		// 	}
	}

	return
}
func GetTxRPC(txid string) map[string]interface{} {
	//var Tx Tx
	if txid == GenesisTx {
		return nil
	}
	res_tx, err := CallZECRPC(URL, "getrawtransaction", 1, []interface{}{txid, 1})
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
