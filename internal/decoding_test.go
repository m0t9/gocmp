package internal

import (
	"bytes"
	"io"
	"os"
	"testing"
)

const TestBufferSize = 1 << 16

func checkFilesEquality(f1, f2 *os.File) bool {
	_, _ = f1.Seek(0, io.SeekStart)
	_, _ = f2.Seek(0, io.SeekStart)
	for {
		b1 := make([]byte, TestBufferSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, TestBufferSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				panic("files comparison check failure")
			}
		}
		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func TestHuffmanEncodeDecode(t *testing.T) {
	td := t.TempDir()
	for _, tt := range []struct {
		name           string
		fileToCompress string
	}{
		{
			name:           "VimBookPDF",
			fileToCompress: "../test/vimbook.pdf",
		},
		{
			name:           "DoraJPG",
			fileToCompress: "../test/dora.jpg",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inf, _ := os.Open(tt.fileToCompress)
			cf, _ := os.CreateTemp(td, tt.name)
			df, _ := os.CreateTemp(td, tt.name)

			hed := NewHuffmanEncoderDecoder()
			if err := hed.Encode(inf, cf); err != nil {
				t.Fatalf("Unexpected encoding error: %s", err)
			}
			_, _ = cf.Seek(0, io.SeekStart)
			if err := hed.Decode(cf, df); err != nil {
				t.Fatalf("Unexpected decoding error: %s", err)
			}
			if !checkFilesEquality(inf, df) {
				t.Fatalf("Initial and uncompressed files are different!")
			}
		})
	}
}
