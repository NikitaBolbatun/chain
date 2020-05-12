package bulba_chain

import (
	"golang.org/x/crypto/ed25519"
	"time"
)

type Block struct {
	BlockNum      uint64
	Timestamp     int64
	Transactions  []Transaction
	BlockHash     string `json:"-"`
	PrevBlockHash string
	StateHash     string
	Signature     []byte `json:"-"`
}

func NewBlock(num uint64, transaction []Transaction, prevBlockHash string) *Block {
	return &Block{
		BlockNum:      num,
		Timestamp:     time.Now().Unix(),
		Transactions:  transaction,
		PrevBlockHash: prevBlockHash,
	}
}

func (c *Node) GetBlockByNumber(ID uint64) Block {
	return c.blocks[ID]
}
//Исправить баг, при хенд, валидаторы не заполняются из за разницы типов
func (c *Node) Insertblock(b Block) error {
	//validator := c.validators[int(c.lastBlockNum)%len(c.validators)]
	//validatorAddr, err := PubKeyToAddress(validator)
	//if err != nil {
	//	return err
	//}
	for _, v := range b.Transactions {
		c.state[v.From] = c.state[v.From] - v.Amount - v.Fee
		c.state[v.To] = c.state[v.To] + v.Amount
	//	c.state[validatorAddr] += v.Fee
	}
	c.blocks = append(c.blocks, b)
	c.lastBlockNum++
	return nil
}

func (c *Node) GetValidator(n uint64) (string, error) {
	validatorKey := c.validators[int(n%uint64(len(c.validators)))]
	return PubKeyToAddress(validatorKey)
}

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
