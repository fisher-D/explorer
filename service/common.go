package service

type Address struct {
	Address        string `json:"address"`
	Txs            []Txs  `json:"Tx"`
	TotalRecCount  int
	TotalSentCount int
	TotalReceived  uint64
	TotalSent      uint64
	Balance        uint64
	FirstSeen      uint64
	LastSeen       uint64
}
type Txs struct {
	Index    uint32 `json:"index"`
	Txid     string `json:"txid"`
	Value    uint64 `json:"value"`
	Spent    bool   `json:"spent"`
	Currency string `json:"currency"`
}
type UTXO struct {
	Index    uint32 `json:"index"`
	Utxo     string `json:"utxo"`
	Address  string `json:"address"`
	Value    uint64 `json:"value"`
	Currency string `json:"currency"`
	Spent    bool   `json:"spent"`
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
	Spent    bool   `json:"spent"`
}

type Vin struct {
	Hash     string `json:"prevout,omitempty"`
	Address  string `json:"address,omitempty"`
	Value    uint64 `json:"value,omitempty"`
	Index    uint32 `json:"index,omitempty"`
	Coinbase string `json:"coinbase,omitempty"`
	Sequence int64  `json:"sequence,omitempty"`
	Currency string `json:"currency"`
	Spent    bool   `json:"spent"`
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
