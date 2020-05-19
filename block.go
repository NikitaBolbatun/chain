package bulba_chain

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"reflect"
	"sync"
	"time"
)


type ChainBlocks struct {
	block[] Block
	lastNum uint64
	blockMutex sync.RWMutex
}
func NewBlocks() *ChainBlocks {
	return &ChainBlocks{
		block:       make([]Block, 0),
		lastNum: 0,
	}
}
type CommitBlock struct {
	Number uint64
	shift uint64
}
type Block struct {
	BlockNum      uint64
	Timestamp     int64
	Transactions  []Transaction
	Reveal		  CommitBlock
	BlockHash     string `json:"-"`
	PrevBlockHash string
	StateHash     string
	Signature     []byte `json:"-"`


}


func NewBlock(num uint64, transaction []Transaction, prevBlockHash string) Block {
	return Block{
		BlockNum:      num,
		Timestamp:     time.Now().Unix(),
		Transactions:  transaction,
		PrevBlockHash: prevBlockHash,
	}
}

func (c *Node) GetBlockByNumber(ID uint64) Block {
	return c.chain.block[ID]
}

func (c *Node) AddBlock() (Block,error) {
	c.blockMutex.Lock()
	defer c.blockMutex.Unlock()
	prevBlock := c.chain.block[len(c.chain.block)-1]
	block := NewBlock(prevBlock.BlockNum+1, nil, prevBlock.BlockHash )
	sing, err := Bytes(block)
	if err!=nil{
		return Block{}, err
	}
	block.Signature=ed25519.Sign(c.key,sing)
	blockHash,err := block.Hash()
	if err!=nil{
		return Block{}, err
	}

	block.BlockHash = blockHash
	return block, nil
}

func (c *Node) Insertblock(b Block) error {
	if !reflect.DeepEqual(b, Block{}){
		if reflect.DeepEqual(c.chain.block[c.lastBlockNum-1], b) {
			return errors.New("block exists")
		}
	}
	validatorAddr, err := c.GetValidator(b.BlockNum)
	if err!=nil {
		return err
	}
	if reflect.DeepEqual(b.Reveal,0) {
		for _, v := range b.Transactions {
			c.state[v.From] = c.state[v.From] - v.Amount - v.Fee
			c.state[v.To] = c.state[v.To] + v.Amount
			c.state[validatorAddr] += v.Fee
		}
	} else {
		c.commit.Commit.Push(b.Reveal.Number)
		if reflect.DeepEqual(b.Reveal.shift,0) {
			c.commit.Commit.PushShift(b.Reveal.shift)
		}

	}
	ok := b.VerifyBlockSign(c.validators[int(c.lastBlockNum%uint64(len(c.validators)))])
	if ok {
		return errors.New("not correct verify block")
	}
	fmt.Println("add block num",c.chain.lastNum,"node",c.address)
	c.AddBlocksNode(b)
	err = c.chain.AddBlockChain(b)
	if err!=nil{
		return err
	}
	return nil
}
func (c *Node) AddBlocksNode (b Block) {
	c.blockMutex.RLock()
	defer c.blockMutex.RUnlock()
	c.chain.block = append(c.chain.block, b)
	c.lastBlockNum++
}
func (c *Node) GetValidator(n uint64) (string,error) {
	validatorKey , err:= PubKeyToAddress(c.validators[int(n%uint64(len(c.validators)))])
	if err!=nil {
		return " ", err
	}
	return validatorKey, nil
}
//Эти функции над использывать
func (bl Block) SignBlock(key ed25519.PrivateKey) error {
	b, err := Bytes(bl.BlockHash)
	if err != nil {
		return err
	}
	bl.Signature = ed25519.Sign(key, b)
	return  nil
}
func (bl *Block) VerifyBlockSign(key ed25519.PublicKey) bool {
	b, err := Bytes(bl.BlockHash)
	if err != nil {
		return false
	}

	return ed25519.Verify(key, b, bl.Signature)
}

func (ChBlock *ChainBlocks) GetLastBlock() uint64 {
	ChBlock.blockMutex.RLock()
	defer ChBlock.blockMutex.RUnlock()
	return ChBlock.lastNum
}
func (ChBlock *ChainBlocks) AddBlockChain(block Block) error {
	ChBlock.blockMutex.Lock()
	defer ChBlock.blockMutex.Unlock()
	ChBlock.block = append(ChBlock.block, block)
	ChBlock.lastNum = block.BlockNum
	return nil
}