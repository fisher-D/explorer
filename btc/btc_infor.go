package btc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	s "github.com/GGBTC/explorer/service"
)

func GetLastBitCoinPrice() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "1")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apikey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	//fmt.Println(resp.Status)
	var infor s.BTCInfo
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBody, &infor)
	//币值
	Price := infor.Data[0].Quote.USD.Price
	//币总额
	Amount := 21000000
	MarketCap := infor.Data[0].Quote.USD.MarketCap
	MaketAmount := infor.Data[0].Quote.USD.Volume24H
	Amou := MaketAmount / Price
	var Info s.Information
	Info.Price = Price
	Info.Amount = Amount
	Info.MarketCap = MarketCap
	Info.MarketAmount = Amou
	s.GetMongo(mongourl)
	Database := s.GlobalS.DB("BTC")
	Height, Difficult := GetBlockInfo(Database)
	Info.Height = Height
	Info.Difficult = Difficult
	collection := Database.C("info")
	collection.Remove(nil)
	err = collection.Insert(Info)
	if err != nil {
		log.Println("Fuck")
		return "Failed"
	}
	return "Success"
}
func GetBlockInfo(Database *mgo.Database) (int, uint64) {
	boll := Database.C("blocks")
	h, _ := boll.Count()
	var height int
	if h <= 1 {
		height = 0
	} else {
		height = h - 1
	}
	var blockinfo s.Blocks
	boll.Find(bson.M{"height": height}).One(&blockinfo)
	Diff := blockinfo.Difficulty
	return height, Diff
}
