package service

type TxDetail struct {
	TxID  string
	Value uint64
	Index int
	Spent bool
}
type OrigonalBlock struct {
	Hash              string   `bson:"hash"`
	Confirmations     uint32   `bson:"confirmations"`
	Strippedsize      uint32   `bson:"strippedsize"`
	Size              uint32   `bson:"size"`
	Weight            uint32   `bson:"weight"`
	Height            uint32   `bson:"height"`
	Version           uint32   `bson:"version"`
	VersionHex        string   `bson:"versionHex"`
	Merkleroot        string   `bson:"merkleroot"`
	Tx                []string `bson:"tx"`
	Time              uint32   `bson:"time"`
	Mediantime        uint32   `bson:"mediantime"`
	Nonce             uint64   `bson:"nonce"`
	Bits              string   `bson:"bits"`
	Difficulty        float64  `bson:"difficulty"`
	Chainwork         string   `bson:"chainwork"`
	NTx               uint32   `bson:"nTx"`
	Previousblockhash string   `bson:"previousblockhash"`
	Nextblockhash     string   `bson:"nextblockhash"`
}
