package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sean-ttm/sbtcoin/db"
	"github.com/sean-ttm/sbtcoin/utils"
)

type Block struct {
	Data 	 	string  `json:"data"`
	Hash	 	string  `json:"hash"`
	PrevHash 	string  `json:"prevHash,omitempty"`
	Height	 	int 	`json:"height"`
	//Added for Proof of Work
	Difficulty 	int		`json:"difficulty"`
	//Nonce is the only thing that can be chaned by miners
	Nonce		int		`json:"nonce"`
	//to find out how long it takes for creating block
	Timestamp 	int		`json:"timestamp"`
}

//db에 block 저장
func (b *Block) persist(){
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

//define not found error 
var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte){
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error){
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine(){
	//0이 difficulty 만큼 연속으로 있는 Hash 찾기 위함
	target := strings.Repeat("0", b.Difficulty)
	for {
		//Block string으로 바꿈
		//blockAsString := fmt.Sprint(b)
		//sha256 이용해서 Hash 생성
		//hash := fmt.Sprintf("%x", sha256.Sum256([]byte(blockAsString)))
		//fmt.Printf("Block as String:%s\nHash:%s\nNonce:%d\n\n\n", blockAsString, hash, b.Nonce )
		
		hash := utils.Hash(b)
		fmt.Printf("Target:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		//string.Hasprefix -> hash가 target="00"으로 시작하는지 return하는 함수
		if strings.HasPrefix(hash, target){
			//timestamp로 difficulty level check 하기 위함
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func creatBlock(data string, prevHash string, height int) *Block{
	block := &Block{
		Data: data,
		Hash: "",
		PrevHash: prevHash,
		Height: height,
		//const difficulty predefined
		Difficulty: Blockchain().difficulty(),
		Nonce: 0,
	}
	//mining function - proof of work
	block.mine()
	//save block data to DB
	block.persist()
	return block
}