package bulba_chain

import (
	"golang.org/x/crypto/ed25519"
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
	nd.transactionPool = map[string]Transaction{}
	transactions := NewTransaction("one", "two", 10, 5, pubkey, nil)
	err = nd.AddTransaction(transactions)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nd.transactionPool)
}
