package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//zec "github.com/GGBTC/explorer/zcash"

func main() {
	//btc.GetLastBitCoinPrice()
	Price()
}

// Return last USD price from BitcoinAverage API
func Price() (price float64, err error) {
	resp, err := http.Get("https://coinmarketcap.com/watchlist/")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
	//price, err = strconv.ParseFloat(string(body), 10)
	return
}
