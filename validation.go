package bulba_chain

import (
	"fmt"
	"time"
)

func (c *Node) ValidTurn() bool {
	validator,err := c.GetValidator(c.chain.lastNum)
	if err!= nil {
		return false
	}
	if validator != c.address {
		return false
	}
	return true
}

func (c *Node) Validation() error {
	go func() {
		for {
			switch {
			case c.ValidTurn():
				block, err:=c.AddBlock()
				if err!= nil {
					return
				}
				err = c.Insertblock(block)
				if err!=nil {
					return
				}
				fmt.Println(c.address, "add block", c.lastBlockNum)
				return
			default:

			}
		}
		time.Sleep(time.Second * 1)
	}()

	return nil
}

