package pkgs

import (
	"encoding/json"
	"fmt"

	"github.com/GGBTC/explorer/service"
)

func GetBlockCountRPC() int {
	res, err := service.CallBitcoinRPC(URL, "getblockcount", 1, []interface{}{})
	if err != nil {
		fmt.Println("Error")
	}
	count, _ := res["result"].(json.Number).Int64()
	//jsoninfo := res["result"].(map[string]interface{})
	return int(count)
}

//GetBlockHashRPC Done
func GetBlockHashRPC(height int) string {
	// Get the block hash
	res, err := service.CallBitcoinRPC(URL, "getblockhash", 1, []interface{}{height})
	if err != nil {
		fmt.Println("error")
	}
	tar := res["result"].(string)
	return tar
}
