package blockchain

import "time"

type Blockchain struct {
	nodes map[string]struct{}

	difficultyTarget    string
	Chain               Chain
	CurrentTransactions []Transaction
}

type Block struct {
	Index                int           `json:"index"`
	Timestamp            time.Time     `json:"timestamp"`
	Transactions         []Transaction `json:"transactions"`
	Nonce                int           `json:"nonce"`
	HashOfPreviousBlocks string        `json:"hash_of_previous_blocks"`
}

type Chain []Block

type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}
