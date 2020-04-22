package bulba_chain

import (
	"bytes"
	"encoding/gob"
	"golang.org/x/crypto/ed25519"
)

type Transaction struct {
	From   string
	To     string
	Amount uint64
	Fee    uint64
	PubKey ed25519.PublicKey

	Signature []byte `json:"-"`
}

func (t Transaction) Hash() (string, error) {
	b, err := Bytes(t)
	if err != nil {
		return "", err
	}
	return Hash(b)
}

func (t Transaction) Bytes() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(t)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func NewTransaction(from, to string, amount, fee uint64, key ed25519.PublicKey, signature []byte) Transaction {
	return Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		PubKey:    key,
		Signature: signature,
	}
}

func (c *Node) AddTransaction(transaction Transaction) error {
	hash, err := transaction.Hash()
	if err != nil {
		return err
	}
	c.transactionPool[hash] = transaction
	return nil
}

func (c *Node) SignTransaction(transaction Transaction) (Transaction, error) {
	b, err := transaction.Bytes()
	if err != nil {
		return Transaction{}, err
	}

	transaction.Signature = ed25519.Sign(c.key, b)
	return transaction, nil
}
