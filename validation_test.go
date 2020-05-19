package bulba_chain

import (
	"crypto"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"reflect"
	"testing"
	"time"
)

func Test_validation(t *testing.T) {
	numOfPeers := 5
	initialBalance := uint64(10000)
	peers := make([]*Node, numOfPeers)

	genesis := Genesis{
		Alloc:      make(map[string]uint64),
		Validators: make([]crypto.PublicKey, 0, numOfPeers),
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
		genesis.Validators = append(genesis.Validators, key.Public())
	}

	for i := 0; i < numOfPeers; i++ {
		var err error
		peers[i], err = NewNode(keys[i], genesis)
		peers[i].validTurn = i+1
		if err != nil {
			t.Error(err)
		}
		peers[i].AddBlocksNode(genesis.ToBlock())
		for _, validator := range genesis.Validators {
			peers[i].validators = append(peers[i].validators, validator.(ed25519.PublicKey))
		}

	}
	for i := 0; i < len(peers); i++ {
		for j := i + 1; j < len(peers); j++ {
			err := peers[i].AddPeer(peers[j])
			if err != nil {
				t.Error(err)
			}
		}
	}
	for i := 0; i < 2; i++ {
		for _, peer := range peers {
			err := peer.Validation()
			if err != nil {
				t.Error(err)
			}
			time.Sleep(time.Second)
		}

	}
	for i := 0; i < numOfPeers; i++ {
		fmt.Println(peers[i].lastBlockNum)
	}

	for i := 0; i < numOfPeers; i++ {
		for j := i+1; i < numOfPeers; i++ {
			if !reflect.DeepEqual(peers[i].chain.block, peers[j].chain.block){
				t.Fatal()
			}
		}
	}
}
