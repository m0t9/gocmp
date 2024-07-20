package internal

import (
	"slices"
	"strings"
	"testing"
)

func TestNewHuffmanTree(t *testing.T) {
	for _, tt := range []struct {
		name              string
		input             string
		expectedTree      HuffmanTree
		expectedEncodings map[byte][]bool
	}{
		{
			name:  "UsualInput",
			input: "abacaba",
			expectedTree: HuffmanTree{
				nodes: []HuffmanNode{
					{
						left:   -1,
						right:  -1,
						parent: 4,
						char:   'a',
					},
					{
						left:   -1,
						right:  -1,
						parent: 3,
						char:   'b',
					},
					{
						left:   -1,
						right:  -1,
						parent: 3,
						char:   'c',
					},
					{
						left:   2,
						right:  1,
						parent: 4,
					},
					{
						left:   3,
						right:  0,
						parent: -1,
					},
				},
			},
			expectedEncodings: map[byte][]bool{
				'a': {true},
				'b': {false, true},
				'c': {false, false},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fa, err := NewFrequencyArray(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Unexpected error during FA building: %s", err)
			}
			f := NewForest(fa)
			tree := NewHuffmanTree(f)
			if len(tree.nodes) != len(tt.expectedTree.nodes) {
				t.Fatalf("Tree sizes differ: expected %d, got %d", len(tt.expectedTree.nodes), len(tree.nodes))
			}
			for i, node := range tree.nodes {
				if node != tt.expectedTree.nodes[i] {
					t.Errorf("%d-th node of HT differs: expected %v, got %v.",
						i, tt.expectedTree.nodes[i], node)
					continue
				}
				c := node.char
				if node.isLeaf() &&
					!slices.Equal(tt.expectedEncodings[c], tree.CharEncoding(c)) {
					t.Errorf("Encodings for `%c` differs: expected %v, got %v", c,
						tt.expectedEncodings[c], tree.nodeEncodings[c])
				}
			}
		})
	}
}
