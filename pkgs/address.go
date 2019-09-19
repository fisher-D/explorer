package pkgs

import (
	"fmt"

	"github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2/bson"
)

type CurrectAddress struct {
	Address   string `json:"address"`
	Txdetails []struct {
		Index    int    `json:"index"`
		Spent    bool   `json:"spent"`
		Txid     string `json:"txid"`
		Value    int    `json:"value"`
		currecny string `json:"coinName"`
	} `json:"txdetails"`
}
type Address struct {
	Result    CurrectAddress
	Detail    AddrDetail
	FirstSeen uint64
	LastSeen  uint64
}
type AddrDetail struct {
	TotalRecCount  int
	TotalSentCount int
	TotalReceived  uint64
	TotalSent      uint64
	Balance        uint64
}
type Txins struct {
	TxID    string
	VoutNum int
}
type Voutinfo struct {
	Address string
	Txid    string
	Index   int
	Value   uint64
	Spent   bool
}

func GetAddressRPC(TxID string) {
	VinMesses, VoutMess, Time := QueryTest(TxID)
	Vout := ProcessVoutMess(VoutMess)
	Vin := ProcessVinMess(VinMesses)
	VoutAddress := SaveAddressData(Vout)
	VinAddress := SaveAddressData(Vin)
	for _, k := range VinAddress {
		VoutAddress = append(VoutAddress, k)
	}
	CompleteAddress(VoutAddress, Time)
}

func SaveAddressData(Address []service.Address) []string {
	service.GetMongo(mongourl)
	session := service.GlobalS.DB("GGBTC").C("addresrelatetx")
	var Addres []string
	for _, k := range Address {
		Addr := k.Address
		Addres = append(Addres, Addr)
		selector := bson.M{"address": k.Address}
		pretx := bson.M{"txdetails": k.TxDetails[0]}
		data := bson.M{"$addToSet": pretx}
		_, err := session.Upsert(selector, data)
		if err != nil {
			fmt.Println(err)
		}
	}
	return Addres
}

func CompleteAddress(address []string, Time uint64) {
	service.GetMongo(mongourl)
	session := service.GlobalS.DB("GGBTC").C("addresrelatetx")
	AddressSession := service.GlobalS.DB("GGBTC").C("address")
	for _, k := range address {
		query := bson.M{"address": k}
		var q []CurrectAddress
		var Final Address
		session.Find(query).All(&q)
		Final.Result = q[0]
		if len(Final.Result.Txdetails) > 1 {
			Final.LastSeen = Time
		} else {
			Final.FirstSeen = Time
			Final.LastSeen = Time
		}
		var totalreceive, totalsent int
		balance, sent, receive := 0, 0, 0
		for _, k := range Final.Result.Txdetails {
			if k.Spent != true {
				totalreceive += k.Value
				receive++
			} else {
				totalsent += k.Value
				sent++
			}
		}
		prebalance := totalreceive - totalsent
		if prebalance < 0 {
			balance = 0
		} else {
			balance = prebalance
		}
		Final.Detail.Balance = uint64(balance)
		Final.Detail.TotalRecCount = receive
		Final.Detail.TotalReceived = uint64(totalreceive)
		Final.Detail.TotalSent = uint64(totalsent)
		Final.Detail.TotalSentCount = sent
		AddressSession.Insert(Final)
	}
}

func ProcessVinMess(VinMesses [][]Voutinfo) []service.Address {
	var Res service.Address
	var Ress []service.Address
	for _, k := range VinMesses {
		for _, m := range k {
			Res = ProcessData(m)
			Ress = append(Ress, Res)
		}
	}
	return Ress
}

func ProcessVoutMess(VoutMess []Voutinfo) []service.Address {
	var Res service.Address
	var Ress []service.Address
	for _, k := range VoutMess {
		Res = ProcessData(k)
		Ress = append(Ress, Res)
	}
	return Ress
}

func ProcessData(VinMess Voutinfo) service.Address {
	var NewAddress service.Address
	var NewTxDetails service.TxDetail
	NewAddress.Address = VinMess.Address
	NewTxDetails.Index = VinMess.Index
	NewTxDetails.Spent = VinMess.Spent
	NewTxDetails.Value = VinMess.Value
	NewTxDetails.TxID = VinMess.Txid
	NewAddress.TxDetails = append(NewAddress.TxDetails, NewTxDetails)
	return NewAddress
}

//----------------------------------------------------------------
func QueryTest(txid string) ([][]Voutinfo, []Voutinfo, uint64) {

	service.GetMongo(mongourl)
	toll := service.GlobalS.DB("GGBTC").C("transaction")
	var q []Tx
	toll.Find(bson.M{"txid": txid}).All(&q)
	//var Txs []Tx
	var vin []*TxIn
	var Vout []Voutinfo
	var Vindels [][]Voutinfo
	//var Txin Txins
	//var VinTx []Txins
	var Time uint64

	//	data, _ := json.Marshal(q)
	//json.Unmarshal(data, &Txs)

	for _, k := range q {
		vin = k.TxIns
		Time = k.BlockTime
	}
	for _, k := range vin {

		Vout = GetVinFromTx(k.Address, k.Index, true)
		Vindels = append(Vindels, Vout)
	}
	Voutinform := GetVoutFromTx(txid, false)
	//fmt.Println("=================================", Voutinform)
	return Vindels, Voutinform, Time
}

func GetVinFromTx(txid string, voutNum uint32, spent bool) []Voutinfo {
	service.GetMongo(mongourl)
	toll := service.GlobalS.DB("GGBTC").C("transaction")
	var q []Tx
	toll.Find(bson.M{"txid": txid}).All(&q)
	//var Txs []AutoTx
	var Vout []*TxOut

	var Voutdel Voutinfo
	var Voutdels []Voutinfo
	//data, _ := json.Marshal(q)
	//json.Unmarshal(data, &Txs)
	//data1, _ := json.Marshal(Txs)
	//fmt.Println(string(data1))
	for _, k := range q {
		Vout = k.TxOuts
	}
	for _, k := range Vout {
		if k.Index == voutNum {
			Voutdel.Address = k.Addr
			// for _, k := range m {
			// 	Voutdel.Address = k
			// }
			Voutdel.Value = k.Value
			Voutdel.Index = int(k.Index)
			Voutdel.Spent = spent
			Voutdel.Txid = txid
			Voutdels = append(Voutdels, Voutdel)
		}
	}

	return Voutdels
}

func GetVoutFromTx(txid string, spent bool) []Voutinfo {
	service.GetMongo(mongourl)
	toll := service.GlobalS.DB("GGBTC").C("transaction")
	var q []Tx
	toll.Find(bson.M{"txid": txid}).All(&q)
	//var Txs []Tx
	var Vout []*TxOut

	var Voutdel Voutinfo
	var Voutdels []Voutinfo
	//	data, _ := json.Marshal(q)
	//	json.Unmarshal(data, &Txs)
	for _, k := range q {
		Vout = k.TxOuts
	}
	for _, k := range Vout {
		Voutdel.Address = k.Addr
		// for _, k := range m {
		// 	Voutdel.Address = k
		// }
		Voutdel.Value = k.Value
		Voutdel.Index = int(k.Index)
		Voutdel.Spent = spent
		Voutdel.Txid = txid
		Voutdels = append(Voutdels, Voutdel)
	}

	return Voutdels
}
