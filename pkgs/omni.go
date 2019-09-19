package pkgs

import (
	"fmt"
	"strconv"
)

// func main() {
// 	//	fmt.Println(1 + 1)
// 	Omni := "6f6d6e69000000000000001f0000000649534e00"
// 	res, err := OmniProcesser(Omni)
// 	fmt.Println(err)
// 	fmt.Println(res)
// }
func OmniProcesser(Omni string) (*OmniTx, error) {

	tar := new(OmniTx)
	if len(Omni) != 40 {

		return nil, fmt.Errorf("Unknown Type of Omni String")
	}
	if Omni[8:12] == "0000" {
		tar.Version = "0"
	} else {
		tar.Version = "Unknown"
	}
	if Omni[12:16] == "0000" {
		tar.TxType = "Simple Send"
	} else {
		tar.TxType = "Unknown"
	}
	if (Omni[22:24] == "1f") || (Omni[22:24] == "31") {
		tar.TokenName = "TetherUS"
	} else {
		tar.TokenName = "Unknown"
	}
	//int, err := strconv.Atoi(string)

	Value16 := Omni[24:len(Omni)]
	n, err := strconv.ParseUint(Value16, 16, 64)
	if err != nil {
		panic(err)
	}
	tar.Value = n
	return tar, nil
}

type OmniTx struct {
	OP_RETURN string
	TokenName string
	Version   string
	TxType    string
	Value     uint64
}
