package bits

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func stringBit(s string, idx int) bool {
	return s[idx/byteSize]&(1<<(idx%byteSize)) > 0
}

func TestBitReader(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
	}{
		{
			name: "EnglishAlphabet",
			data: "abcdefghijklmnopqrstuvwxyz",
		},
		{
			name: "SpecialCharacters",
			data: "!!@#$%^&*(),.<>/`[]",
		},
		{
			name: "OneLetter",
			data: "a",
		},
		{
			name: "NoLetters",
			data: "",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.data)
			br := NewBitReader(r)
			for idx := 0; idx < len(tt.data)*byteSize; idx++ {
				if b, err := br.ReadBit(); err != nil {
					t.Fatalf("Unexpected error on data reading: %s",
						err)
				} else if exp := stringBit(tt.data, idx); exp != b {
					t.Errorf("Expected and read bit differs: expected %t, got %t",
						exp, b)
				}
			}
			if _, err := br.ReadBit(); !errors.Is(err, io.EOF) {
				t.Errorf("For some reason, BitReader has not reached EOF")
			}
		})
	}
}

func TestBitWriteRead(t *testing.T) {
	td := t.TempDir()
	for _, tt := range []struct {
		name string
		data []bool
	}{
		{
			name: "8Bits",
			data: []bool{true, false, false, true, true, false, false, true},
		},
		{
			name: "5Bits",
			data: []bool{true, false, false, true, true},
		},
		{
			name: "9Bits",
			data: []bool{true, false, true, true, true, false, false, true, true},
		},
		{
			name: "0Bits",
			data: []bool{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.CreateTemp(td, tt.name)
			bw := NewBitWriter(f)
			if err := bw.WriteBits(tt.data...); err != nil {
				t.Fatalf("Unexpected error during write: %s", err)
			}
			if err := bw.Flush(); err != nil {
				t.Fatalf("Unexpected error during flushing: %s", err)
			}
			_, _ = f.Seek(0, 0)

			br := NewBitReader(f)
			for i, b := range tt.data {
				if rb, err := br.ReadBit(); err != nil {
					t.Fatalf("Unexpected error during read: %s", err)
				} else if rb != b {
					t.Errorf("%d-th read bit differs from expected: want %t, got %t", i, b, rb)
				}
			}
			for b, err := br.ReadBit(); ; b, err = br.ReadBit() {
				if errors.Is(err, io.EOF) {
					return
				} else if err != nil {
					t.Fatalf("Unexpected error while reading remaining bits: %s", err)
				} else if b != false {
					t.Fatalf("Remaining bits are non-zero")
				}
			}
		})
	}
}
