package bulba_chain

import (
	"crypto"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"reflect"
	"testing"
	"time"
)


//func Test_Handshake(t *testing.T) {
//	numOfPeers := 5
//	initialBalance := uint64(10000)
//	peers := make([]*Node, numOfPeers)
//
//	genesis := Genesis{
//		Alloc:      make(map[string]uint64),
//		Validators: make([]crypto.PublicKey, 0, numOfPeers),
//	}
//
//	keys := make([]ed25519.PrivateKey, numOfPeers)
//	for i := range keys {
//		_, key, err := ed25519.GenerateKey(nil)
//		if err != nil {
//			t.Fatal(err)
//		}
//		keys[i] = key
//		address, err := PubKeyToAddress(key.Public())
//		if err != nil {
//			t.Error(err)
//		}
//		genesis.Alloc[address] = initialBalance
//		genesis.Validators = append(genesis.Validators, key.Public())
//	}
//
//	for i := 0; i < numOfPeers; i++ {
//		var err error
//		peers[i], err = NewNode(keys[i], genesis)
//		if err != nil {
//			t.Error(err)
//		}
//		peers[i].AddBlocksNode(genesis.ToBlock())
//		for _, validator := range genesis.Validators {
//			peers[i].validators = append(peers[i].validators, validator.(ed25519.PublicKey))
//		}
//
//	}
//	for i := 0; i < len(peers); i++ {
//		block, err := peers[0].AddBlock()
//		if err != nil {
//			t.Fatal()
//		}
//		err = peers[0].Insertblock(block)
//		if err != nil {
//			t.Fatal()
//		}
//	}
//
//	for i := 0; i < len(peers); i++ {
//		for j := i + 1; j < len(peers); j++ {
//			err := peers[i].AddPeer(peers[j])
//			if err != nil {
//				t.Error(err)
//			}
//		}
//	}
//
//	time.Sleep(time.Second* 3)
//
//	for i:=0; i<len(peers); i++ {
//		fmt.Println(peers[i].lastBlockNum)
//	}
//	for i := 0; i < len(peers); i++ {
//		for j := i +1 ; j < len(peers); j++ {
//			if !reflect.DeepEqual(peers[i].chain, peers[j].chain) {
//				t.Fatal()
//			}
//		}
//	}
//}


func Test_SendTransaction (t *testing.T) {
	numOfPeers := 5
	initialBalance := uint64(10000)
	peers := make([]*Node, numOfPeers)

	genesis := Genesis{
		Alloc: make(map[string]uint64),
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
	var err error
	for i := 0; i < numOfPeers; i++ {
		peers[i], err = NewNode(keys[i], genesis)
		if err != nil {
			t.Error(err)
		}
		peers[i].Insertblock(genesis.ToBlock())
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
	for i:=0; i<len(peers); i++ {
		fmt.Println(peers[0].validators)
	}
	tr := Transaction{
		From:   peers[1].NodeAddress(),
		To:     peers[2].NodeAddress(),
		Amount: 10,
		Fee:    10,
		PubKey: keys[3].Public().(ed25519.PublicKey),
	}



	tr, err = peers[3].SignTransaction(tr)
	if err != nil {
		t.Error(err)
	}
	err = peers[3].AddTransaction(tr)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second* 5)
	fmt.Println(len(peers[0].transactionPool))
	for i := 0; i < len(peers); i++ {
		for j := i +1 ; j < len(peers); j++ {
			if !reflect.DeepEqual(peers[i].transactionPool, peers[j].transactionPool) {
				t.Fatal()
			}
		}
	}

}


func Test_SendBlock(t *testing.T) {
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
		peers[i].AddBlocksNode(genesis.ToBlock())
		if err != nil {
			t.Fatal(err)
		}
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
	block, err:= peers[0].AddBlock()
	if err != nil {
		t.Fatal()
	}
	err = peers[0].Insertblock(block)
	if err!=nil {
		t.Fatal()
	}
	time.Sleep(time.Second)
	for i := 0; i < len(peers); i++ {
		fmt.Println(peers[i].lastBlockNum)
	}
	for i := 0; i < len(peers); i++ {
		for j := i +1 ; j < len(peers); j++ {
			if !reflect.DeepEqual(peers[i].lastBlockNum, peers[j].lastBlockNum) {
				t.Fatal()
			}
		}
	}

}
