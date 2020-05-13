package bulba_chain

import (
	"crypto"
	"golang.org/x/crypto/ed25519"
	"testing"
	"time"
)

func Test_validation(t *testing.T) {
	numOfPeers := 5
	initialBalance := uint64(10000)
	peers := make([]*Node, numOfPeers)

	genesis := Genesis{
		Alloc:      make(map[string]uint64),
		Validators: make([]crypto.PublicKey, numOfPeers),
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
		genesis.Validators = append(genesis.Validators, address)
	}
	for _, peer := range peers{
		err:= peer.Validation()
		if err!=nil {
			t.Error(err)
		}
	}
	//todo: тесты должны быть автоматизированные. желательно без sleep (как и код)
	time.Sleep(time.Second*5)
	t.Log(peers)
}