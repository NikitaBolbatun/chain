package bulba_chain

import (
	"testing"
)

func TestExampleNewMerkleTree(t *testing.T) {
	tests := [][]byte{[]byte("MerkleTree"), []byte("Blockchain"), []byte("RootHash"), []byte("StateHash"), []byte("123312")}
	mt := NewMerkleTree(tests[5])

	t.Log(mt)
}
