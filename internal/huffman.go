package internal

import (
	"encoding/binary"
	"io"
)

type huffmanNode struct {
	left   int16
	right  int16
	parent int16
	char   byte
}

func readNewHuffmanNode(r io.Reader) (*huffmanNode, error) {
	node := &huffmanNode{}
	if err := binary.Read(r, binary.LittleEndian, &node.left); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &node.right); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &node.parent); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &node.char); err != nil {
		return nil, err
	}
	return node, nil
}

func (hn *huffmanNode) isLeaf() bool {
	return hn.left == hn.right && hn.left == -1
}

func (hn *huffmanNode) writeTo(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, hn.left); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, hn.right); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, hn.parent); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, hn.char); err != nil {
		return err
	}
	return nil
}

type huffmanTree struct {
	nodes         []huffmanNode
	nodeEncodings [bytesCount][]bool
}

func (ht *huffmanTree) getNode(idx int) *huffmanNode {
	return &ht.nodes[idx]
}

func (ht *huffmanTree) root() *huffmanNode {
	return ht.getNode(len(ht.nodes) - 1)
}

func (ht *huffmanTree) buildEncodings() {
	ht.encodingDfs(ht.root(), nil)
}

func (ht *huffmanTree) encodingDfs(node *huffmanNode, encoding []bool) {
	if node.left == node.right && node.left == -1 {
		ht.nodeEncodings[node.char] = make([]bool, len(encoding))
		copy(ht.nodeEncodings[node.char], encoding)
		return
	}
	if node.left != -1 {
		encoding = append(encoding, false)
		ht.encodingDfs(&ht.nodes[node.left], encoding)
		encoding = encoding[:len(encoding)-1]
	}
	if node.right != -1 {
		encoding = append(encoding, true)
		ht.encodingDfs(&ht.nodes[node.right], encoding)
		encoding = encoding[:len(encoding)-1]
	}
}

func newHuffmanTree(f *forest) *huffmanTree {
	ht := &huffmanTree{nodes: make([]huffmanNode, 0, f.size())}
	for _, t := range f.trees {
		ht.nodes = append(ht.nodes, huffmanNode{
			left:   -1,
			right:  -1,
			parent: -1,
			char:   t.char,
		})
	}

	for ; f.size() > 1; f.popTree() {
		m1, m2 := f.findTwoWithMinFrequency()
		newRoot := int16(len(ht.nodes))
		ht.nodes[m1.root].parent = newRoot
		ht.nodes[m2.root].parent = newRoot

		ht.nodes = append(ht.nodes, huffmanNode{
			left:   m1.root,
			right:  m2.root,
			parent: -1,
		})

		m1.frequency += m2.frequency
		m1.root = newRoot

		*m2 = f.trees[f.size()-1]
	}

	ht.buildEncodings()
	return ht
}

func (ht *huffmanTree) charEncoding(char byte) []bool {
	return ht.nodeEncodings[char]
}

func (ht *huffmanTree) writeTo(w io.Writer) error {
	tsz := int16(len(ht.nodes))
	if err := binary.Write(w, binary.LittleEndian, tsz); err != nil {
		return err
	}
	for _, node := range ht.nodes {
		if err := node.writeTo(w); err != nil {
			return err
		}
	}
	return nil
}

func readNewHuffmanTree(r io.Reader) (*huffmanTree, error) {
	var tsz int16
	if err := binary.Read(r, binary.LittleEndian, &tsz); err != nil {
		return nil, err
	}
	ht := &huffmanTree{nodes: make([]huffmanNode, tsz)}
	for i := int16(0); i < tsz; i++ {
		if nodePtr, err := readNewHuffmanNode(r); err == nil {
			ht.nodes[i] = *nodePtr
		} else {
			return nil, err
		}
	}
	ht.buildEncodings()
	return ht, nil
}
