package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sean-ttm/sbtcoin/blockchain"
)
const (
	port string = ":4000"
	templateDir string = "explorer/templates/"
)

type homeData struct {
	PageTitle string
	Blocks []*blockchain.Block
}

//http.ResponseWriter - response data , pointing data - http.Request
func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home",data)
}

 func add(rw http.ResponseWriter, r *http.Request){
	 switch r.Method {
		case "GET":
			templates.ExecuteTemplate(rw, "add", nil)
		case "POST":
			r.ParseForm()
			data := r.Form.Get("blockData")
			blockchain.GetBlockchain().AddBlock(data)
			http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	 }
 }

var templates *template.Template

func Start(){
		//go는 **/*.gohtml 안됨, 폴더 지정 해줘야 함
		templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
		//initialize 했으니, templates로 쓰면 됨
		templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
		//set up
		http.HandleFunc("/", home) 
		http.HandleFunc("/add", add)
		fmt.Printf("Listening on http:localhost%s\n", port)
		//this is how to start server in Go. log.Fatal logs error(1)
		log.Fatal(http.ListenAndServe(port, nil))
}