package btc

import (
	"fmt"

	"github.com/GGBTC/explorer/service"
	"gopkg.in/mgo.v2"
)

const (
	LTCurl    = "Example"
	URL       = "EXAMPLE"
	GenesisTx = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"
	mongourl  = "localhost:27017"
)

type BTCconfig struct {
	Name      string
	URL       string
	GenesisTx string
	DataBase  *mgo.Database
}

// type LTCconfig struct {
// 	Name      string
// 	URL       string
// 	GenesisTx string
// 	mongourl  string
// }
type Config struct {
	BTC BTCconfig
	//LTC LTCconfig
}

func CreateDB() *mgo.Database {
	service.GetMongo(mongourl)
	database := service.GlobalS.DB("BTC")
	return database
}
func BuildConfig(url, genesisTx string) *Config {
	config := new(Config)
	config.BTC.Name = "BTC"
	config.BTC.URL = url
	config.BTC.GenesisTx = genesisTx
	config.BTC.DataBase = CreateDB()
	return config
}

func main() {
	configs := BuildConfig(URL, GenesisTx)
	btc := configs.BTC
	fmt.Println(btc)
}
