package bulba_chain

import (
	"context"
	"fmt"
	"time"
)

func (c *Node) ValidTurn() bool {
	validator,err := c.GetValidator(c.chain.lastNum+1)
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
				c.Broadcast(context.TODO(), Message{
					From: c.address,
					Data: BlockSend{
						NodeName:    c.address,
						Block: 		 block,
					},
				})
				time.Sleep(time.Second*4)
				return

			default:

			}
		}

	}()

	return nil
}

