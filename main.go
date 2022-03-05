package main

import (
	"github.com/sean-ttm/sbtcoin/cli"
	"github.com/sean-ttm/sbtcoin/db"
)

func main() {	
	//defer executes when the main closed
	defer db.Close()
	cli.Start()
}
