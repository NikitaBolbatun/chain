package bulba_chain

import (
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

func (c *Node) AddBlock(b Block) error {
	for _, v := range b.Transactions {
		c.state[v.From] = c.state[v.From] - v.Amount - v.Fee
		c.state[v.To] = c.state[v.To] + v.Amount
		validatorAddr, err := c.GetValidator(b.BlockNum)
		if err != nil {
			return err
		}
		c.state[validatorAddr] += v.Fee
	}
	c.blocks = append(c.blocks, b)
	c.lastBlockNum++
	return nil
}

func (c *Node) GetValidator(n uint64) (string, error) {
	validatorKey := c.validators[int(n%uint64(len(c.validators)))]
	return PubKeyToAddress(validatorKey)
}
