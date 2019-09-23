package pkgs

//TODO
//Test the VinRemove And VoutInsert
// type Unspent struct {
// 	Txid     string
// 	Address  string
// 	Value    uint64
// 	Currency string
// 	Index    uint32
// }

// // func CountTx() {
// // 	session := service.GlobalS.DB("LTC").C("txbyheight")
// // 	unspetsess := service.GlobalS.DB("LTC").C("unspent")

// // }
// func GetUnSpentTransaction(Txid string) error {
// 	session := service.GlobalS.DB("LTC").C("transaction")
// 	var q Tx
// 	session.Find(bson.M{"txid": Txid}).One(&q)
// 	//fmt.Println(q)
// 	var Vin []*Vin
// 	var Vout []*Vout
// 	Vin = q.TxIns
// 	Vout = q.TxOuts
// 	UTXOI := drawVinOut(Vin)
// 	UTXOT := drawVoutOut(Vout, Txid)
// 	//fmt.Println(UTXOI)
// 	RemoveVin(UTXOI)

// 	InsertVout(UTXOT)

// 	return nil
// 	//fmt.Println("========================")
// 	//fmt.Println(UTXOT)

// }

// func drawVinOut(Vin []*TxIn) *Unspent {
// 	UTXO := new(Unspent)
// 	if Vin == nil {
// 		return nil
// 	}

// 	if Vin[0].Coinbase != "" {
// 		return nil
// 	}
// 	for _, k := range Vin {
// 		UTXO.Address = k.Address
// 		UTXO.Txid = k.Hash
// 		UTXO.Value = k.Value
// 		UTXO.Currency = k.Currency
// 		UTXO.Index = k.Index
// 		return UTXO
// 	}
// 	return nil
// }
// func drawVoutOut(Vout []*TxOut, Txid string) *Unspent {
// 	UTXO := new(Unspent)
// 	for _, k := range Vout {
// 		UTXO.Address = k.Addr
// 		UTXO.Txid = Txid
// 		UTXO.Value = k.Value
// 		UTXO.Currency = k.Currency
// 		UTXO.Index = k.Index
// 		return UTXO
// 	}
// 	return nil
// }

// func InsertVout(UTXOT *Unspent) {
// 	if UTXOT != nil {

// 		session := service.GlobalS.DB("LTC").C("unspent")
// 		session.Insert(UTXOT)
// 		//return nil
// 	} else {
// 		//return fmt.Errorf("Coinbase Tx does not need to be instert")
// 		fmt.Println("Coinbase Tx does not need to be instert")
// 	}
// }
// func RemoveVin(UTXOI *Unspent) {
// 	if UTXOI != nil {
// 		session := service.GlobalS.DB("LTC").C("unspent")
// 		session.Remove(UTXOI)
// 		//return nil
// 	} else {
// 		//return fmt.Errorf("Coinbase Tx does not need to be instert")
// 		fmt.Println("Coinbase Tx does not need to be instert")
// 	}
// }
