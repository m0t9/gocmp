package bits

import (
	"encoding/binary"
	"io"
)

type BitReader interface {
	ReadBit() (bool, error)
}

const (
	noReadWriteAttempt = -1
	byteSize           = 8
)

type BitReaderImpl struct {
	readBits int
	lastByte byte
	r        io.Reader
}

func NewBitReader(r io.Reader) BitReader {
	return &BitReaderImpl{readBits: noReadWriteAttempt, r: r}
}

func (br *BitReaderImpl) readByte() error {
	if !(br.readBits == noReadWriteAttempt || br.readBits == byteSize) {
		return nil
	}
	if err := binary.Read(br.r, binary.LittleEndian, &br.lastByte); err != nil {
		return err
	}
	br.readBits = 0
	return nil
}

func (br *BitReaderImpl) ReadBit() (bool, error) {
	if err := br.readByte(); err != nil {
		return false, err
	}
	bit := br.lastByte&(1<<br.readBits) > 0
	br.readBits++
	return bit, nil
}

var _ BitReader = &BitReaderImpl{}
