package chain

import "crypto/sha256"

type MerkleSheet struct {
	Parent *MerkleSheet
	First  *MerkleSheet
	Two    *MerkleSheet
	Data   []byte
	Hash   []byte
}

type MerkleTree struct {
	MainHash *MerkleSheet
}

func NewMerkleTree(data ...[]byte) *MerkleTree {
	var sheets []*MerkleSheet

	for _, datum := range data{
		sheets = append(sheets, NewMerkleSheets(datum))
	}

	for len(sheets) > 1 {
		var parents []*MerkleSheet
		var hash [32]byte
		for i := 0; i+1 < len(sheets); i += 2 {
			hash = sha256.Sum256(append(sheets[i].Hash, sheets[i+1].Hash...))
			sheet := MerkleSheet{First: sheets[i], Two: sheets[i+1], Hash: hash[:]}
			if sheet.First != nil {
				sheet.Parent = &sheet
			}
			if sheet.Two != nil {
				sheet.Parent = &sheet
			}

			parents = append(parents, sheet.First)
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

func NewMerkleSheets(data []byte) *MerkleSheet {
	hash := sha256.Sum256(data)
	d := append([]byte(nil), data...)
	sheet := MerkleSheet{ Hash: hash[:], Data: d}
	return &sheet
}
