package blockchain

import (
	"fmt"
	"sync"

	"github.com/sean-ttm/sbtcoin/db"
	"github.com/sean-ttm/sbtcoin/utils"
)

type blockchain struct {
	//slice of pointers of block (block에 data 모두 불러올 필요 없기 때문)
	//blocks []*Block
	NewestHash 	string	 `jason:"newestHash"`
	Height 		int		 `jason:"height"`
}

//singleton Pattern https://refactoring.guru/design-patterns/singleton/go/example#example-0
var b *blockchain
//여러군데에서 동시에 호출에도 한번만 실행할 수 있게 해주는 library
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist(){
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := creatBlock(data, b.NewestHash, b.Height +1 )
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func Blockchain()*blockchain {
	if b == nil {
		//do it only once
		once.Do(func() {
			//initialize blockchain
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				// restore b from bytes
				b.restore(checkpoint)
			}
		})
	}
	fmt.Println(b.NewestHash)
	return b
}