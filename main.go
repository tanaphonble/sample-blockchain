package main

import (
	"app/blockchain"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	bc := blockchain.New("0000")
	hdlr := blockchain.NewHandler(bc)

	e.GET("/blockchain", hdlr.GetBlockChain)
	e.GET("/mine", hdlr.MineBlock)
	e.POST("/transactions", hdlr.AddTransaction)

	e.Logger.Fatal(e.Start(":1323"))
}
