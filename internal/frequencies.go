package internal

import (
	"bufio"
	"errors"
	"io"
)

const BytesCount = 256

type FrequencyCounter interface {
	FrequencyOf(byte) uint64
	Total() uint64
}

type FrequencyArray struct {
	byteFrequency [BytesCount]uint64
	totalCount    uint64
}

func NewFrequencyArray(r io.Reader) (*FrequencyArray, error) {
	fa := &FrequencyArray{}
	br := bufio.NewReader(r)
	b, err := br.ReadByte()
	for ; err == nil; b, err = br.ReadByte() {
		fa.byteFrequency[b]++
		fa.totalCount++
	}
	if !errors.Is(err, io.EOF) {
		return nil, err
	}
	return fa, nil
}

func (fa *FrequencyArray) FrequencyOf(b byte) uint64 {
	return fa.byteFrequency[b]
}

func (fa *FrequencyArray) Total() uint64 {
	return fa.totalCount
}

var _ FrequencyCounter = &FrequencyArray{}

type ForestTree struct {
	frequency uint64
	root      int16
	char      byte
}

type Forest struct {
	trees []ForestTree
}

func NewForest(fc FrequencyCounter) *Forest {
	f := &Forest{
		trees: make([]ForestTree, 0, BytesCount),
	}
	for b := 0; b < BytesCount; b++ {
		if fc.FrequencyOf(byte(b)) > 0 {
			f.trees = append(f.trees, ForestTree{
				frequency: fc.FrequencyOf(byte(b)),
				root:      int16(len(f.trees)),
				char:      byte(b),
			})
		}
	}
	return f
}

func (f *Forest) FindTwoWithMinFrequency() (*ForestTree, *ForestTree) {
	var m1, m2 *ForestTree
	for i := 0; i < f.Size(); i++ {
		if m1 == nil || m1.frequency >= f.trees[i].frequency {
			m2 = m1
			m1 = &f.trees[i]
		} else if m2 == nil || m2.frequency >= f.trees[i].frequency {
			m2 = &f.trees[i]
		}
	}
	return m1, m2
}

func (f *Forest) Size() int {
	return len(f.trees)
}

func (f *Forest) PopTree() {
	if f.Size() > 0 {
		f.trees = f.trees[:f.Size()-1]
	}
}
