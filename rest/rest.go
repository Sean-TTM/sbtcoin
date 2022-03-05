package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sean-ttm/sbtcoin/blockchain"
	"github.com/sean-ttm/sbtcoin/utils"
)

var port string

//PORT 에서 message로 호출하는 것 받아오기 위함
type addBlockBody struct {
	Message string
}

type url string

//Json으로 encoding할 때, url 부분만 별도로 실행하게 함
func (u url) MarshalText() ([]byte, error){
	url := fmt.Sprintf("http:localhost%s%s", port, u)
	return []byte(url), nil
}

type URLDescription struct {
	//field struc tag: `jason: "name"`
	//json 에서는 저렇게 표현한다 라는 뜻임
	URL 		url	   `json:"url"`
	Method 		string `json:"method"`
	Description string `json:"description"`
	//,omitempty: hide when empty (띄어쓰기 안해야됨 주의)
	Payload 	string `json:"payload,omitempty"`
	//`json:"-"` : ignore json
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

//String()는 내장되어 있어서, 별도 implement 없이 바로 사용 가능
//Stringer interface - URLDescription Format 할 때 method 호출됨
// func (u URLDescription) String() string {
//	return "Hello I'm the URL Description"
//}

func documentation(rw http.ResponseWriter, r *http.Request){
	data := []URLDescription{
		{
			URL: 			url("/"),
			Method: 		"GET",
			Description:	"See Documentation",
		},
		{
			URL: 			url("/blocks"),
			Method: 		"POST",
			Description: 	"Add a Block",
			Payload: 		"data:string",
		},
		{
			URL: 			url("/blocks/{hash}"),
			Method: 		"GET",
			Description: 	"See a Block",
		},
	}
	//sending json response (middleware로 구현함)
	//rw.Header().Add("Content-Type", "application/json")
	
	//turn data into json (harder way)
	//b, err := json.Marshal(data)
	//utils.HandleErr(err)
	//fmt.Fprintf(rw, "%s",b)
	// easier way
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request){
	switch r.Method {
	case "GET":
		//web에 contents가 json이라고 알려줌 (Middleware로 구현함)
		//rw.Header().Add("Content-Type", "application/json")
		//Allblocks를 Json으로 encode 함
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var addBlockBody addBlockBody
		//r.body에서 POST로 Json data를 decode함, 여기선 Message를 가져와서 addBlockBody에 넣음
		//이 pointer로 써야 하고, err handler 있어야 함
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		//Blockchain package에 있는 block data 넣는 함수 호출
		blockchain.Blockchain().AddBlock(addBlockBody.Message)
		//POST일 때에는 statusCreated (==201) 넣음
		rw.WriteHeader(http.StatusCreated)
	}	
}

func block(rw http.ResponseWriter, r *http.Request){
	//gorillaMux - Vars는 http.Request에서 변수를 map으로 가져옴 
	//여기서는 map[id:number] - /{id:[0-9]} 
	vars := mux.Vars(r)
	//해당 id number 변수로 가져옴
	//height := vars["height"]

	//strconv: string 변환 library, Atoi는 string to ineger	
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

//middleware
func jsonContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw,r)
	})
}

func Start (aPort int) {
	//HandleFunc이 동시에 다루어지기 때문에, 별도의 Mux (Multiplexer)를 생성해서, default Mux 대신 쓰이게 함 
	//handler := http.NewServeMux()
	
	port = fmt.Sprintf(":%d", aPort)

	//router , dispatcher 다루기 위해서 gorillamux 설치함
	//go get -u github.com/gorilla/mux
	router := mux.NewRouter()
	
	//middleware
	router.Use(jsonContentMiddleware)

	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET","POST")
	
	//gorillaMux - id:number 인 주소 처리 
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}