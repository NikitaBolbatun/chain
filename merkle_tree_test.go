package chain

import (
	"bytes"
	"testing"
)

func TestNewMerkleTree_DifferentTreesProduceDifferentHashes(t *testing.T) {
	test := [][]byte{[]byte("MerkleTree"), []byte("Blockchain"), []byte("RootHash"), []byte("StateHash"), []byte("123312")}

	perms := permutations(len(test))
	results := make([][]byte, len(perms))
	tests := make([][][]byte, len(perms))
	for n, p := range perms {
		tc := make([][]byte, len(test))
		for i, idx := range p {
			tc[i] = test[idx]
		}
		tests[n] = tc

		mt := NewMerkleTree(tests[n]...)
		results[n] = mt.MainHash.Hash
	}

	first := results[0]
	for i, res := range results {
		if i == 0 {
			continue
		}

		if bytes.Equal(first, res) {
			t.Errorf("results for\n\n%v\n\tand\n%v\n\nare equal: %v", asStrings(tests[0]), asStrings(tests[i]), res)
		}
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
