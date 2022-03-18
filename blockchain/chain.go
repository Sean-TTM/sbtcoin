package blockchain

import (
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

func (b *blockchain) AddBlock() {
	block := creatBlock(b.NewestHash, b.Height +1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func Blocks(b *blockchain) []*Block {
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

func persistBlockchain(b *blockchain){
	db.SaveCheckpoint(utils.ToBytes(b))
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

// 특정 address에서 사용되지 않은 TX만 취해서 리턴
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	//map - Key-value data base - make function 으로 map 방식으로 선언
	creatorTxs := make(map[string]bool)
	//input 값으로 받은 address 에서 input 으로 사용 된 tx들 찾아서 true로 마킹
	for _, block := range Blocks(b){
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxId] = true
				}
			}
			for index, output := range tx.TxOuts{
				if output.Owner == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut){
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
//transactions 중에 특정 address 만 찾아서 취합
	// var ownedTxOuts []*TxOut
	// txOuts := b.txOuts()
	// for _, tx := range txOuts {
	// 	if tx.Owner == address {
	// 		ownedTxOuts = append(ownedTxOuts, tx)
	// 	}
	// }
	// return ownedTxOuts
}

// //Blockchain 내 모든 transaction return
// func (b *blockchain) txOuts() []*TxOut {
// 	var txOuts []*TxOut
// 	blocks := b.Blocks()
// 	for _, block := range blocks {
// 		for _, tx := range block.Transactions {
// 			txOuts = append(txOuts, tx.TxOuts...)	
// 		}
// 	}
// 	return txOuts
// }

//adjust difficulty of mining by mining 걸린 시간 조사
func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
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

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height % difficultyInterval == 0 {
		//recalculate the difficulty
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func Blockchain() *blockchain {
		//do it only once
		once.Do(func() {
			//initialize blockchain
			b = &blockchain{
				Height : 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				// restore b from bytes
				b.restore(checkpoint)
			}
		})
	return b
}