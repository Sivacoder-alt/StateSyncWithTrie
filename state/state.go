package state

import (
	"bytes"
	"fmt"
	"state-sync/trie"
	"state-sync/utils"
)

type StateManager struct {
	Trie *trie.Trie
}

func NewStateManager(t *trie.Trie) *StateManager {
	return &StateManager{Trie: t}
}

func (sm *StateManager) SyncStates(states map[string]string) {
	for k, v := range states {
		sm.Trie.Update([]byte(k), []byte(v))
	}
}

func (sm *StateManager) GenerateProof(key []byte) ([][]byte, error) {
	return sm.Trie.Prove(key)
}

func (sm *StateManager) VerifyProof(rootHash, key []byte, proof [][]byte) ([]byte, bool) {
	if len(proof) == 0 {
		fmt.Println("proof is empty")
		return nil, false
	}
	h := utils.Keccak256(proof[0])
	if !bytes.Equal(h, rootHash) {
		fmt.Printf("proof: root hash mismatch, expected %x, got %x\n", rootHash, h)
		return nil, false
	}
	var value []byte
	i := 0
	nibbles := trie.KeyToNibbles(key)
	for i < len(proof) {
		nodeSer := proof[i]
		nodeType := trie.NodeType(nodeSer[0])
		fmt.Printf("proof: node %d, type %v\n", i, nodeType)
		pos := 1
		switch nodeType {
		case trie.Leaf:
			keyLen := int(nodeSer[pos])
			pos++
			nodeKey := nodeSer[pos : pos+keyLen]
			pos += keyLen
			valueLen := int(nodeSer[pos])
			pos++
			value = nodeSer[pos : pos+valueLen]
			if !bytes.Equal(nodeKey, nibbles) {
				fmt.Printf("proof: leaf key mismatch, expected %x, got %x\n", nibbles, nodeKey)
				return nil, false
			}
			fmt.Printf("proof: leaf with value %s\n", value)
			return value, true
		case trie.Extension:
			keyLen := int(nodeSer[pos])
			pos++
			nodeKey := nodeSer[pos : pos+keyLen]
			pos += keyLen
			childHash := nodeSer[pos : pos+32]
			prefixLen := trie.CommonPrefix(nodeKey, nibbles)
			if prefixLen != len(nodeKey) {
				fmt.Printf("proof: extension prefix mismatch, expected len %d, got %d\n", len(nodeKey), prefixLen)
				return nil, false
			}
			nibbles = nibbles[prefixLen:]
			i++
			if i >= len(proof) {
				fmt.Println("proof: Ran out of Proof nodes for extension")
				return nil, false
			}
			h = utils.Keccak256(proof[i])
			if !bytes.Equal(h, childHash) {
				fmt.Printf("proof: extension child hash mismatch, expected %x, got %x\n", childHash, h)
				return nil, false
			}
		case trie.Branch:
			var childHashes [16][]byte
			for j := 0; j < 16; j++ {
				childHashes[j] = nodeSer[pos : pos+32]
				pos += 32
			}
			valueLen := int(nodeSer[pos])
			pos++
			value = nodeSer[pos : pos+valueLen]
			if len(nibbles) == 0 {
				if valueLen == 0 {
					fmt.Println("proof: no value in branch for empty key")
					return nil, false
				}
				fmt.Printf("proof: found branch value %s\n", value)
				return value, true
			}
			idx := nibbles[0]
			nibbles = nibbles[1:]
			childHash := childHashes[idx]
			i++
			if i >= len(proof) {
				fmt.Println("proof: ran out of proof nodes for branch") 
				return nil, false
			}
			h = utils.Keccak256(proof[i])
			if !bytes.Equal(h, childHash) {
				fmt.Printf("proof: branch child hash mismatch, expected %x, got %x\n", childHash, h)
				return nil, false
			}
		default:
			fmt.Printf("proof: invalid node type %v\n", nodeType)
			return nil, false
		}
	}
	fmt.Println("proof: incomplete   path")
	return nil, false
}