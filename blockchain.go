package bulba_chain

import (
	"context"
	"crypto"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"log"
	"sync"
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
	//handshake state
	handshake  bool
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

	blockMutex sync.RWMutex
	transMutex sync.RWMutex
	peersMutex sync.RWMutex

}

func (c *Node) NodeKey() crypto.PublicKey {
	return c.key.Public()
}

func (c *Node) Connection(address string, in chan Message, out chan Message) chan Message {
	if out == nil {
		out = make(chan Message, MSGBusLen)
	}
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
			//c.Broadcast(ctx, msg)
		}
	}
}

func (c *Node) processMessage(address string, message Message,ctx context.Context) error {
	switch m := message.Data.(type) {
	//case NodeInfoResp:
	//	return c.NodeInfoResp(address,m,ctx)
	//case BlockHandshake:
	//	return c.BlocksHandshake(address, m, ctx)
	case BlockSend:
		return c.BlockMessage(m,ctx)
	case TransactionSend:
		return c.TransactionMessage(m)
	}
	return nil
}


//func (c *Node) NodeInfoResp(address string, message NodeInfoResp, ctx context.Context) error {
//	c.blockMutex.Lock()
//	defer c.blockMutex.Unlock()
//
//	if c.lastBlockNum < message.BlockNum && !c.handshake{
//		c.handshake = true
//		c.peers[address].Send(ctx, Message{
//			From: c.address,
//			Data: BlockHandshake{
//				NodeName: c.address,
//				BlockNum: c.lastBlockNum,
//			},
//		})
//	}
//	return nil
//}
//
//
//func (c *Node) BlocksHandshake(address string, message BlockHandshake, ctx context.Context) error {
//	c.blockMutex.Lock()
//	defer c.blockMutex.Unlock()
//	fmt.Println(c.address, "connected to ", address, "need sync", c.lastBlockNum < message.BlockNum)
//	if message.BlockNum<c.lastBlockNum {
//		c.blocks = append(c.blocks, message.Block)
//		c.lastBlockNum++
//		if c.lastBlockNum < message.LastBlockName {
//			c.peers[address].Send(ctx, Message{
//				From: c.address,
//				Data: BlockHandshake{
//					NodeName: c.address,
//					BlockNum: c.lastBlockNum,
//				},
//			})
//		}
//	} else {
//		c.handshake = false
//		c.peers[address].Send(ctx, Message{
//			From: c.address,
//			Data: BlockHandshake{
//				NodeName:      c.address,
//				BlockNum:      message.BlockNum,
//				LastBlockName: c.lastBlockNum,
//				Block:         c.blocks[message.BlockNum],
//			},
//		})
//	}
//	return nil
//}


func (c *Node) BlockMessage(message BlockSend, ctx context.Context) error {
	c.blocks = append(c.blocks,message.Block)
	c.lastBlockNum++
	c.Broadcast(ctx, Message{
		From: c.address,
		//Data: ,
	})

	return nil
}

func (c *Node) TransactionMessage(message TransactionSend) error {
	fmt.Println("Transaction Message ", message.NodeName)
	err := c.addtransaction(message.Transaction)
	if err != nil {
		return err
	}
	return nil
}
func (c *Node) addtransaction(transaction Transaction)  error {
	hash, err := transaction.Hash()
	if err != nil {
		return err
	}
	c.transactionPool[hash] = transaction
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
	c.peersMutex.Lock()
	defer c.peersMutex.Unlock()
	for _, v := range c.peers {
		if v.Address != c.address && v.Address != msg.From{
			v.Send(ctx, msg)
		}
	}
}