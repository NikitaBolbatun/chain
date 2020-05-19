package bulba_chain

import (
	"golang.org/x/crypto/ed25519"
	"testing"
)

//func TestBlockProcessing(t *testing.T) {
//	pubkey, _, err := ed25519.GenerateKey(nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	validatorAddr, err := PubKeyToAddress(pubkey)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	nd := &Node{
//		validators: []ed25519.PublicKey{pubkey},
//	}
//	nd.state = map[string]uint64{
//		"one":         200,
//		"two":         50,
//		validatorAddr: 50,
//	}
//	nd.AddBlocksNode(nd.genesis.ToBlock())
//	block, errr :=nd.AddBlock()
//	if errr!=nil {
//		t.Fatal()
//	}
//	err = nd.Insertblock(block)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if nd.state["one"] != 90 {
//		t.Error()
//	}
//	if nd.state["two"] != 150 {
//		t.Error()
//	}
//
//}
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

func TestBlock_Verify(t *testing.T) {
	pubkey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	block := &Block{
		BlockNum:      1,
		BlockHash:     "0",
		PrevBlockHash: "0",
		StateHash:     "",
		Timestamp:    1000,
		Transactions: []Transaction{},
		Signature:    nil,
	}
	block.BlockHash, err = block.Hash()
	if err != err {
		t.Fatal(err)
	}
	err = block.SignBlock(privateKey)
	if err != err {
		t.Fatal(err)
	}

	ok := block.VerifyBlockSign(pubkey)
	if err != nil {
		t.Fatal()
	}
	if ok {
		t.Fatal()
	}

}
