package bitbuffer

import (
	"encoding/binary"
	"errors"
	"io"
)

const ByteSize = 8

type BitBuffer interface {
	At(int) (bool, error)
	AddBits(...bool)
	Len() int
	WriteTo(io.Writer) error
}

type MemoryBitBuffer struct {
	buffer   []byte
	lastByte []bool
}

func NewMemoryBitBuffer() *MemoryBitBuffer {
	return &MemoryBitBuffer{}
}

func (mbb *MemoryBitBuffer) Len() int {
	return len(mbb.buffer)*ByteSize + len(mbb.lastByte)
}

// Converts slice of booleans to little endian, having higher order bits (from left to right)
func boolSliceToByte(arr []bool) byte {
	if len(arr) > ByteSize {
		panic("len of bool slice greater than 8, can not convert to byte!")
	}
	var b byte
	for i := range arr {
		if arr[i] {
			b |= 1 << (ByteSize - i - 1)
		}
	}
	return b
}

func byteToBoolSlice(b byte) []bool {
	bits := make([]bool, ByteSize)
	for i := 0; i < ByteSize; i++ {
		bits[i] = (b & (1 << (ByteSize - 1 - i))) > 0
	}
	return bits
}

func (mbb *MemoryBitBuffer) AddBits(bs ...bool) {
	for _, b := range bs {
		mbb.lastByte = append(mbb.lastByte, b)
		if len(mbb.lastByte) == ByteSize {
			mbb.buffer = append(mbb.buffer, boolSliceToByte(mbb.lastByte))
			mbb.lastByte = mbb.lastByte[:0]
		}
	}
}

func (mbb *MemoryBitBuffer) At(idx int) (bool, error) {
	if idx < 0 || idx >= len(mbb.buffer)*ByteSize+len(mbb.lastByte) {
		return false, ErrOutOfRange
	}
	blockIdx := idx / ByteSize
	inBlockIdx := ByteSize - (idx % ByteSize) - 1
	if blockIdx < len(mbb.buffer) {
		return mbb.buffer[blockIdx]&(1<<inBlockIdx) > 0, nil
	} else {
		return mbb.lastByte[idx%ByteSize], nil
	}
}

// WriteTo writes buffer to writer in the way of (# of bits in last byte, bytes... in LittleEndian)
func (mbb *MemoryBitBuffer) WriteTo(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, byte(len(mbb.lastByte)))
	if err != nil {
		return err
	}
	for _, b := range mbb.buffer {
		if err = binary.Write(w, binary.LittleEndian, b); err != nil {
			return err
		}
	}
	if len(mbb.lastByte) > 0 {
		return binary.Write(w, binary.LittleEndian, boolSliceToByte(mbb.lastByte))
	}
	return nil
}

var _ BitBuffer = &MemoryBitBuffer{}

func ReadNewMemoryBitBuffer(r io.Reader) (*MemoryBitBuffer, error) {
	var lastByteSize byte
	if err := binary.Read(r, binary.LittleEndian, &lastByteSize); err != nil {
		return nil, err
	}
	mbb := NewMemoryBitBuffer()
	var lastReadByte byte
	err := binary.Read(r, binary.LittleEndian, &lastReadByte)
	for ; err == nil; err = binary.Read(r, binary.LittleEndian, &lastReadByte) {
		mbb.buffer = append(mbb.buffer, lastReadByte)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if lastByteSize > 0 {
		lastReadByte = mbb.buffer[len(mbb.buffer)-1]
		mbb.buffer = mbb.buffer[:len(mbb.buffer)-1]
		mbb.lastByte = byteToBoolSlice(lastReadByte)
	}
	return mbb, nil
}
