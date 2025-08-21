package trie

import (
	"bytes"
	"errors"
	"fmt"
)

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{root: &Node{Type: Empty}}
}

func (t *Trie) Root() *Node {
	return t.root
}

func (t *Trie) SetRoot(root *Node) {
	t.root = root
}

func KeyToNibbles(key []byte) []byte {
	nibbles := make([]byte, len(key)*2)
	for i, b := range key {
		nibbles[i*2] = b >> 4
		nibbles[i*2+1] = b & 0x0f
	}
	return nibbles
}

func CommonPrefix(a, b []byte) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return minLen
}

func (t *Trie) Update(rawKey, value []byte) {
	key := KeyToNibbles(rawKey)
	t.root = t.update(t.root, key, value)
}

func (t *Trie) update(node *Node, key, value []byte) *Node {
	if len(key) == 0 {
		if node.Type == Branch {
			node.Children[16] = &Node{Type: Leaf, Value: value}
			return node
		}
		return &Node{Type: Leaf, Value: value}
	}
	switch node.Type {
	case Empty:
		return &Node{Type: Leaf, Key: key, Value: value}
	case Leaf:
		prefixLen := CommonPrefix(node.Key, key)
		if prefixLen == len(node.Key) && prefixLen == len(key) {
			node.Value = value
			return node
		}
		newBranch := &Node{Type: Branch}
		if prefixLen == len(node.Key) {
			newBranch.Children[16] = &Node{Type: Leaf, Value: node.Value}
			remainingKey := key[prefixLen:]
			if len(remainingKey) > 0 {
				newBranch.Children[remainingKey[0]] = &Node{Type: Leaf, Key: remainingKey[1:], Value: value}
			} else {
				newBranch.Children[16] = &Node{Type: Leaf, Value: value}
			}
		} else {
			nodeRemaining := node.Key[prefixLen:]
			keyRemaining := key[prefixLen:]
			newBranch.Children[nodeRemaining[0]] = &Node{Type: Leaf, Key: nodeRemaining[1:], Value: node.Value}
			newBranch.Children[keyRemaining[0]] = &Node{Type: Leaf, Key: keyRemaining[1:], Value: value}
		}
		if prefixLen > 0 {
			return &Node{Type: Extension, Key: node.Key[:prefixLen], Child: newBranch}
		}
		return newBranch
	case Extension:
		prefixLen := CommonPrefix(node.Key, key)
		if prefixLen == len(node.Key) {
			node.Child = t.update(node.Child, key[prefixLen:], value)
			return node
		}
		newBranch := &Node{Type: Branch}
		nodeRemaining := node.Key[prefixLen:]
		keyRemaining := key[prefixLen:]
		newBranch.Children[nodeRemaining[0]] = &Node{Type: Extension, Key: nodeRemaining[1:], Child: node.Child}
		newBranch.Children[keyRemaining[0]] = &Node{Type: Leaf, Key: keyRemaining[1:], Value: value}
		if prefixLen > 0 {
			return &Node{Type: Extension, Key: node.Key[:prefixLen], Child: newBranch}
		}
		return newBranch
	case Branch:
		idx := key[0]
		child := node.Children[idx]
		if child == nil {
			child = &Node{Type: Empty}
		}
		node.Children[idx] = t.update(child, key[1:], value)
		return node
	}
	return node
}

func (t *Trie) RootHash() []byte {
	return t.root.ComputeHash()
}

func (t *Trie) Prove(rawKey []byte) ([][]byte, error) {
	key := KeyToNibbles(rawKey)
	proof := [][]byte{}
	node := t.root
	for {
		if node.Type == Empty {
			return nil, errors.New("key not found")
		}
		proof = append(proof, node.Serialize())
		fmt.Printf("Prove: added node type %v, serialized: %x\n", node.Type, node.Serialize())
		switch node.Type {
		case Leaf:
			if !bytes.Equal(node.Key, key) {
				return nil, errors.New("key mismatch in leaf")
			}
			return proof, nil
		case Extension:
			prefixLen := CommonPrefix(node.Key, key)
			if prefixLen != len(node.Key) {
				return nil, errors.New("key mismatch in extension")
			}
			key = key[prefixLen:]
			node = node.Child
		case Branch:
			if len(key) == 0 {
				if node.Children[16] == nil {
					return nil, errors.New("no value at branch end")
				}
				return proof, nil
			}
			idx := key[0]
			key = key[1:]
			node = node.Children[int(idx)]
			if node == nil {
				return nil, errors.New("key not found in branch")
			}
		default:
			return nil, errors.New("invalid node type")
		}
	}
}