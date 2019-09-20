package main

import (
	"github.com/GGBTC/explorer/service"
)

const (
	URL       = "EXAMPLE"
	GenesisTx = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"
	mongourl  = "localhost:27017"
)

func main() {
	// txid := "4e92b3d85db42892e772a5590ecd822f114aae2b09fd058d162df092eaa63cd9"
	// VinMesses, VoutMess, Time := pkgs.QueryTx(txid)
	// Vout := pkgs.ProcessVoutMess(VoutMess)
	// Vin := pkgs.ProcessVinMess(VinMesses)
	// fmt.Println(Vout)
	// VoutAddress := pkgs.SaveAddressData(Vout)
	// VinAddress := pkgs.SaveAddressData(Vin)
	// for _, k := range VinAddress {
	// 	VoutAddress = append(VoutAddress, k)
	// }
	// pkgs.CompleteAddress(VoutAddress, Time)
	// fmt.Println(VoutAddress, Time)
	ClearTx()
}

func ClearTx() {
	// GetUnSpentTransaction(Txid)
	service.GetMongo(mongourl)
	service.GlobalS.DB("GGBTC").C("unspent").DropCollection()
	service.GlobalS.DB("GGBTC").C("transaction").DropCollection()
	service.GlobalS.DB("GGBTC").C("status").DropCollection()
	service.GlobalS.DB("GGBTC").C("address").DropCollection()
	service.GlobalS.DB("GGBTC").C("addresrelatetx").DropCollection()
}
