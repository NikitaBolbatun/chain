package bulba_chain

import "crypto/sha256"

type MerkleNode struct {
	Parent *MerkleNode
	Left   *MerkleNode
	Right  *MerkleNode
	Hash   []byte
}

type MerkleTree struct {
	Root *MerkleNode
}

func NewMerkleTree(data ...[]byte) *MerkleTree {
	var nodes []*MerkleNode

	for _, datum := range data {
		nodes = append(nodes, newMerkleNode(nil, nil, datum))
	}

	for len(nodes) > 1 {
		var parents []*MerkleNode

		for i := 0; i+1 < len(nodes); i += 2 {
			node := newMerkleNode(nodes[i], nodes[i+1], append(nodes[i].Hash, nodes[i+1].Hash...))
			parents = append(parents, node)
		}

		if len(nodes)%2 != 0 {
			parents = append(parents, nodes[len(nodes)-1])
		}

		nodes = parents
	}

	if len(nodes) == 1 {
		return &MerkleTree{Root: nodes[0]}
	}
	return nil
}
func newMerkleNode(left *MerkleNode, right *MerkleNode, data []byte) *MerkleNode {
	var hash [32]byte

	if left == nil && right == nil {
		hash = sha256.Sum256(data)
	} else {
		hash = sha256.Sum256(append(left.Hash, right.Hash...))
	}

	node := MerkleNode{Left: left, Right: right, Hash: hash[:]}
	if left != nil {
		left.Parent = &node
	}
	if right != nil {
		right.Parent = &node
	}

	return &node
}
