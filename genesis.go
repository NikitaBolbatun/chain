package bulba_chain

import (
	"crypto"
	"sort"
)

// first block with blockchain settings
type Genesis struct {
	//Account -> funds
	Alloc map[string]uint64
	//list of validators public keys
	Validators []crypto.PublicKey
}

func NewGenesis() Genesis {
	return Genesis{
		Alloc:      make(map[string]uint64),
		Validators: []crypto.PublicKey{},
	}
}

func (g Genesis) ToBlock() Block {
	var transaction []Transaction
	for account, amount := range g.Alloc {
		transaction = append(transaction, Transaction{
			From:      "",
			To:        account,
			Amount:    amount,
			Fee:       0,
			PubKey:    nil,
			Signature: nil,
		})
	}
	sort.Slice(transaction, func(i, j int) bool {
		return transaction[i].To < transaction[j].To
	})
	block := Block{
		BlockNum:      0,
		Timestamp:     0,
		Transactions:  transaction,
		BlockHash:     "",
		PrevBlockHash: "",
		Signature:     nil,
	}
	blockbyte, err := Bytes(block)
	if err != nil {
		return Block{}
	}
	block.BlockHash, err = Hash(blockbyte)
	if err != nil {
		return Block{}
	}
	block.PrevBlockHash, err = Hash([]byte("0"))
	if err != nil {
		return Block{}
	}

	return block
}
