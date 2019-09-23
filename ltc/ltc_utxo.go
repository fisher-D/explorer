package ltc

import (
	"fmt"
	"log"

	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func LTCUnspent(txid string, Database *mgo.Database) ([]*s.UTXO, []*s.UTXO) {
	TxCollection := Database.C("txs")
	var Txtar s.Tx
	query := bson.M{"txid": txid}
	TxCollection.Find(query).One(&Txtar)
	var Remove []*s.UTXO
	var Store []*s.UTXO
	for _, k := range Txtar.Vin {
		listi := VinUTXO(k)
		Remove = append(Remove, listi)
	}
	//fmt.Println(Remove)
	for _, k := range Txtar.Vout {
		listo := VoutUTXO(k, txid)
		Store = append(Store, listo)
	}
	//	fmt.Println(Store)
	UTXOCollection := Database.C("utxo")
	log.Print("Remove used UTXOs")
	for _, r := range Remove {
		UTXOCollection.Remove(r)
	}
	log.Print("Insert new UTXOs")
	for _, i := range Store {
		UTXOCollection.Remove(i)
		err := UTXOCollection.Insert(i)
		if err != nil {
			fmt.Println("Fuck Store")
		}
	}
	return Remove, Store
}

func VinUTXO(Vi s.Vin) *s.UTXO {
	InUTXO := new(s.UTXO)
	if Vi.Address != "" {
		return nil
	}
	InUTXO.Address = Vi.Address
	InUTXO.Index = Vi.Vout
	InUTXO.TxID = Vi.Txid
	InUTXO.Value = Vi.Value
	InUTXO.Currency = "LTC"
	//InUTXO.Spent = true
	return InUTXO

}

func VoutUTXO(Vo s.VoutNew, txid string) *s.UTXO {
	OutUTXO := new(s.UTXO)
	OutUTXO.Address = Vo.Addresses
	OutUTXO.Index = Vo.N
	//OutUTXO.Spent = false
	OutUTXO.TxID = txid
	OutUTXO.Value = Vo.Value
	OutUTXO.Currency = "LTC"
	return OutUTXO
}
