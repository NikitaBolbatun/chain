package bulba_chain

import (
	"crypto"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestCommit_Send (t *testing.T) {
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
		for j := i + 1; j < 5; j++ {
			err := peers[i].AddPeer(peers[j])
			if err != nil {
				t.Error(err)
			}
		}
	}


	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		number := rand.Uint64()
		shiftCommit := Encrypt(number,5)
		err = peers[i].AddCommit(shiftCommit ,peers[i].address)
		if err!= nil{
			t.Fatal()
		}
	}

	time.Sleep(time.Second* 5)

	for i := 0; i < len(peers); i++ {
		fmt.Println(peers[i].commit.Number)
	}
	for i := 0; i < len(peers); i++ {
		for j := i +1 ; j < len(peers); j++ {
			if !reflect.DeepEqual(peers[i].commit, peers[j].commit) {
				t.Fatal()
			}
		}
	}

}

func TestNode_Reveal(t *testing.T) {
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
		for j := i + 1; j < 5; j++ {
			err := peers[i].AddPeer(peers[j])
			if err != nil {
				t.Error(err)
			}
		}
	}


	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		number := rand.Uint64()/rand.Uint64()
		shiftCommit := Encrypt(number,5)
		err = peers[i].AddCommit(shiftCommit ,peers[i].address)
		if err!= nil{
			t.Fatal()
		}
	}

	time.Sleep(time.Second* 5)

	for i := 0; i < len(peers); i++ {
		fmt.Println(peers[i].commit.Number)
	}

	for i := 0; i < len(peers); i++ {
		peers[i].random, err =peers[i].Reveal(peers[i].commit)
		if err!= nil{
			t.Fatal()
		}
		fmt.Println(peers[i].random)
	}
	for i := 0; i < len(peers); i++ {
		for j := i +1 ; j < len(peers); j++ {
			if !reflect.DeepEqual(peers[i].random, peers[j].random) {
				t.Fatal()
			}
		}
	}

}

func Test_FullCommit_Reveal(t *testing.T) {

}