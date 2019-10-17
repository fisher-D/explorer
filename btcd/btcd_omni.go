package btcd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GGBTC/explorer/service"
	"gopkg.in/fatih/set.v0"
)

func OmniProcesser(opreturn string) (*OmniTx, error) {
	Omni := strings.Split(opreturn, " ")[1]
	//fmt.Println(Omni)
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
		return nil, fmt.Errorf("Unable to Transfer Value")
	}
	tar.OP_RETURN = Omni
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

func OmniProcess1(vin []*service.Vin, vouts []*service.VoutNew) string {
	NoOmni := []*service.VoutNew{}
	equal546 := []*service.VoutNew{}
	notequal := []*service.VoutNew{}
	//list3 := []string{}
	for _, k := range vouts {
		if k.Value != 0 {
			NoOmni = append(NoOmni, k)
		}
	}
	//如果出Omni OutPut之外只有一笔交易，则确认为目标结果
	if len(NoOmni) == 1 {
		return NoOmni[0].Addr
	}
	//若有多笔交易，判断546值
	for _, m := range NoOmni {
		//判断有等于546的有几笔
		if m.Value == 546 {
			equal546 = append(equal546, m)
		} else {
			notequal = append(notequal, m)
		}
	}
	if len(equal546) != 0 {
		//若仅有一笔，则确认为目标结果
		if len(equal546) == 1 {
			return equal546[0].Addr
		}
		//移除可能是找零的地址
		noDups := DropDups(equal546, vin)
		if len(noDups) == 1 {
			return noDups[0]
		}
		return "Unknown"

	}
	//若有多笔，进行Vin与Vout地址去重
	//若地址合集长度为1
	if len(notequal) != 0 {
		noDups := DropDups(notequal, vin)
		if len(noDups) == 1 {
			return noDups[0]
		}
		return "Unknown"
	}
	return "FinalUnknown"
}

//}
func DropDups(new []*service.VoutNew, old []*service.Vin) []string {
	Ol := set.New(set.ThreadSafe)
	for _, k := range old {
		Ol.Add(k.Address)
	}

	Ne := set.New(set.ThreadSafe)
	for _, k := range new {
		Ne.Add(k.Addr)
	}
	differ := set.Difference(Ol, Ne)
	res := differ.List()
	tar := make([]string, len(res))
	for v, k := range res {
		tar[v] = k.(string)
	}
	//res := tar
	return tar
	//	fmt.Println(res)
}
