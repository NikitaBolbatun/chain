package bulba_chain

import (
	"context"
	"crypto"
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

// first block with blockchain settings
type Genesis struct {
	//Account -> funds
	Alloc map[string]uint64
	//list of validators public keys
	Validators []crypto.PublicKey
}

func (g Genesis) ToBlock() Block {
	//todo impliment me
	// алфавитный порядок порядок genesis.Alloc
	return Block{}
}

type Message struct {
	From string
	Data interface{}
}

type NodeInfoResp struct {
	NodeName string
	BlockNum uint64
}

type connectedPeer struct {
	Address string
	In      chan Message
	Out     chan Message
	cancel  context.CancelFunc
}

func (cp connectedPeer) Send(ctx context.Context, m Message) {
	//todo timeout using context + done check
	cp.Out <- m
}
