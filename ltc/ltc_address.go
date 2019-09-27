package ltc

import (
	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetAddress(time uint64, in []*s.UTXO, out []*s.UTXO, Database *mgo.Database) {
	addressCollection := Database.C("address")
	for _, k := range in {
		predata := VinInfo(k)
		FinishAddress(time, predata, addressCollection)
	}
	for _, k := range out {
		predata1 := VoutInfo(k)
		FinishAddress(time, predata1, addressCollection)
	}
}
func VinInfo(InUTXO *s.UTXO) *s.Address {
	if InUTXO == nil {
		return nil
	}
	if InUTXO.Address == "" {
		return nil
	}
	var Txi s.Txs
	var Txis []s.Txs
	Addre := new(s.Address)
	Addre.Address = InUTXO.Address
	Txi.Index = InUTXO.Index
	Txi.Txid = InUTXO.Utxo
	Txi.Value = InUTXO.Value
	Txi.Currency = "LTC"
	InUTXO.Spent = true
	Txis = append(Txis, Txi)
	Addre.Txs = Txis
	return Addre
}

func VoutInfo(OutUTXO *s.UTXO) *s.Address {
	var Txi s.Txs
	var Txis []s.Txs
	Addre := new(s.Address)
	Txi.Index = OutUTXO.Index
	Txi.Txid = OutUTXO.Utxo
	Txi.Value = OutUTXO.Value
	Txi.Currency = "LTC"
	OutUTXO.Spent = false
	Addre.Address = OutUTXO.Address
	Txis = append(Txis, Txi)
	Addre.Txs = Txis
	return Addre
}

func FillParas(addreinfo *s.Address) *s.Address {
	for _, k := range addreinfo.Txs {
		if k.Spent != true {
			addreinfo.TotalRecCount++
			addreinfo.TotalReceived = +k.Value
		} else {
			addreinfo.TotalSentCount++
			addreinfo.TotalSent = +k.Value
		}
	}
	addreinfo.Balance = addreinfo.TotalReceived - addreinfo.TotalSent

	return addreinfo
}

//Not Found
func CompleteAddress(Time uint64, addreinfo *s.Address) *s.Address {
	addreinfo.FirstSeen = Time
	addreinfo.LastSeen = Time
	res := FillParas(addreinfo)
	if res.Balance < 0 {
		res.Balance = 0
	}
	return res
}
func UpdateAddress(Time uint64, olds s.Address, news *s.Address) *s.Address {
	news.FirstSeen = olds.FirstSeen
	news.LastSeen = Time
	for _, k := range olds.Txs {
		news.Txs = append(news.Txs, k)
	}
	res := FillParas(news)
	res.TotalReceived = res.TotalReceived + olds.TotalReceived
	res.TotalSent = res.TotalSent + olds.TotalSent
	res.Balance = res.Balance + olds.Balance
	if res.Balance < 0 {
		res.Balance = 0
	}
	return res

}
func FinishAddress(Time uint64, addressinfo *s.Address, collection *mgo.Collection) bool {
	if addressinfo != nil {
		var olds s.Address
		query := bson.M{"address": addressinfo.Address}
		err := collection.Find(query).One(&olds)
		if err != nil {
			firstdata := CompleteAddress(Time, addressinfo)
			collection.Insert(firstdata)
		} else {
			updateresult := UpdateAddress(Time, olds, addressinfo)
			collection.Remove(olds)
			collection.Insert(updateresult)
		}
		return true
	}
	return false
}
