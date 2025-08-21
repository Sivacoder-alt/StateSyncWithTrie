package trie

import (
	"bytes"
	"state-sync/utils"
)

type NodeType int

const (
	Empty     NodeType = 0
	Leaf      NodeType = 1
	Extension NodeType = 2
	Branch    NodeType = 3
)

type Node struct {
	Type     NodeType
	Key      []byte    
	Value    []byte    
	Child    *Node     
	Children [17]*Node 
}

func (n *Node) Serialize() []byte {
	var b bytes.Buffer
	b.WriteByte(byte(n.Type))
	switch n.Type {
	case Leaf:
		b.WriteByte(byte(len(n.Key)))
		b.Write(n.Key)
		b.WriteByte(byte(len(n.Value)))
		b.Write(n.Value)
	case Extension:
		b.WriteByte(byte(len(n.Key)))
		b.Write(n.Key)
		b.Write(n.Child.ComputeHash())
	case Branch:
		for i := 0; i < 16; i++ {
			if n.Children[i] != nil {
				b.Write(n.Children[i].ComputeHash())
			} else {
				b.Write(make([]byte, 32)) 
			}
		}
		if n.Children[16] != nil {
			b.WriteByte(byte(len(n.Children[16].Value)))
			b.Write(n.Children[16].Value)
		} else {
			b.WriteByte(0)
		}
	}
	return b.Bytes()
}

func (n *Node) ComputeHash() []byte {
	if n.Type == Empty {
		return make([]byte, 32)
	}
	return utils.Keccak256(n.Serialize())
}