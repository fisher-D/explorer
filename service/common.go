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
	TxID     string
	Index    uint32
	Address  string
	Value    uint64
	Spent    bool
	Currency string
}
type Tx struct {
	Txid      string `json:"txid"`
	Version   uint32 `json:"version"`
	BlockHash string `json:"blockhash"`
	//BlockHeight uint       `json:"blockheight"`
	BlockTime uint64     `json:"blocktime"`
	Vin       []*Vin     `json:"Vin"`
	Vout      []*VoutNew `json:"Vout"`
}

type VoutNew struct {
	Addr     string `json:"address"`
	Value    uint64 `json:"value"`
	Index    uint32 `json:"index"`
	Currency string `json:"currency"`
}

type Vin struct {
	Hash     string `json:"prevout,omitempty"`
	Address  string `json:"address,omitempty"`
	Value    uint64 `json:"value,omitempty"`
	Index    uint32 `json:"index,omitempty"`
	Coinbase string `json:"coinbase,omitempty"`
	Sequence int64  `json:"sequence,omitempty"`
	Currency string `json:"currency"`
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
	//Hex string `json:"hex"`
	//	ReqSigs   int      `json:"reqSigs"`
	Type string `json:"type"`
}
type Vout struct {
	N            uint64       `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
	Value        uint64       `json:"value"`
	//Currency     string       `json:"currency"`
}
type Blocks struct {
	Height int    `json:"height"`
	Hash   string `json:"hash"`
	//	Confirmations int      `json:"confirmations"`
	Difficulty uint64 `json:"difficulty"`
	Version    int    `json:"version"`
	//VersionHex string   `json:"versionHex"`
	Time int      `json:"time"`
	NTx  int      `json:"NTx"`
	Tx   []string `json:"tx"`
}
