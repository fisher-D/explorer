package pkgs

import (
	"fmt"
	"strconv"
)

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
	para := Omni[22:24]
	switch para {
	case "if":
		tar.TokenName = "TetherUS"
	case "31":
		tar.TokenName = "TetherUS"
	default:
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
