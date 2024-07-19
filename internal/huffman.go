package internal

type HuffmanNode struct {
	left   int16
	right  int16
	parent int16
	char   byte
}

type HuffmanTree struct {
	nodes         []HuffmanNode
	nodeEncodings [BytesCount][]bool
}

func (ht *HuffmanTree) buildEncodings() {
	ht.encodingDfs(&ht.nodes[len(ht.nodes)-1], nil)
}

func (ht *HuffmanTree) encodingDfs(node *HuffmanNode, encoding []bool) {
	if node.left == node.right && node.left == -1 {
		ht.nodeEncodings[node.char] = make([]bool, 0, len(encoding))
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
		ht.nodes = append(ht.nodes, HuffmanNode{
			left:   m1.root,
			right:  m2.root,
			parent: -1,
		})
		newRoot := int16(len(ht.nodes) - 1)
		ht.nodes[m1.root].parent = newRoot
		ht.nodes[m2.root].parent = newRoot

		m1.frequency += m2.frequency
		m1.root = newRoot

		*m2 = f.trees[f.Size()-1]
	}

	ht.buildEncodings()
	return ht
}
