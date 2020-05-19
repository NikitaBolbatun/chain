package bulba_chain

import (
	"crypto"
	"golang.org/x/crypto/ed25519"
	"reflect"
	"testing"
)

func TestGenesis_ToBlockSorting(t *testing.T) {
	genesis := NewGenesis()
	genesis.Alloc = map[string]uint64{
		"one":   100,
		"two":   110,
		"fair":  10,
		"qwert": 5,
		"abc":   100,
		"bcd":   545,
		"000":   545,
		"938":   545,
	}
	block := genesis.ToBlock()
	block2 := genesis.ToBlock()
	if !reflect.DeepEqual(block, block2) {
		t.Fatal()
	}

}

func TestGenesis_ToBlock(t *testing.T) {
	genesis := NewGenesis()
	genesis.Alloc = map[string]uint64{
		"one": 100,
		"two": 110,
	}
	pubkey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	validatorAddr, err := PubKeyToAddress(pubkey)
	if err != nil {
		t.Fatal(err)
	}

	genesis.Validators = []crypto.PublicKey{
		validatorAddr,
	}
	block := genesis.ToBlock()
	if reflect.DeepEqual(block, Block{}){
		t.Fatal()
	}
}
