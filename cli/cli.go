package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/sean-ttm/sbtcoin/explorer"
	"github.com/sean-ttm/sbtcoin/rest"
)


func usage(){
	fmt.Printf("Welcome to SBT Coin\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port=4000: 	Set the PORT of the server\n")
	fmt.Printf("-mode=rest: 	Choose among 'html', 'rest', or 'both'\n")
	fmt.Printf("-port2=5000:	Set the PORT of the html (if both)")
	runtime.Goexit()
}

func Start(){
		//os.Args - string slices로 command-line arguments 받아옴
		if len(os.Args) ==1 {
			usage()	
		}
	
		port := flag.Int("port", 4000, "Set port of the server")
		port2 := flag.Int("port2", 6000, "Set port of the html server (if both)")
		mode := flag.String("mode", "rest", "Choose among 'html', 'rest', or both")
		flag.Parse()
	
		switch *mode{
		case "rest":
			rest.Start(*port)
		case "html":
			explorer.Start(*port)
		case "both":
			rest.Start(*port)
			explorer.Start(*port2)
		default:
			usage()
		}
}