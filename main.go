package main

import (
	"log"

	zec "github.com/GGBTC/explorer/zcash"
)

func main() {
	height := zec.GetzecCountRPC()
	log.Println(height)
}
