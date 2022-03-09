package main

import "github.com/sean-ttm/sbtcoin/cli"

func main() {	
	//defer executes when the main closed
	//defer db.Close()
	cli.Start()
}

