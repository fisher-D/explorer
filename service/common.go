package service

import (
	"time"
)

type Status struct {
	//LastBlockHeight uint64
	//	TxCount         uint64
	//AddressCount    uint64
	//	UTXOcount       uint64
	//ID        bson.ObjectId `bosn:"id"`
	Currecny  string `json:"currency"`
	UtxoBlock uint64 `json:"ublock"`
	UtxoTx    string `json:"utx"`
}
type Address struct {
	Address        string `json:"address"`
	Txs            []Txs  `json:"Tx"`
	TotalRecCount  int
	TotalSentCount int
	TotalReceived  uint64
	TotalSent      uint64
	Balance        uint64
	FirstSeen      uint64
	LastSeen       uint64 `json:"lastseen"`
}
type Txs struct {
	Index    uint32 `json:"index"`
	Txid     string `json:"txid"`
	Value    uint64 `json:"value"`
	Spent    string `json:"spent"`
	Currency string `json:"currency"`
}
type UTXO struct {
	//Id       bson.ObjectId `json:"id"        bson:"_id"`
	Index    uint32 `json:"index"`
	Utxo     string `json:"utxo"`
	Address  string `json:"address"`
	Value    uint64 `json:"value"`
	Currency string `json:"currency"`
	Spent    string `json:"spent,omitempty"`
}
type Tx struct {
	Txid      string     `json:"txid"`
	Version   uint32     `json:"version"`
	BlockHash string     `json:"blockhash"`
	BlockTime uint64     `json:"blocktime"`
	Totalin   uint64     `json:"valuein"`
	Totalout  uint64     `json:"valueout"`
	Type      string     `json:"type"`
	Fee       uint64     `json:"fee"`
	Vin       []*Vin     `json:"Vin"`
	Vout      []*VoutNew `json:"Vout"`
}

type VoutNew struct {
	Index    uint32 `json:"index"`
	Addr     string `json:"address"`
	Value    uint64 `json:"value"`
	Currency string `json:"currency"`
	Spent    string `json:"spent"`
}

type Vin struct {
	Hash     string `json:"prevout,omitempty"`
	Address  string `json:"address,omitempty"`
	Value    uint64 `json:"value,omitempty"`
	Index    uint32 `json:"index,omitempty"`
	Coinbase string `json:"coinbase,omitempty"`
	Sequence int64  `json:"sequence,omitempty"`
	Currency string `json:"currency"`
	Spent    string `json:"spent"`
}

type TxOld struct {
	Txid      string `json:"txid"`
	Blockhash string `json:"blockhash"`
	Blocktime uint64 `json:"blocktime"`
	Version   uint64 `json:"version"`
	Vin       []Vin  `json:"vin"`
	Vout      []Vout `json:"vout"`
}

type ScriptPubKey struct {
	Addresses []string `json:"addresses"`
	Asm       string   `json:"asm"`
	Type      string   `json:"type"`
}

type Vout struct {
	N            uint64       `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
	Value        uint64       `json:"value"`
}
type Blocks struct {
	Height     int      `json:"height"`
	Hash       string   `json:"hash"`
	Difficulty uint64   `json:"difficulty"`
	Version    int      `json:"version"`
	Time       int      `json:"time"`
	NTx        int      `json:"NTx"`
	Tx         []string `json:"tx"`
}

// Float rounding to precision 8
func FloatToUint(x float64) uint64 {
	return uint64(int64((x * float64(100000000.0)) + float64(0.5)))
}

type Information struct {
	Price        float64
	MarketCap    float64
	MarketAmount float64
	Amount       int
	Height       int
	Difficult    uint64
}
type BTCInfo struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
		Notice       interface{} `json:"notice"`
	} `json:"status"`
	Data []struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Symbol            string      `json:"symbol"`
		Slug              string      `json:"slug"`
		NumMarketPairs    int         `json:"num_market_pairs"`
		DateAdded         time.Time   `json:"date_added"`
		Tags              []string    `json:"tags"`
		MaxSupply         int         `json:"max_supply"`
		CirculatingSupply int         `json:"circulating_supply"`
		TotalSupply       int         `json:"total_supply"`
		Platform          interface{} `json:"platform"`
		CmcRank           int         `json:"cmc_rank"`
		LastUpdated       time.Time   `json:"last_updated"`
		Quote             struct {
			USD struct {
				Price            float64   `json:"price"`
				Volume24H        float64   `json:"volume_24h"`
				PercentChange1H  float64   `json:"percent_change_1h"`
				PercentChange24H float64   `json:"percent_change_24h"`
				PercentChange7D  float64   `json:"percent_change_7d"`
				MarketCap        float64   `json:"market_cap"`
				LastUpdated      time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}
