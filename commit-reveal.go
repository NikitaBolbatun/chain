package bulba_chain

import (
	"context"
	"strconv"
	"sync"
	"time"
)

type Commit struct {
	Number []uint64
	shift uint64

	commitMutex sync.RWMutex
}
func (c *Node) AddCommit (n uint64, nodeName string) error{
	c.commit.commitMutex.Lock()
	defer c.commit.commitMutex.Unlock()
	c.commit.Push(n)
	c.commit.Number = append(c.commit.Number,Encrypt(n,5))
	lastCommit := c.commit.Pop()
	if lastCommit == Encrypt(n,5) && nodeName != c.address {
		return nil
	}
	time.Sleep(time.Second)
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
