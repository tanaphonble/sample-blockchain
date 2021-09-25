package blockchain

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	blockChain *BlockChain
}

func NewHandler(blockChain *BlockChain) *Handler {
	return &Handler{
		blockChain: blockChain,
	}
}

func (h *Handler) GetBlockChain(c echo.Context) error {
	return c.JSON(http.StatusOK, GetFullBlockChainResponse{
		Chain:  h.blockChain.Chain,
		Length: len(h.blockChain.Chain),
	})
}

func (h *Handler) MineBlock(c echo.Context) error {
	lastBlockHash := h.blockChain.HashBlock(h.blockChain.LastBlock())
	index := len(h.blockChain.Chain)
	nonce := h.blockChain.ProofOfWork(index, lastBlockHash, h.blockChain.CurrentTransactions)

	block := h.blockChain.AppendBlock(nonce, lastBlockHash)

	return c.JSON(http.StatusOK, MineBlockResponse{
		Message: "New Block Mined",
		Block:   block,
	})
}

func (h *Handler) AddTransaction(c echo.Context) error {
	var req AddTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Missing fields",
		})
	}

	index := h.blockChain.AddTransaction(
		req.Sender,
		req.Recipient,
		req.Amount,
	)

	return c.JSON(http.StatusCreated, map[string]string{
		"message": fmt.Sprintf("Transaction will be added to Block %d", index),
	})
}
