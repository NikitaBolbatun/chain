package bulba_chain

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)
type CommitPool struct {
	Commit Commit
}
type Commit struct {
	RevealNumber []uint64
	Number []uint64
	shift []uint64

	commitMutex sync.RWMutex
}
func (c *Node) ValidationCommitReveal(block Block) error {
	go func() {
		for {
			switch {
			case c.ValidTurn():
				err := c.Insertblock(block)
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




func (c *Node) AddCommit (n uint64, nodeName string) error{
	c.commit.Commit.commitMutex.Lock()
	defer c.commit.Commit.commitMutex.Unlock()
	c.commit.Commit.Push(n)
	c.commit.Commit.Number = append(c.commit.Commit.Number,Encrypt(n,5))
	lastCommit := c.commit.Commit.Pop()
	if lastCommit == Encrypt(n,5) && nodeName != c.address {
		return nil
	}
	ctx := context.Background()
	c.Broadcast(ctx, Message{
		From: c.address,
		Data: CommitSend{
			NodeName:    c.address,
			commit: n,
		},
	})
	return nil
}
func (c *Node) Reveal (reveal Commit) (string,error) {
	var str string
	for i:=0; i < len(reveal.Number); i++ {
		reveal.Number[i] = Decryptor(reveal.Number[i],5)
		str += strconv.Itoa(int(reveal.Number[i]))
	}
	hashStr, err := c.Hash(str)
	if err!=nil{
		return "" , err
	}
	return hashStr,nil
}
func (c *Commit) Push(n uint64) {
	c.Number = append(c.Number, n)
}
func (c *Commit) PushShift(n uint64) {
	c.shift = append(c.Number, n)
}

func (c *Commit) Pop() uint64 {
	n := len(c.Number) - 1
	result := c.Number[n]
	c.Number = c.Number[:n]
	return result
}

func Encrypt(number uint64, shift uint64) uint64 {
	encrypt :=number + shift
	return encrypt
}

func Decryptor(number uint64, shift uint64) uint64 {
	decryptor :=number  - shift
	return decryptor
}
func (c *Node) Hash(number string) (string, error) {
	b, err := Bytes(number)
	if err != nil {
		return "", err
	}
	return Hash(b)
}
