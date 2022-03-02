package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data 	 string
	Hash	 string
	PrevHash string
}

type blockchain struct {
	//slice of pointers of block (block에 data 모두 불러올 필요 없기 때문)
	blocks []*Block
}

//singleton Pattern https://refactoring.guru/design-patterns/singleton/go/example#example-0
var b *blockchain
//여러군데에서 동시에 호출에도 한번만 실행할 수 있게 해주는 library
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks -1 ].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string){
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain()*blockchain {
	if b == nil {
		//do it only once
		once.Do(func() {
			//initialize blockchain
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}