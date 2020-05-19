package bulba_chain

import "C"
import (
	"context"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"sync"
	"time"
)


// todo Сделать еще одну функцию добавки блока для фикса бага рассылки\
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

type Block struct {
	BlockNum      uint64
	Timestamp     int64
	Transactions  []Transaction
	Reveal		  uint64
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
	// добавление транзакций в блок
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

//Исправить баг, при хенд, валидаторы не заполняются из за разницы типов
func (c *Node) Insertblock(b Block) error {


	//validatorAddr, err := c.GetValidator(c.chain.lastNum)
	//if err != nil {
	//	return err
	//}
	for _, v := range b.Transactions {
		c.state[v.From] = c.state[v.From] - v.Amount - v.Fee
		c.state[v.To] = c.state[v.To] + v.Amount
		//c.state[validatorAddr] += v.Fee
	}

	fmt.Println("Добавляем блок",b.BlockNum,"пользователю",c.address)
	c.AddBlocksNode(b)
	c.chain.AddBlockChain(b)
	fmt.Println(c.chain.lastNum)
	c.Broadcast(context.TODO(), Message{
		From: c.address,
		Data: BlockSend{
			NodeName:    c.address,
			Block: b,
		},
	})
	return nil
}
func (c *Node) AddBlocksNode (b Block) {
	c.blockMutex.RLock()
	defer c.blockMutex.RUnlock()
	c.chain.block = append(c.chain.block, b)
	c.lastBlockNum++
}
func (c *Node) GetValidator(n uint64) (string,error) {
	fmt.Println(len(c.validators))
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
func (bl *Block) VerifyBlockSign(key ed25519.PublicKey) (bool, error) {
	b, err := Bytes(bl.BlockHash)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(key, b, bl.Signature), nil
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