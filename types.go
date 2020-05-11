package bulba_chain

import (
	"context"
	"errors"
)

func (bl *Block) Hash() (string, error) {
	if bl == nil {
		return "", errors.New("empty block")
	}
	b, err := Bytes(bl)
	if err != nil {
		return "", err
	}
	return Hash(b)
}

type Message struct {
	From string
	Data interface{}
}

type NodeInfoResp struct {
	NodeName string
	BlockNum uint64
}
type BlockInfoResp struct {
	NodeName string
	BlockNum uint64
}

type connectedPeer struct {
	Address string
	In      chan Message
	Out     chan Message
	cancel  context.CancelFunc
}

type BlockHandshake struct {
	NodeName  string
	BlockNum  uint64
	BlockTo   uint64
	lastBlock uint64
	Block     Block
}
func (cp connectedPeer) Send(ctx context.Context, m Message) {
	//todo timeout using context + done check
	cp.Out <- m
}
