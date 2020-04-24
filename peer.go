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
	in := peer.Connection(c.address, out)
	c.Connection(remoteAddress, in)
	return nil
}



func (c *Node) RemovePeer(peer Blockchain) error {
	delete(c.peers, peer.NodeAddress())
	return nil
}


func (c *Node) GetBalance(account string) (uint64, error) {
	panic("implement me")
}