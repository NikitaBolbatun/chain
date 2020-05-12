package bulba_chain

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
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
//Исправить костыль с 2 функцией addtransaction
func (c *Node) AddTransaction(transaction Transaction) error {
	c.transMutex.Lock()
	defer c.transMutex.Unlock()

	hash, err := transaction.Hash()
	if err != nil {
		return err
	}

	err = c.CheckTransaction(transaction)
	if err != nil {
		return err
	}

	c.transactionPool[hash] = transaction

	ctx := context.Background()

	fmt.Println("Add transaction transaction pool")

	c.Broadcast(ctx, Message{
		From: c.address,
		Data: TransactionSend{
			NodeName:    c.address,
			Transaction: transaction,
		},
	})
	return nil
}

func (c *Node) SignTransaction(transaction Transaction)  (Transaction,error) {
	b, err := transaction.Bytes()
	if err != nil {
		return  transaction,err
	}

	transaction.Signature = ed25519.Sign(c.key, b)

	return  transaction,nil
}

func (c *Node) CheckTransaction(transaction Transaction) error {
	if transaction.To == "" || transaction.From == "" {
		return errors.New("username not correct")
	}
	balance:= c.state[transaction.From]
	if balance < transaction.Fee + transaction.Amount {
		return  errors.New("balance not correct" )
	}
	return nil
}