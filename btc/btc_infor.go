package btc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/GGBTC/explorer/service"
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
	var infor BTCInfo
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBody, &infor)
	//币值
	Price := infor.Data[0].Quote.USD.Price
	//币总额
	Amount := 21000000
	MarketCap := infor.Data[0].Quote.USD.MarketCap
	MaketAmount := infor.Data[0].Quote.USD.Volume24H
	Amou := MaketAmount / Price
	var Info Information
	Info.Price = Price
	Info.Amount = Amount
	Info.MarketCap = MarketCap
	Info.MarketAmount = Amou
	service.GetMongo(mongourl)
	Database := service.GlobalS.DB("BTC")
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
	var blockinfo service.Blocks
	boll.Find(bson.M{"height": height}).One(&blockinfo)
	Diff := blockinfo.Difficulty
	return height, Diff
}

type Information struct {
	Price        float64
	MarketCap    float64
	MarketAmount float64
	Amount       int
	Height       int
	Difficult    uint64
}
type BTCInfo struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
		Notice       interface{} `json:"notice"`
	} `json:"status"`
	Data []struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Symbol            string      `json:"symbol"`
		Slug              string      `json:"slug"`
		NumMarketPairs    int         `json:"num_market_pairs"`
		DateAdded         time.Time   `json:"date_added"`
		Tags              []string    `json:"tags"`
		MaxSupply         int         `json:"max_supply"`
		CirculatingSupply int         `json:"circulating_supply"`
		TotalSupply       int         `json:"total_supply"`
		Platform          interface{} `json:"platform"`
		CmcRank           int         `json:"cmc_rank"`
		LastUpdated       time.Time   `json:"last_updated"`
		Quote             struct {
			USD struct {
				Price            float64   `json:"price"`
				Volume24H        float64   `json:"volume_24h"`
				PercentChange1H  float64   `json:"percent_change_1h"`
				PercentChange24H float64   `json:"percent_change_24h"`
				PercentChange7D  float64   `json:"percent_change_7d"`
				MarketCap        float64   `json:"market_cap"`
				LastUpdated      time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}
