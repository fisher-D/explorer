package pkgs

import (
	s "github.com/GGBTC/explorer/service"
)

func DataName() []string {
	//url := "192.168.72.250:27017"
	s.GetMongo(s.Mongourl)
	session, _ := s.GlobalS.DatabaseNames()
	CoinName := []string{}
	for _, k := range session {
		if k != "admin" && k != "config" && k != "local" {
			CoinName = append(CoinName, k)
		}
	}
	if len(CoinName) != 0 {
		return CoinName
	}
	return []string{"No DataBase Exist"}
}
