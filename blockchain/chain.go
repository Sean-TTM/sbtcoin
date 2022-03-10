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
	CurrentDifficulty int `json:"currentDifficulty"`
}

const (
	defaultDifficulty   int = 2
	difficultyInterval  int = 5
	blockInterval		int = 2
	allowedRange		int = 2
)

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
	b.CurrentDifficulty = block.Difficulty
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

//adjust difficulty of mining by mining 걸린 시간 조사
func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval - 1]
	actualTime := (newestBlock.Timestamp - lastRecalculatedBlock.Timestamp)/60 
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - -1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height % difficultyInterval == 0 {
		//recalculate the difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func Blockchain()*blockchain {
	if b == nil {
		//do it only once
		once.Do(func() {
			//initialize blockchain
			b = &blockchain{
				Height : 0,
			}
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