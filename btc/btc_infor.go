package btc

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type ChainInfor struct {
	CurrentPrice uint64 `json:"currentprice"`
	MarketPrice  uint64 `json:"marketprice"`
	TxAmount     uint64 `json:"txAmount"`
	TotalPrice   uint64 `json:"totalPirce"`
	Height       uint64 `json:"height"`
	Difficulty   uint64 `json:"difficult"`
}

//GetLastBitcoinPrice Get BTC Pirce
func GetLastBitcoinPrice() (price float64, err error) {
	resp, err := http.Get("https://api.bitcoinaverage.com/ticker/global/USD/last")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	price, err = strconv.ParseFloat(string(body), 10)
	return
}

// func main() {
// 	p, _ := GetLastBitcoinPrice()
// 	fmt.Printf("%v", p)
// }
