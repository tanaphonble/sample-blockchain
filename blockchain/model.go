package blockchain

type AddTransactionRequest struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

type GetFullBlockChainResponse struct {
	Chain  []Block `json:"chain"`
	Length int     `json:"length"`
}

type MineBlockResponse struct {
	Message string `json:"message"`
	Block   Block  `json:"block,inline"`
}
