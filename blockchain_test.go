package bulba_chain

import (
	"golang.org/x/crypto/ed25519"
	"testing"
)




func TestNode_Sync(t *testing.T) {
	numOfPeers := 5
	initialBalance := uint64(100000)
	peers := make([]*Node, numOfPeers)

	genesis := Genesis{
		Alloc: make(map[string]uint64),
	}

	keys := make([]ed25519.PrivateKey, numOfPeers)
	for i := range keys {
		_, key, err := ed25519.GenerateKey(nil)
		if err != nil {
			t.Fatal(err)
		}
		keys[i] = key
		address, err := PubKeyToAddress(key.Public())
		if err != nil {
			t.Error(err)
		}
		genesis.Alloc[address] = initialBalance
	}
	var err error
	for i := 0; i < numOfPeers; i++ {
		peers[i], err = NewNode(keys[i], genesis)
		if err != nil {
			t.Error(err)
		}
		err := peers[i].AddBlock(genesis.ToBlock())
		if err != nil {
			t.Fatal(err)
		}
	}

	if err != nil {
		t.Error(err)
	}

}