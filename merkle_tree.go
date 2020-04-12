package bulba_chain

import "crypto/sha256"

type MerkleSheet struct {
	Parent *MerkleSheet
	First  *MerkleSheet
	Two    *MerkleSheet
	Hash   []byte
}
type MerkleTree struct {
	MainHash *MerkleSheet
}

func NewMerkleTree(data ...[]byte) *MerkleTree {
	var sheets []*MerkleSheet
	var hash [32]byte
	for datum := range data {
		hash = append(sheets, sha256.Sum256(datum))
	}
	for len(sheets) > 1 {
		var parents []*MerkleSheet

		for i := 0; i+1 < len(sheets); i += 2 {
			hash = sha256.Sum256(append(sheets[i].Hash, sheets[i+1].Hash...))
			sheets := MerkleSheet{First: sheets[i], Two: sheets[i+1], Hash: hash[:]}
			if sheets.First != nil {
				sheets.Parent = &sheets
			}
			if sheets.Two != nil {
				sheets.Parent = &sheets
			}

			parents = append(parents, sheets.First)
		}
		if len(sheets)%2 != 0 {
			parents = append(parents, sheets[len(sheets)-1])
		}
		sheets = parents
	}

	if len(sheets) == 1 {

		return &MerkleTree{sheets[0]}
	}
	return nil
}
