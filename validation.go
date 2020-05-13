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
	// todo: зачем тут горутина. ты же не доживаешься валидации и всегда возвращаешь nil ошибку
	go func() {
		for {
			// todo: не понимаю, что это за конструкция с switch - case true и зачем она нужна
			switch {
			case true:
				// todo: Insertblock внутри себя заполняет c.blocks, но при вызове Insertblock мы читаем c.blocks
				// todo: итого c.blocks мы читаем до того, как отработает Insertblock. c.blocks всегда nil
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