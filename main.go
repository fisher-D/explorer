package main

import (
	// "fmt"
	// "io/ioutil"
	// "net/http"
	"log"

	"github.com/GGBTC/explorer/btc"
)

//zec "github.com/GGBTC/explorer/zcash"

func main() {
	res := btc.GetLastBitCoinPrice()
	log.Println(res)
	//return true
}
