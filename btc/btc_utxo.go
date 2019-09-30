package btc

import (
	"log"

	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//TODO
//将之拆分为两个并行的方法
func BTCUnspent(txid string, Database *mgo.Database) ([]*s.UTXO, []*s.UTXO) {
	TxCollection := Database.C("txs")
	var Txtar s.Tx
	query := bson.M{"txid": txid}
	TxCollection.Find(query).One(&Txtar)
	//Add uinque Index for Collection UTXO
	utxoIndex := mgo.Index{
		Key:    []string{"utxo"},
		Unique: true,
	}
	UTXOCollection := Database.C("utxos")
	UTXOCollection.EnsureIndex(utxoIndex)
	Remove := RmoveProcss(Txtar, UTXOCollection)
	Store := InsertProcess(txid, Txtar, UTXOCollection)
	return Remove, Store
}

func RmoveProcss(Txtar s.Tx, UTXOCollection *mgo.Collection) []*s.UTXO {
	var Remove []*s.UTXO
	for _, k := range Txtar.Vin {
		listi := VinUTXO(k)
		Remove = append(Remove, listi)
	}
	for _, r := range Remove {
		if r != nil {
			UTXOCollection.Remove(r)
		}
	}
	log.Print("Remove used UTXOs")
	return Remove
}
func InsertProcess(txid string, Txtar s.Tx, UTXOCollection *mgo.Collection) []*s.UTXO {
	var Store []*s.UTXO
	for _, k := range Txtar.Vout {
		listo := VoutUTXO(k, txid)
		Store = append(Store, listo)
	}
	for _, i := range Store {
		if i != nil {
			err := UTXOCollection.Insert(i)
			if err != nil {
				log.Println("No need to store nil infor")
			}
		}
	}
	log.Print("Insert new UTXOs")
	return Store
}
func VinUTXO(Vi *s.Vin) *s.UTXO {
	InUTXO := new(s.UTXO)
	if Vi.Address == "" || Vi.Currency != "BTC" {
		log.Println("Do not Hand USDT OR OMNI,Vin")
		return nil
	}
	InUTXO.Address = Vi.Address
	InUTXO.Index = Vi.Index
	InUTXO.Utxo = Vi.Hash
	InUTXO.Value = Vi.Value
	InUTXO.Currency = "BTC"
	return InUTXO

}

func VoutUTXO(Vo *s.VoutNew, txid string) *s.UTXO {
	if Vo.Currency != "BTC" {
		log.Println("Do not Hand USDT OR OMNI,Vout")
		return nil
	}
	OutUTXO := new(s.UTXO)
	OutUTXO.Address = Vo.Addr
	OutUTXO.Index = Vo.Index
	OutUTXO.Utxo = txid
	OutUTXO.Value = Vo.Value
	OutUTXO.Currency = "BTC"
	return OutUTXO
}
