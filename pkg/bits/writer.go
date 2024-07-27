package bits

import (
	"io"
)

type BitWriter interface {
	WriteBits(...bool) error
	Flush() error
}

type BitWriterImpl struct {
	lastByte  byte
	w         io.Writer
	writeBits int
}

func (bw *BitWriterImpl) writeByte() error {
	if bw.writeBits == 0 {
		return nil
	}

	if _, err := bw.w.Write([]byte{bw.lastByte}); err != nil {
		return err
	}
	bw.lastByte = 0
	bw.writeBits = 0
	return nil
}

func (bw *BitWriterImpl) WriteBits(bs ...bool) error {
	for _, b := range bs {
		if bw.writeBits == byteSize {
			if err := bw.writeByte(); err != nil {
				return err
			}
		}
		if b {
			bw.lastByte |= 1 << bw.writeBits
		}
		bw.writeBits++
	}
	return nil
}

func (bw *BitWriterImpl) Flush() error {
	return bw.writeByte()
}

func NewBitWriter(w io.Writer) BitWriter {
	return &BitWriterImpl{writeBits: 0, w: w}
}

var _ BitWriter = &BitWriterImpl{}
