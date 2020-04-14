package chain

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
	panic("implement me")
}

func (c *Node) AddBlock(b Block) error {
	for _, v := range b.Transactions {
		c.state[v.From] = c.state[v.From] - v.Amount - v.Fee
		c.state[v.To] = c.state[v.To] + v.Amount
		validatorKey := c.validators[int(b.BlockNum%uint64(len(c.validators)))]
		validatorAddr, err := PubKeyToAddress(validatorKey)
		if err != nil {
			return err
		}
		c.state[validatorAddr] += v.Fee
	}
	c.blocks = append(c.blocks, b)
	return nil
}
