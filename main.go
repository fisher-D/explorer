package main

import (
	"fmt"

	s "github.com/GGBTC/explorer/service"
)

func main() {
	//url := "192.168.72.250:27017"
	s.GetMongo(s.Mongourl)
	session, _ := s.GlobalS.DatabaseNames()
	CoinName := []string{}
	for _, k := range session {
		if k != "admin" && k != "config" && k != "local" {
			CoinName = append(CoinName, k)
		}
	}
	fmt.Println(CoinName)
}
