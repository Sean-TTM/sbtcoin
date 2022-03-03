package main

import (
	"github.com/sean-ttm/sbtcoin/explorer"
	"github.com/sean-ttm/sbtcoin/rest"
)


func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
