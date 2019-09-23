package pkgs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/GGBTC/explorer/service"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Temps struct {
	//ID     int
	Height uint32
}
type TxArrays struct {
	Id uint32
	Tx []string
}

// Fetch a transaction via bticoind RPC API
// func GetTxRPC(txid string, height int) (he map[string]interface{}, err error) {
// 	// Hard coded genesis tx since it's not included in bitcoind RPC API
// 	if txid == GenesisTx {
// 		return
// 		//return TxData{GenesisTx, []TxIn{}, []TxOut{{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 5000000000}}}, nil
// 	}
// 	// Get the TX from bitcoind RPC API
// 	res_tx, err := CallBitcoinRPC(URL, "getrawtransaction", 1, []interface{}{txid, 1})
// 	if err != nil {
// 		log.Fatalf("Err: %v", err)
// 	}
// 	var txjson map[string]interface{}
// 	txjson = res_tx["result"].(map[string]interface{})
// 	//var he map[string]interface{}
// 	he = make(map[string]interface{})
// 	he = txjson
// 	he["BlockHeight"] = height
// 	return he, nil

// }
func BuildTxsArry(txarr *mgo.Collection, height int) ([]string, uint32, error) {
	var res []TxArrays
	//var q []bson.M
	err := txarr.Find(bson.M{"id": height}).All(&res)
	//data, _ := json.Marshal(q)
	//json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println("Unable to Query Block At height:", height)
		return nil, 0, nil
	}
	if res == nil {
		fmt.Println("Already finish")
		return nil, uint32(height), err
	}
	Height := res[0].Id
	Txs := res[0].Tx
	return Txs, Height, nil

}

func SaveTxRPC(Tx *service.Tx) string {
	// var Txid string
	// //fmt.Sprint(Tx)
	// if Tx != nil {
	// 	Txid = Tx.TxID
	// } else {
	// 	Tx = nil
	// }
	// fmt.Println(time.Now(), "Processing Transactions ||TxID", Txid)
	err := service.Insert("LTC", "transaction", Tx)
	if err != nil {
		return "Failed1"
	}
	//err = GetAddressRPC(Tx.TxID)
	//if err != nil {
	//	return "Failed2"
	//}
	//err = GetUnSpentTransaction(Tx.TxID)
	//if err != nil {
	//	return "Faild3"
	//}
	//
	return "Success"
}

func QueryTxByHeight(height int) {
	query := bson.M{"BlockHeight": height}
	TxQueryOptions(query)
}

func QueryTxByHash(hash string) {
	query := bson.M{"txid": hash}
	TxQueryOptions(query)
}

func TxQueryOptions(query interface{}) {

	var session *mgo.Session
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	boll := session.DB("LTC").C("transaction")
	var q []bson.M
	boll.Find(query).All(&q)
	b, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}

	// for a, info := range q {
	// 	b, err := json.Marshal(info)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//fmt.Println("Result Count :", a+1)
	fmt.Println("Result:", string(b))

	//}
}

//CheckBlockHeightInMongo
//Get the lastest block height in mongodb
//Remove it in case that record is incomplete
//Get tx will start from that height

func CheckBlockHeightInMongo() (int, int) {
	var session *mgo.Session
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	boll := session.DB("LTC").C("blocks")
	toll := session.DB("LTC").C("transaction")
	status := session.DB("LTC").C("status")
	var a []bson.M
	var startTime int
	status.Find(nil).All(&a)
	if a == nil {
		startTime = 0
	} else {
		startTime = 1
	}
	count, _ := boll.Count()
	var q []bson.M
	quer := bson.M{"BlockHeight": startTime}
	toll.Find(quer).All(&q)
	//res, _ := TxQuery(coll)
	if len(q) == 0 {
		fmt.Println("Catched Up Tx in Block Height: ", startTime)
		toll.Remove(bson.M{"BlockHeight": startTime})
		return startTime, count

	}
	return 0, 0
}

func AutoDrawTx() string {

	boll := service.GlobalS.DB("LTC").C("blocks")
	lastNum, _ := boll.Count()
	fmt.Println(lastNum)
	startNum, _ := service.GlobalS.DB("LTC").C("txbyheight").Count()
	fmt.Println(startNum)
	if startNum > 1 {
		service.GlobalS.DB("LTC").C("txbyheight").Remove(bson.M{"id": startNum})
	}
	for i := startNum; i <= lastNum; i++ {
		DrawOutTxs(boll, i)
	}
	return "Success"

}
func DrawOutTxs(boll *mgo.Collection, order int) {
	//lastNum, _ := boll.Count()
	var TT TxArrays
	var q []service.OrigonalBlock
	boll.Find(bson.M{"height": order}).All(&q)
	for _, k := range q {
		TT.Id = k.Height
		TT.Tx = k.Tx
		session := service.GlobalS.DB("LTC").C("txbyheight")
		session.Insert(TT)
		fmt.Println("Success To Build txbyheight for BlockHeight:", order)
	}

}

//main function
func CatchUpTx() string {
	service.GetMongo(mongourl)

	txarr := service.GlobalS.DB("LTC").C("txbyheight")
	res, err := GetAndSaveTx(txarr)
	if err != nil {
		fmt.Println("Error happened")
	}
	//fmt.Println(res)
	return res
}

func GetStartTime() (int, int) {
	service.GetMongo(mongourl)
	temp := service.GlobalS.DB("LTC").C("status")
	tar := service.GlobalS.DB("LTC").C("txbyheight")
	var test Temps
	var tt Temps
	//test.Height = 32
	temp.Insert(test)
	temp.Find(nil).One(&tt)
	time := int(tt.Height)
	endtime, _ := tar.Count()
	return time, endtime
}

//GetAndSaveTx ...
func GetAndSaveTx(txarr *mgo.Collection) (string, error) {
	startNum, endNum := GetStartTime()
	var tem Temps
	//fmt.Println(startNum, endNum)
	for i := startNum + 1; i <= endNum; i++ {
		Txs, Height, err := BuildTxsArry(txarr, i)
		if err != nil {
			//	fmt.Println("Unable to Get Transactions for Block At height", Height)
			return "", err
		}
		countTx := len(Txs)
		//fmt.Println("Blockheight :", Height, "||Has:", countTx, "Number of Txs")
		if countTx != 0 {
			tem.Height = Height
			service.GlobalS.DB("LTC").C("status").Update(nil, tem)
			for i := 0; i < countTx; i++ {
				//Test Propuse Delete Height
				rawTx := GetClearTx(Txs[i])
				// if err != nil {
				// 	fmt.Println("Unable to Get Transaction,Tx Hash:", Txs[i])
				// 	return "", err
				// }
				result := SaveTxRPC(rawTx)
				fmt.Println(result)
				//if result != "Success" {
				//	return "Faild", fmt.Errorf("Error Happened")
				//}

			}

		}
	}
	return "Success", nil
}

// Helper to make call to bitcoind RPC API
func CallBitcoinRPC(address string, method string, id interface{}, params []interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return nil, err
	}
	resp, err := http.Post(address,
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalf("ReadAll: %v", err)
	//	return nil, err
	//}
	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&result)
	//err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}
	return result, nil
}
