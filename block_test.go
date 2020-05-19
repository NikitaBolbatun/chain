package bulba_chain

import (
	"golang.org/x/crypto/ed25519"
	"testing"
)

func TestBlockProcessing(t *testing.T) {
	pubkey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	validatorAddr, err := PubKeyToAddress(pubkey)
	if err != nil {
		t.Fatal(err)
	}

	nd := &Node{
		validators: []ed25519.PublicKey{pubkey},
	}
	nd.state = map[string]uint64{
		"one":         200,
		"two":         50,
		validatorAddr: 50,
	}

	err = nd.Insertblock(Block{
		BlockNum: 1,
		Transactions: []Transaction{
			{
				From:   "one",
				To:     "two",
				Fee:    10,
				Amount: 100,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if nd.state["one"] != 90 {
		t.Error()
	}
	if nd.state["two"] != 150 {
		t.Error()
	}

}
func TestNode_GetValidator(t *testing.T) {
	pubkey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	validatorAddr, err := PubKeyToAddress(pubkey)
	if err != nil {
		t.Fatal(err)
	}

	nd := &Node{
		validators: []ed25519.PublicKey{pubkey},
	}
	nd.state = map[string]uint64{
		"one":         200,
		"two":         50,
		validatorAddr: 50,
	}
	_, err = nd.GetValidator(1)
	if err != nil {
		t.Fatal(err)
	}
}
