package btcd

import (
	s "github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//TODO
//进行并发设计
func GetAddress(Tx *s.Tx, Database *mgo.Database) string {
	if Tx != nil {
		time := Tx.BlockTime
		in := Tx.Vin
		out := Tx.Vout
		txid := Tx.Txid
		addressIndex := mgo.Index{
			Key:    []string{"address"},
			Unique: true,
		}
		addressCollection := Database.C("address")
		addressCollection.EnsureIndex(addressIndex)
		//addressCollection.EnsureIndex(addressIndex2)
		for _, k := range in {
			predata := VinInfo(k)
			FinishAddress(time, predata, addressCollection)
		}
		for _, k := range out {
			//fmt.Println(k)
			predata1 := VoutInfo(k, txid)
			FinishAddress(time, predata1, addressCollection)
		}
	}
	return "Success"
}

func VinInfo(InUTXO *s.Vin) *s.Address {
	if InUTXO.Address == "" {
		return nil
	}
	var Txi s.Txs
	var Txis []s.Txs
	Addre := new(s.Address)
	Addre.Address = InUTXO.Address
	Txi.Index = InUTXO.Index
	Txi.Txid = InUTXO.Hash
	Txi.Value = InUTXO.Value
	Txi.Currency = "BTCD"
	InUTXO.Spent = "true"
	Txi.Spent = "Ture"
	Txis = append(Txis, Txi)
	Addre.Txs = Txis
	return Addre
}

func VoutInfo(OutUTXO *s.VoutNew, txid string) *s.Address {
	if OutUTXO.Addr == "" {
		return nil
	}
	var Txi s.Txs
	var Txis []s.Txs
	Addre := new(s.Address)
	Txi.Index = OutUTXO.Index
	Txi.Txid = txid
	Txi.Value = OutUTXO.Value
	Txi.Currency = "BTCD"
	OutUTXO.Spent = "false"
	Txi.Spent = "False"
	Addre.Address = OutUTXO.Addr
	Txis = append(Txis, Txi)
	Addre.Txs = Txis
	return Addre
}

func FillParas(addreinfo *s.Address) *s.Address {
	for _, k := range addreinfo.Txs {
		if k.Spent != "true" {
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
	//news Txs is not empty
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

func FinishAddress(Time uint64, addressinfo *s.Address, collection *mgo.Collection) {
	if addressinfo != nil {
		var olds s.Address
		query := bson.M{"address": addressinfo.Address}
		err := collection.Find(query).One(&olds)
		if err != nil {
			firstdata := CompleteAddress(Time, addressinfo)
			collection.Insert(firstdata)
		} else {
			updateresult := UpdateAddress(Time, olds, addressinfo)
			collection.Remove(bson.M{"address": olds.Address})
			collection.Insert(updateresult)
		}
		//return true
	}
	//return false
}
