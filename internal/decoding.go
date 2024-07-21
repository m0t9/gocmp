package internal

import (
	"bufio"
	"go-compressor/pkg/bitbuffer"
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

func NewHuffmanEncoderDecoder() *HuffmanEncoderDecoder {
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
	mbb := bitbuffer.NewMemoryBitBuffer()
	for b, err := br.ReadByte(); err == nil; b, err = br.ReadByte() {
		mbb.AddBits(ht.charEncoding(b)...)
	}
	if err != nil {
		return err
	}
	if err := mbb.WriteTo(bw); err != nil {
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
	mbb, err := bitbuffer.ReadNewMemoryBitBuffer(br)
	if err != nil {
		return err
	}
	bw := bufio.NewWriterSize(w, BufferSize)
	node := ht.root()
	var bts []byte
	for i := 0; i < mbb.Len(); i++ {
		if node.isLeaf() {
			bts = append(bts, node.char)
			node = ht.root()
		}
		b, _ := mbb.At(i)
		if !b {
			node = ht.getNode(int(node.left))
		} else {
			node = ht.getNode(int(node.right))
		}
	}

	if _, err := bw.Write(bts); err != nil {
		return err
	}
	return bw.Flush()
}

var _ EncoderDecoder = NewHuffmanEncoderDecoder()
