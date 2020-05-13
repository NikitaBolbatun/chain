package bulba_chain

import (
	"context"
	"fmt"
	"time"
)

func (c *Node) ValidTurn() bool {
	return true
}

func (c *Node) Validation() error {
	go func() {
		for {
			switch {
			case true:
				block := c.Insertblock(*NewBlock(c.lastBlockNum+1, nil, c.blocks[c.lastBlockNum].BlockHash))
				fmt.Println(c.address, "add block", c.lastBlockNum)
				ctx := context.Background()
				c.Broadcast(ctx, Message{
					From: c.address,
					Data: block,
				})
				return
			default:

			}
		}
		time.Sleep(time.Second * 1)
	}()

	return nil
}