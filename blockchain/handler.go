package blockchain

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	blockChain *Blockchain
}

func NewHandler(blockChain *Blockchain) *Handler {
	return &Handler{
		blockChain: blockChain,
	}
}

func (h *Handler) GetBlockchain(c echo.Context) error {
	return c.JSON(http.StatusOK, GetFullBlockchainResponse{
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
	var request AddTransactionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Missing fields",
		})
	}

	index := h.blockChain.AddTransaction(
		request.Sender,
		request.Recipient,
		request.Amount,
	)

	return c.JSON(http.StatusCreated, map[string]string{
		"message": fmt.Sprintf("Transaction will be added to Block %d", index),
	})
}

func (h *Handler) AddNodes(c echo.Context) error {
	var request AddNodesRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Error Missing node(s) info",
		})
	}

	h.blockChain.AddNodes(request.Nodes)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "New nodes added",
		"nodes":   request.Nodes,
	})
}

func (h *Handler) SyncNodes(c echo.Context) error {
	if updated := h.blockChain.UpdateBlockchain(); updated {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "The blockchain has been updated to the latest",
			"blockchain": h.blockChain.Chain,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "Our blockchain is the latest",
		"blockchain": h.blockChain.Chain,
	})
}
