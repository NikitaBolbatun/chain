package bulba_chain

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"testing"
)

//func TestMerkleNode_FindNode(t *testing.T) {
//	pubkey, _, err := ed25519.GenerateKey(nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	nd := &Node{
//		validators: []ed25519.PublicKey{pubkey},
//	}
//	nd.state = map[string]uint64{
//		"one":         20000,
//		"two":         20000,
//	}
//	nd.transactionPool = map[string]Transaction{}
//	for i:=0; i < 5; i++ {
//		transactions := NewTransaction("one", "two", uint64(i+10), uint64(i), nil, nil)
//		err = nd.AddTransaction(transactions)
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//	merkleBytes,err :=Bytes(nd.transactionPool)
//	if err != nil {
//		t.Fatal(err)
//	}
//	merkle := nd.NewMerkleTree(merkleBytes)
//	transaction := NewTransaction("one", "two", uint64(12), uint64(2), nil, nil)
//	hash, err := transaction.Hash()
//	if err != nil {
//		t.Fatal(err)
//	}
//	merkleTransaction := merkle.Root.FindNode(hash)
//	fmt.Println(merkleTransaction.FindNode(f))
//}
func TestMerkleNode(t *testing.T) {
	pubkey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	nd := &Node{
		validators: []ed25519.PublicKey{pubkey},
	}
	nd.state = map[string]uint64{
		"one":         20000,
		"two":         20000,
	}
	nd.transactionPool = map[string]Transaction{}
	for i:=0; i < 5; i++ {
		transactions := NewTransaction("one", "two", uint64(i+10), uint64(i), nil, nil)
		err = nd.AddTransaction(transactions)
		if err != nil {
			t.Fatal(err)
		}
	}
	merkleBytes,err :=Bytes(nd.transactionPool)
	if err != nil {
		t.Fatal(err)
	}
	merkle := nd.NewMerkleTree(merkleBytes)

	fmt.Println(merkle)
}

func TestNewMerkleTree_DifferentTreesProduceDifferentHashes(t *testing.T) {
	test := [][]byte{[]byte("MerkleTree"), []byte("Blockchain"), []byte("RootHash"), []byte("StateHash"), []byte("123312")}
	pubkey, _, _:= ed25519.GenerateKey(nil)
	nd := &Node{
		validators: []ed25519.PublicKey{pubkey},
	}
	perms := permutations(len(test))
	results := make([][]byte, len(perms))
	tests := make([][][]byte, len(perms))
	for n, p := range perms {
		tc := make([][]byte, len(test))
		for i, idx := range p {
			tc[i] = test[idx]
		}
		tests[n] = tc

		mt := nd.NewMerkleTree(tests[n]...)
		results[n] = mt.Root.Hash
	}

	for j := range results {
		j := j
		t.Run(fmt.Sprintf("first word with an index %d - %q", j, string(tests[j][0])), func(t *testing.T) {
			first := results[j]
			for i, res := range results {
				if i == j {
					continue
				}

				if bytes.Equal(first, res) {
					t.Errorf("results for\n\n%v\n\tand\n%v\n\nare equal: %v", asStrings(tests[j]), asStrings(tests[i]), res)
				}
			}
		})
	}
}

func asStrings(b [][]byte) []string {
	res := make([]string, len(b))
	for i, slice := range b {
		res[i] = string(slice)
	}
	return res
}

func permutations(n int) [][]int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}

	var helper func([]int, int)
	var res [][]int

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
