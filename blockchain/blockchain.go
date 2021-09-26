package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"time"

	"github.com/labstack/gommon/log"
)

func New(difficultyTarget string, nodes []string) *Blockchain {
	bc := &Blockchain{
		Chain:               Chain{},
		CurrentTransactions: []Transaction{},
		difficultyTarget:    difficultyTarget,
		nodes:               make(map[string]struct{}),
	}

	genesisHash := bc.HashBlock("genesis_block")
	bc.AppendBlock(bc.ProofOfWork(0, genesisHash, []Transaction{}), genesisHash)
	bc.AddNodes(nodes)

	return bc
}

func (bc *Blockchain) LastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) HashBlock(b interface{}) string {
	blockEncoded, _ := json.Marshal(b)
	hasher := sha256.New()
	hasher.Write(blockEncoded)

	return hex.EncodeToString(hasher.Sum(nil))
}

func (bc *Blockchain) AppendBlock(nonce int, hashOfPreviousBlocks string) Block {
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

func (bc *Blockchain) AddTransaction(sender, recipient string, amount float64) int {
	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})

	return bc.LastBlock().Index + 1
}

func (bc *Blockchain) ProofOfWork(index int, hashOfPreviousBlocks string, transactions []Transaction) int {
	nonce := 0
	for !bc.validProof(index, hashOfPreviousBlocks, transactions, nonce) {
		nonce++
	}

	return nonce
}

func (bc *Blockchain) validProof(index int, hashOfPreviousBlocks string, transactions []Transaction, nonce int) bool {
	transactionsEncoded, _ := json.Marshal(transactions)
	content := fmt.Sprintf("%d%s%s%d", index, hashOfPreviousBlocks, transactionsEncoded, nonce)

	hasher := sha256.New()
	hasher.Write([]byte(content))
	contentHash := hex.EncodeToString(hasher.Sum(nil))

	return contentHash[:len(bc.difficultyTarget)] == bc.difficultyTarget
}

func (bc *Blockchain) AddNodes(nodes []string) {
	for _, node := range nodes {
		bc.nodes[node] = struct{}{}
	}

	fmt.Printf("current nodes: %v", bc.nodes)
}

func (bc *Blockchain) ValidChain(chain []Block) bool {
	lastBlock := chain[0]
	currentIndex := 1

	for currentIndex < len(chain) {
		block := chain[currentIndex]
		if block.HashOfPreviousBlocks != bc.HashBlock(lastBlock) {
			return false
		}

		if !bc.validProof(
			currentIndex,
			block.HashOfPreviousBlocks,
			block.Transactions,
			block.Nonce,
		) {
			return false
		}

		lastBlock = block
		currentIndex++
	}

	return true
}

func (bc *Blockchain) UpdateBlockchain() bool {
	var newChain []Block
	maxLength := len(bc.Chain)

	for node := range bc.nodes {
		url := fmt.Sprintf("http://%s/blockchain", node)
		fmt.Println("URL:", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("cannot get blockchain at node %s: %v", node, err)
			return false
		}

		if resp.StatusCode != 200 {
			return false
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}

		var bcResp GetFullBlockchainResponse
		if err := json.Unmarshal(body, &bcResp); err != nil {
			return false
		}

		length := bcResp.Length
		chain := bcResp.Chain

		if length > maxLength && bc.ValidChain(chain) {
			maxLength = length
			newChain = chain
		}
	}

	if newChain != nil {
		bc.Chain = newChain
		return true
	}

	return false
}
