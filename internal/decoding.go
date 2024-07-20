package internal

import (
	"bufio"
	"go-compressor/pkg/bitbuffer"
	"io"
)

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

type HuffmanMemoryEncoderDecoder struct {
}

func NewHuffmanMemoryEncoderDecoder() *HuffmanMemoryEncoderDecoder {
	return &HuffmanMemoryEncoderDecoder{}
}

func (hmed *HuffmanMemoryEncoderDecoder) Encode(r io.ReadSeeker, w io.Writer) error {
	fa, err := NewFrequencyArray(r)
	if err != nil {
		return err
	}
	f := NewForest(fa)
	ht := NewHuffmanTree(f)
	if err := ht.WriteTo(w); err != nil {
		return err
	}
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	br := bufio.NewReader(r)
	mbb := bitbuffer.NewMemoryBitBuffer()
	for b, err := br.ReadByte(); err == nil; b, err = br.ReadByte() {
		mbb.AddBits(ht.CharEncoding(b)...)
	}
	if err != nil {
		return err
	}
	return mbb.WriteTo(w)
}

func (hmed *HuffmanMemoryEncoderDecoder) Decode(r io.Reader, w io.Writer) error {
	ht, err := ReadNewHuffmanTree(r)
	if err != nil {
		return err
	}
	mbb, err := bitbuffer.ReadNewMemoryBitBuffer(r)
	if err != nil {
		return err
	}
	bw := bufio.NewWriter(w)
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
	return nil
}

var _ EncoderDecoder = NewHuffmanMemoryEncoderDecoder()
