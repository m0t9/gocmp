package internal

import (
	"bufio"
	"errors"
	"io"
)

const bytesCount = 256

type frequencyCounter interface {
	frequencyOf(byte) uint64
	total() uint64
}

type frequencyArray struct {
	byteFrequency [bytesCount]uint64
	totalCount    uint64
}

func newFrequencyArray(r io.Reader) (*frequencyArray, error) {
	fa := &frequencyArray{}
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

func (fa *frequencyArray) frequencyOf(b byte) uint64 {
	return fa.byteFrequency[b]
}

func (fa *frequencyArray) total() uint64 {
	return fa.totalCount
}

var _ frequencyCounter = &frequencyArray{}

type forestTree struct {
	frequency uint64
	root      int16
	char      byte
}

type forest struct {
	trees []forestTree
}

func newForest(fc frequencyCounter) *forest {
	f := &forest{
		trees: make([]forestTree, 0, bytesCount),
	}
	for b := 0; b < bytesCount; b++ {
		if fc.frequencyOf(byte(b)) > 0 {
			f.trees = append(f.trees, forestTree{
				frequency: fc.frequencyOf(byte(b)),
				root:      int16(len(f.trees)),
				char:      byte(b),
			})
		}
	}
	return f
}

func (f *forest) findTwoWithMinFrequency() (*forestTree, *forestTree) {
	var m1, m2 *forestTree
	for i := 0; i < f.size(); i++ {
		if m1 == nil || m1.frequency > f.trees[i].frequency {
			m2 = m1
			m1 = &f.trees[i]
		} else if m2 == nil || m2.frequency > f.trees[i].frequency {
			m2 = &f.trees[i]
		}
	}
	return m1, m2
}

func (f *forest) size() int {
	return len(f.trees)
}

func (f *forest) popTree() {
	if f.size() > 0 {
		f.trees = f.trees[:f.size()-1]
	}
}
