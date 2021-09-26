package main

import (
	"app/blockchain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()

	e := echo.New()

	nodes := viper.GetStringSlice("NODES")

	bc := blockchain.New("0000", nodes)
	hdlr := blockchain.NewHandler(bc)

	e.Use(middleware.Logger())

	e.GET("/blockchain", hdlr.GetBlockchain)
	e.GET("/mine", hdlr.MineBlock)
	e.POST("/transactions", hdlr.AddTransaction)
	e.POST("/nodes", hdlr.AddNodes)
	e.GET("/nodes/sync", hdlr.SyncNodes)

	e.Logger.Fatal(e.Start(":1323"))
}
