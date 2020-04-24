package bulba_chain

import (
	"context"
	"crypto"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"log"
	"reflect"
)

const MSGBusLen = 100

func NewNode(key ed25519.PrivateKey, genesis Genesis) (*Node, error) {
	address, err := PubKeyToAddress(key.Public())
	if err != nil {
		return nil, err
	}
	return &Node{
		key:             key,
		address:         address,
		genesis:         genesis,
		blocks:          make([]Block, 0),
		lastBlockNum:    0,
		peers:           make(map[string]connectedPeer, 0),
		state:           make(map[string]uint64),
		transactionPool: make(map[string]Transaction),
	}, err
}

type Node struct {
	key          ed25519.PrivateKey
	address      string
	genesis      Genesis
	lastBlockNum uint64

	//state
	blocks []Block
	//peer address - > peer info
	peers map[string]connectedPeer
	//hash(state) - хеш от упорядоченного слайса ключ-значение
	//todo hash()
	state      map[string]uint64
	validators []ed25519.PublicKey

	//transaction hash - > transaction
	transactionPool map[string]Transaction
}

func (c *Node) NodeKey() crypto.PublicKey {
	return c.key.Public()
}

func (c *Node) Connection(address string, in chan Message) chan Message {
	out := make(chan Message, MSGBusLen)
	ctx, cancel := context.WithCancel(context.Background())
	c.peers[address] = connectedPeer{
		Address: address,
		Out:     out,
		In:      in,
		cancel:  cancel,
	}

	go c.peerLoop(ctx, c.peers[address])
	return c.peers[address].Out
}

func (c *Node) peerLoop(ctx context.Context, peer connectedPeer) {
	//todo handshake
	peer.Send(ctx, Message{
		From: c.address,
		Data: NodeInfoResp{
			NodeName: c.address,
			BlockNum: c.lastBlockNum,
		},
	})
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-peer.In:
			err := c.processMessage(peer.Address, msg,ctx)
			if err != nil {
				log.Println("Process peer error", err)
				continue
			}

			//broadcast to connected peers
			c.Broadcast(ctx, msg)
		}
	}
}
func (c *Node) processMessage(address string, msg Message,ctx context.Context) error {
	switch m := msg.Data.(type) {
	case NodeInfoResp:
		fmt.Println(c.address, "connected to ", address, "need sync", c.lastBlockNum < m.BlockNum)
		if c.lastBlockNum < m.BlockNum{
			c.peers[address].Send(ctx, Message{
				From: c.address,
				Data: BlockInfoResp{
					NodeName: c.address,
					BlockNum: c.lastBlockNum + 1,
				},
			})
		}

	case BlockByNumResp:
		if !reflect.DeepEqual(m.Block, Block{}) {
			err := c.AddBlock(m.Block)
			if err != nil {
				return err
			}
			if c.lastBlockNum < m.LastBlockNum {
				c.peers[address].Send(ctx, Message{
					From: c.address,
					Data: BlockByNumResp{
						NodeName: c.address,
						BlockNum: c.lastBlockNum + 1,
					},
				})
			}
		} else {
			c.peers[address].Send(ctx, Message{
				From: c.address,
				Data: BlockByNumResp{
					NodeName:     c.address,
					BlockNum:     m.BlockNum,
					LastBlockNum: c.lastBlockNum,
					Block:        c.GetBlockByNumber(m.BlockNum),
				},
			})
		}
	}
	return nil
}


func (c *Node) NodeInfo() NodeInfoResp {
	return NodeInfoResp{
		NodeName: c.address,
		BlockNum: c.lastBlockNum,
	}
}

func (c *Node) NodeAddress() string {
	return c.address
}


func (c *Node) Broadcast(ctx context.Context, msg Message) {
	for _, v := range c.peers {
		if v.Address != c.address {
			v.Send(ctx, msg)
		}
	}
}