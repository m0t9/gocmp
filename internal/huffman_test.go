package internal

import (
	"io"
	"os"
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
		{
			name:  "WikiTest",
			input: "aaaaaaaaaaaaaaabbbbbbbccccccddddddeeeee",
			expectedTree: HuffmanTree{
				nodes: []HuffmanNode{
					{
						left:   -1,
						right:  -1,
						parent: 8,
						char:   'a',
					},
					{
						left:   -1,
						right:  -1,
						parent: 6,
						char:   'b',
					},
					{
						left:   -1,
						right:  -1,
						parent: 5,
						char:   'c',
					},
					{
						left:   -1,
						right:  -1,
						parent: 6,
						char:   'd',
					},
					{
						left:   -1,
						right:  -1,
						parent: 5,
						char:   'e',
					},
					{
						left:   4,
						right:  2,
						parent: 7,
					},
					{
						left:   3,
						right:  1,
						parent: 7,
					},
					{
						left:   5,
						right:  6,
						parent: 8,
					},
					{
						left:   0,
						right:  7,
						parent: -1,
					},
				},
			},
			expectedEncodings: map[byte][]bool{
				'a': {false},
				'b': {true, true, true},
				'c': {true, false, true},
				'd': {true, true, false},
				'e': {true, false, false},
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
					t.Errorf("%d-th node of HT differs: expected %+v, got %+v.",
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

func TestHuffmanTreeWriteRead(t *testing.T) {
	td := t.TempDir()
	for _, tt := range []struct {
		name  string
		input string
	}{
		{
			name:  "UsualInput",
			input: "abacaba",
		},
		{
			name:  "MoreComplexInput",
			input: "\u0000\u0001\u0002\u0003\u0004\u0000\u0001\u0002\u0003\u0004AAAAAAAAAAAAAA",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src, _ := os.CreateTemp(td, tt.name)
			fa, err := NewFrequencyArray(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Unexpected error during FA building: %s", err)
			}
			f := NewForest(fa)
			tree := NewHuffmanTree(f)
			if err = tree.WriteTo(src); err != nil {
				t.Fatalf("Unexpected error during HT writing: %s", err)
			}
			_, _ = src.Seek(0, io.SeekStart)
			readTree, err := ReadNewHuffmanTree(src)
			if err != nil {
				t.Fatalf("Unexpected error during HT reading: %s", err)
			}
			if len(readTree.nodes) != len(tree.nodes) {
				t.Fatalf("Size of read HT differs from expected: got %d, expected %d",
					len(readTree.nodes), len(tree.nodes))
			}
			for i := range readTree.nodes {
				if readTree.nodes[i] != tree.nodes[i] {
					t.Errorf("%d-th node in read HT differs from expected: %v got, expected %v",
						i, readTree.nodes[i], tree.nodes[i])
				}
			}
		})
	}
}
