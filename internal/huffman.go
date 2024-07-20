package internal

import (
	"encoding/binary"
	"io"
)

type HuffmanNode struct {
	left   int16
	right  int16
	parent int16
	char   byte
}

func readNewHuffmanNode(r io.Reader) (*HuffmanNode, error) {
	node := &HuffmanNode{}
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

func (hn *HuffmanNode) isLeaf() bool {
	return hn.left == hn.right && hn.left == -1
}

func (hn *HuffmanNode) writeTo(w io.Writer) error {
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

type HuffmanTree struct {
	nodes         []HuffmanNode
	nodeEncodings [BytesCount][]bool
}

func (ht *HuffmanTree) getNode(idx int) *HuffmanNode {
	return &ht.nodes[idx]
}

func (ht *HuffmanTree) root() *HuffmanNode {
	return ht.getNode(len(ht.nodes) - 1)
}

func (ht *HuffmanTree) buildEncodings() {
	ht.encodingDfs(ht.root(), nil)
}

func (ht *HuffmanTree) encodingDfs(node *HuffmanNode, encoding []bool) {
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

func NewHuffmanTree(f *Forest) *HuffmanTree {
	ht := &HuffmanTree{nodes: make([]HuffmanNode, 0, f.Size())}
	for _, t := range f.trees {
		ht.nodes = append(ht.nodes, HuffmanNode{
			left:   -1,
			right:  -1,
			parent: -1,
			char:   t.char,
		})
	}

	for ; f.Size() > 1; f.PopTree() {
		m1, m2 := f.FindTwoWithMinFrequency()
		newRoot := int16(len(ht.nodes))
		ht.nodes[m1.root].parent = newRoot
		ht.nodes[m2.root].parent = newRoot

		ht.nodes = append(ht.nodes, HuffmanNode{
			left:   m1.root,
			right:  m2.root,
			parent: -1,
		})

		m1.frequency += m2.frequency
		m1.root = newRoot

		*m2 = f.trees[f.Size()-1]
	}

	ht.buildEncodings()
	return ht
}

func (ht *HuffmanTree) CharEncoding(char byte) []bool {
	return ht.nodeEncodings[char]
}

func (ht *HuffmanTree) WriteTo(w io.Writer) error {
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

func ReadNewHuffmanTree(r io.Reader) (*HuffmanTree, error) {
	var tsz int16
	if err := binary.Read(r, binary.LittleEndian, &tsz); err != nil {
		return nil, err
	}
	ht := &HuffmanTree{nodes: make([]HuffmanNode, tsz)}
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
