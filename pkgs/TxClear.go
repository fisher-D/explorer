package pkgs

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/GGBTC/explorer/service"
)

var (
	BTC  = "BTC"
	Omni = "Omni"
)

type Tx struct {
	TxID        string   `json:"TxID"`
	Index       uint32   `json:"-"`
	BlockHash   string   `json:"block_hash"`
	BlockHeight uint     `json:"block_height"`
	BlockTime   uint64   `json:"block_time"`
	TxIns       []*TxIn  `json:"Vin"`
	TxOuts      []*TxOut `json:"Vout"`
}

type TxOut struct {
	Addr     string `json:"address"`
	Value    uint64 `json:"value"`
	Index    uint32 `json:"index"`
	Currency string `json:"currency"`
}

type TxIn struct {
	Hash     string `json:"prevout,omitempty"`
	Address  string `json:"address,omitempty"`
	Value    uint64 `json:"value,omitempty"`
	Index    uint32 `json:"index,omitempty"`
	Coinbase string `json:"coinbase,omitempty"`
	Sequence int64  `json:"sequence,omitempty"`
	Currency string `json:"currency"`
}

func GetClearTx(txid string, height int) *Tx {
	//	res_tx, err := pkgs.GetTxRPC("1fcc63d2edbfe98b797234f5eed92b87ff2b155262d1ce07ee3db572cdc5d0c9", 595227)
	//fmt.Println(res)
	if txid == GenesisTx {
		return nil
		//return TxData{GenesisTx, []TxIn{}, []TxOut{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
	}
	//fmt.Println(txid)
	res_tx, err := CallBitcoinRPC(URL, "getrawtransaction", 1, []interface{}{txid, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	tx := new(Tx)
	//fmt.Println(res_tx)
	txjson := res_tx["result"].(map[string]interface{})
	//fmt.Println(txjson)
	blocktime, _ := txjson["blocktime"].(json.Number).Int64()
	tx.BlockTime = uint64(blocktime)
	tx.BlockHeight = uint(height)
	tx.BlockHash = txjson["blockhash"].(string)

	total_tx_out := uint64(0)
	total_tx_in := uint64(0)

	for _, txijson := range txjson["vin"].([]interface{}) {
		_, coinbase := txijson.(map[string]interface{})["coinbase"]
		if !coinbase {
			txi := new(TxIn)
			//txinjsonprevout := new(PrevOut)
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
				txi.Currency = "BTC"
			} else {
				prevout, _ := GetTxOutRPC(txi.Hash, txi.Index)
				txi.Address = prevout.Addr
				txi.Value = prevout.Value
				txi.Currency = "BTC"
			}

			total_tx_in += uint64(txi.Value)

			//txi.txi = txi

			tx.TxIns = append(tx.TxIns, txi)

			// TODO handle txi from this TX
		} else {
			txi := new(TxIn)
			txi.Coinbase = txijson.(map[string]interface{})["coinbase"].(string)
			txi.Sequence, _ = txijson.(map[string]interface{})["sequence"].(json.Number).Int64()
			tx.TxIns = append(tx.TxIns, txi)
			txi.Currency = "BTC"
		}
	}
	for _, txojson := range txjson["vout"].([]interface{}) {
		txo := new(TxOut)
		txoval, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
		txo.Value = uint64(txoval * 1e8)
		n, _ := txojson.(map[string]interface{})["n"].(json.Number).Float64()
		txo.Index = uint32(n)
		//txo.Addr = txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})[0].(string)
		if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nulldata" {
			txodata, txoisinterface := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})
			if txoisinterface {
				txo.Addr = txodata[0].(string)
				txo.Currency = "BTC"
			} else {
				//TODO Currecny
				txo.Addr = ""
			}
		} else {
			res := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["asm"].(string)
			Omni := "OP_RETURN 6f6d6e69"
			if strings.Contains(res, Omni) == true {
				txo.Addr = res
				txo.Currency = "USDT"
			}
		}
		tx.TxOuts = append(tx.TxOuts, txo)
		//	txospent := new(TxoSpent)
		//	txospent.Spent = false
		//	txo.Spent = txospent
		total_tx_out += uint64(txo.Value)
	}

	//tx.TxOutCnt = uint32(len(tx.TxOuts))
	//tx.TxInCnt = uint32(len(tx.TxIns))
	tx.TxID = txid
	//tx.TotalOut = uint64(total_tx_out)
	//tx.TotalIn = uint64(total_tx_in)
	return tx
}
func GetTxOutRPC(tx_id string, txo_vout uint32) (txo *TxOut, err error) {
	// Hard coded genesis tx since it's not included in bitcoind RPC API
	if tx_id == GenesisTx {
		return
		//return TxData{GenesisTx, []TxIn{}, []TxOut{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
	}
	// Get the TX from bitcoind RPC API
	res_tx, err := CallBitcoinRPC(URL, "getrawtransaction", 1, []interface{}{tx_id, 1})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	txjson := res_tx["result"].(map[string]interface{})

	txojson := txjson["vout"].([]interface{})[txo_vout]
	txo = new(TxOut)
	valtmp, _ := txojson.(map[string]interface{})["value"].(json.Number).Float64()
	txo.Value = service.FloatToUint(valtmp)
	if txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["type"].(string) != "nulldata" {
		txodata, txoisinterface := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["addresses"].([]interface{})
		if txoisinterface {
			txo.Addr = txodata[0].(string)
			txo.Currency = "BTC"
		} else {
			txo.Addr = ""
		}
	} else {
		res := txojson.(map[string]interface{})["scriptPubKey"].(map[string]interface{})["asm"].(string)
		Omni := "OP_RETURN 6f6d6e69"
		//Temporary
		//I use OP_RETURN 6f6d6e69 as condition to check wether a nulldata is an Omni or Not
		if strings.Contains(res, Omni) == true {
			txo.Addr = res
			txo.Currency = "USDT"
		}
	}
	//txospent := new(TxoSpent)
	//txospent.Spent = false
	//txo.Spent = txospent
	return
}

// type TxoSpent struct {
// 	Spent       bool   `json:"spent"`
// 	BlockHeight uint32 `json:"block_height,omitempty"`
// 	InputHash   string `json:"tx_hash,omitempty"`
// 	InputIndex  uint32 `json:"in_index,omitempty"`
// }
// type TxSample struct {
// 	TxID        string `json:"TxID"`
// 	BlockHash   string `json:"block_hash"`
// 	BlockHeight int    `json:"block_height"`
// 	BlockTime   int    `json:"block_time"`
// 	Vin         []struct {
// 		Txid    string `json:"txid"`
// 		Address string `json:"address"`
// 		Value   int    `json:"value"`
// 		N       int    `json:"n"`
// 	} `json:"Vin"`
// 	Vout []struct {
// 		Address string `json:"address"`
// 		Value   int    `json:"value"`
// 		N       int    `json:"n"`
// 	} `json:"Vout"`
// }
