package btc

import (
	"log"

	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func BTCUnspent(txid string, Database *mgo.Database) ([]*s.UTXO, []*s.UTXO) {
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
		if r != nil {
			UTXOCollection.Remove(r)
		}
	}
	log.Print("Insert new UTXOs")
	for _, i := range Store {
		//i.Spent = ""
		//TODO
		//Build Uinque Index
		//UTXOCollection.Remove(i)
		//fmt.Println(i)
		if i != nil {
			err := UTXOCollection.Insert(i)
			if err != nil {
				log.Println("No need to store nil infor")
			}
		}
	}
	return Remove, Store
}

func VinUTXO(Vi *s.Vin) *s.UTXO {
	InUTXO := new(s.UTXO)
	//UTXO
	//	fmt.Println(Vi.Currency, "2222222222222")
	//Only Handle BTC UTXO
	if Vi.Address == "" || Vi.Currency != "BTC" {
		log.Println("Do not Hand USDT OR OMNI,Vin")
		return nil
	}
	InUTXO.Address = Vi.Address
	InUTXO.Index = Vi.Index
	InUTXO.Utxo = Vi.Hash
	InUTXO.Value = Vi.Value
	InUTXO.Currency = "BTC"
	//InUTXO.Spent = nil
	//InUTXO.Spent = true
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
	//OutUTXO.Spent = false
	OutUTXO.Utxo = txid
	OutUTXO.Value = Vo.Value
	OutUTXO.Currency = "BTC"
	return OutUTXO
}
