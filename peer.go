package bulba_chain

import (
	"errors"
)

func (c *Node) AddPeer(peer Blockchain) error {
	remoteAddress, err := PubKeyToAddress(peer.NodeKey())
	if err != nil {
		return err
	}

	if c.address == remoteAddress {
		return errors.New("self connection")
	}

	if _, ok := c.peers[remoteAddress]; ok {
		return nil
	}

	out := make(chan Message, MSGBusLen)
	in := peer.Connection(c.address, out,nil)
	c.Connection(remoteAddress, in, out)
	return nil
}

func (c *Node) RemovePeer(peer Blockchain) error {
	delete(c.peers, peer.NodeAddress())
	return nil
}

func (c *Node) GetBalance(account string) (uint64, error) {
	balance, err := c.state[account]
	if !err {
		return 0, errors.New("unknown user")
	}

	return balance, nil
}
