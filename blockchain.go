package bulba_chain

import (
	"context"
	"crypto"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"log"
	"reflect"
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
		chain:          *NewBlocks(),
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
	chain ChainBlocks
	//peer address - > peer info
	peers map[string]connectedPeer
	validTurn int
	//hash(state) - хеш от упорядоченного слайса ключ-значение
	state      map[string]uint64
	validators []ed25519.PublicKey
	commit CommitPool
	//transaction hash - > transaction
	transactionPool map[string]Transaction
	random        string

	blockMutex sync.RWMutex
	transMutex sync.RWMutex
	peersMutex sync.RWMutex

}

func (c *Node) NodeKey() crypto.PublicKey {
	return c.key.Public()
}
func (c *Node) AddValidatorToNode() error {
	c.Insertblock(c.genesis.ToBlock())
	for _, validator := range c.genesis.Validators {
		c.validators = append(c.validators, validator.(ed25519.PublicKey))
	}
	return nil
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
			//c.Broadcast(ctx, msg)
		}
	}
}

func (c *Node) processMessage(address string, message Message,ctx context.Context) error {
	var err error
	switch message:= message.Data.(type) {
	case NodeInfoResp:
		err = c.NodeInfoResp(address,message,ctx)
	case BlockHandshake:
		err = c.BlocksHandshake(address, message,ctx)
	case BlockSend:
		err = c.BlockMessage(message)
	case TransactionSend:
		err = c.TransactionMessage(message)
	case CommitSend:
		err = c.CommitMessage(message)
	}
	if err!=nil {
		return err
	}
	return nil
}


func (c *Node) NodeInfoResp(address string, message NodeInfoResp, ctx context.Context) error {
	c.blockMutex.Lock()
	defer c.blockMutex.Unlock()
	if c.lastBlockNum < message.BlockNum && !c.handshake{
		c.handshake = true
		c.peers[address].Send(ctx, Message{
			From: c.address,
			Data: BlockHandshake{
				NodeName: c.address,
				BlockNum: c.lastBlockNum,
			},
		})
	}
	return nil
}


func (c *Node) BlocksHandshake(address string, message BlockHandshake, ctx context.Context) error {
	c.blockMutex.Lock()
	defer c.blockMutex.Unlock()
	fmt.Println(c.address, "connected to ", address, "need sync", c.lastBlockNum < message.BlockNum)
	fmt.Println(message.NumberBlock, message.LastBlockName, c.lastBlockNum, message.BlockNum)
	if !reflect.DeepEqual(message.Block, Block{}) {
		c.chain.block = append(c.chain.block)
		c.lastBlockNum++
		if c.lastBlockNum < message.LastBlockName {
			c.peers[address].Send(ctx, Message{
				From: c.address,
				Data: BlockHandshake{
					NodeName: c.address,
					BlockNum: c.lastBlockNum + 1,
				},
			})
		} else {
			c.handshake = false
		}
	} else {
		c.peers[address].Send(ctx, Message{
			From: c.address,
			Data: BlockHandshake{
				NodeName:      c.address,
				BlockNum:      message.BlockNum,
				LastBlockName: c.lastBlockNum,
				Block:         c.GetBlockByNumber(message.BlockNum),
			},
		})
	}
	return nil
}


func (c *Node) BlockMessage(message BlockSend) error {
	c.blockMutex.RLock()
	defer c.blockMutex.RUnlock()
	fmt.Println("Send block BlockMessage",message.Block.BlockNum)
	err := c.Insertblock(message.Block)
	if err!=nil {
				return err
			}

	return nil
}

func (c *Node) TransactionMessage(message TransactionSend) error {
	fmt.Println("Transaction Message ", message.NodeName)
	err := c.AddTransaction(message.Transaction)
	if err != nil {
		return err
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
	c.peersMutex.Lock()
	defer c.peersMutex.Unlock()
	for _, v := range c.peers {
		if v.Address != c.address{
			v.Send(ctx, msg)
		}
	}
}

func (c *Node) CommitMessage(message CommitSend) error {
	fmt.Println("Commit Message ", message.NodeName)
	err := c.AddCommit(message.commit,message.NodeName)
	if err != nil {
		return err
	}
	return nil
}