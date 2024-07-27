package internal

import (
	"bufio"
	"encoding/binary"
	"go-compressor/pkg/bits"
	"io"
)

const BufferSize = 1 << 16

type Encoder interface {
	Encode(io.ReadSeeker, io.Writer) error
}

type Decoder interface {
	Decode(io.Reader, io.Writer) error
}

type EncoderDecoder interface {
	Encoder
	Decoder
}

type HuffmanEncoderDecoder struct {
}

func NewHuffmanEncoderDecoder() EncoderDecoder {
	return &HuffmanEncoderDecoder{}
}

func (hmed *HuffmanEncoderDecoder) Encode(r io.ReadSeeker, w io.Writer) error {
	bw := bufio.NewWriterSize(w, BufferSize)
	fa, err := newFrequencyArray(r)
	if err != nil {
		return err
	}
	f := newForest(fa)
	ht := newHuffmanTree(f)
	if err := ht.writeTo(bw); err != nil {
		return err
	}
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	br := bufio.NewReaderSize(r, BufferSize)

	// write # of bytes in original file
	if err := binary.Write(bw, binary.LittleEndian, fa.total()); err != nil {
		return err
	}

	bitwr := bits.NewBitWriter(bw)
	for b, err := br.ReadByte(); err == nil; b, err = br.ReadByte() {
		if err := bitwr.WriteBits(ht.charEncoding(b)...); err != nil {
			return err
		}
	}
	if err := bitwr.Flush(); err != nil {
		return err
	}
	return bw.Flush()
}

func (hmed *HuffmanEncoderDecoder) Decode(r io.Reader, w io.Writer) error {
	ht, err := readNewHuffmanTree(r)
	if err != nil {
		return err
	}
	br := bufio.NewReaderSize(r, BufferSize)

	// read original file size
	var bytesCnt uint64
	if err := binary.Read(br, binary.LittleEndian, &bytesCnt); err != nil {
		return err
	}

	bw := bufio.NewWriterSize(w, BufferSize)
	bitr := bits.NewBitReader(br)
	node := ht.root()

	var writeBytes uint64
	for writeBytes < bytesCnt {
		if b, err := bitr.ReadBit(); err != nil {
			return err
		} else if !b {
			node = ht.getNode(int(node.left))
		} else {
			node = ht.getNode(int(node.right))
		}

		if node.isLeaf() {
			if _, err := bw.Write([]byte{node.char}); err != nil {
				return err
			}
			node = ht.root()
			writeBytes++
		}
	}
	return bw.Flush()
}

var _ EncoderDecoder = NewHuffmanEncoderDecoder()
