package bulba_chain

import (
	"golang.org/x/crypto/ed25519"
	"reflect"
	"testing"
)

func TestNode_AddTransaction(t *testing.T) {
	pubkey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	nd := &Node{
		validators: []ed25519.PublicKey{pubkey},
	}
	nd.state = map[string]uint64{
		"one":         200,
		"two":         50,
	}
	nd.transactionPool = map[string]Transaction{}
	transactions := NewTransaction("one", "two", 10, 5, nil, nil)
	err = nd.AddTransaction(transactions)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_SingTransaction(t *testing.T) {
	numOfPeers := 2
	initialBalance := uint64(10000)
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
	}
	tr := Transaction{
		From:   peers[1].NodeAddress(),
		To:     peers[0].NodeAddress(),
		Amount: 100,
		Fee:    10,
		PubKey: keys[0].Public().(ed25519.PublicKey),
	}
	//Проверка что до подписи она была пустая
	var a[]byte
	if !reflect.DeepEqual(tr.Signature,a){
		t.Fatal()
	}
	tr, err = peers[1].SignTransaction(tr)
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(tr.Signature,a){
		t.Fatal()
	}
}
