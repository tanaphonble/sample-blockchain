package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type BlockChain struct {
	difficultyTarget    string
	Chain               []Block
	CurrentTransactions []Transaction
}

type Block struct {
	Index                int           `json:"index"`
	Timestamp            time.Time     `json:"timestamp"`
	Transactions         []Transaction `json:"transactions"`
	Nonce                int           `json:"nonce"`
	HashOfPreviousBlocks string        `json:"hash_of_previous_blocks"`
}

type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

func New(difficultyTarget string) *BlockChain {
	bc := &BlockChain{
		Chain:               []Block{},
		CurrentTransactions: []Transaction{},
		difficultyTarget:    difficultyTarget,
	}

	genesisHash := bc.HashBlock("genesis_block")
	bc.AppendBlock(bc.ProofOfWork(0, genesisHash, []Transaction{}), genesisHash)

	return bc
}

func (bc *BlockChain) LastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *BlockChain) HashBlock(b interface{}) string {
	blockEncoded, _ := json.Marshal(b)
	hasher := sha256.New()
	hasher.Write(blockEncoded)

	return hex.EncodeToString(hasher.Sum(nil))
}

func (bc *BlockChain) AppendBlock(nonce int, hashOfPreviousBlocks string) Block {
	b := Block{
		Index:                len(bc.Chain),
		Timestamp:            time.Now().UTC(),
		Transactions:         bc.CurrentTransactions,
		Nonce:                nonce,
		HashOfPreviousBlocks: hashOfPreviousBlocks,
	}

	bc.CurrentTransactions = []Transaction{}
	bc.Chain = append(bc.Chain, b)

	return b
}

func (bc *BlockChain) AddTransaction(sender, recipient string, amount float64) int {
	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})

	return bc.LastBlock().Index + 1
}

func (bc *BlockChain) ProofOfWork(index int, hashOfPreviousBlocks string, transactions []Transaction) int {
	nonce := 0
	for !bc.validProof(index, hashOfPreviousBlocks, transactions, nonce) {
		nonce++
	}

	return nonce
}

func (bc *BlockChain) validProof(index int, hashOfPreviousBlocks string, transactions []Transaction, nonce int) bool {
	transactionsEncoded, _ := json.Marshal(transactions)
	content := fmt.Sprintf("%d%s%s%d", index, hashOfPreviousBlocks, transactionsEncoded, nonce)

	hasher := sha256.New()
	hasher.Write([]byte(content))
	contentHash := hex.EncodeToString(hasher.Sum(nil))

	return contentHash[:len(bc.difficultyTarget)] == bc.difficultyTarget
}
