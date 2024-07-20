package bitbuffer

import (
	"io"
	"os"
	"testing"
)

func TestMemoryBitBuffer(t *testing.T) {
	for _, tt := range []struct {
		name string
		bits []bool
	}{
		{
			name: "8Bits",
			bits: []bool{false, true, false, false, true, true, false, false},
		},
		{
			name: "9Bits",
			bits: []bool{false, true, true, false, false, true, true, true, true},
		},
		{
			name: "ManyBits",
			bits: []bool{
				false, true, true, false, false, true, true,
				true, true, false, true, true, false, false,
				true, false, false, false, true, false, true,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bb := NewMemoryBitBuffer()
			for _, bit := range tt.bits {
				bb.AddBits(bit)
			}
			if bb.Len() != len(tt.bits) {
				t.Fatalf("Expected written number of bits differ from got: expected %d, got %d",
					len(tt.bits), bb.Len())
			}
			for idx, bit := range tt.bits {
				bbit, err := bb.At(idx)
				if err != nil {
					t.Errorf("Unexpected error reading bit with index %d: %s", idx, err)
				} else if bbit != bit {
					t.Errorf("Written and expected bits differ: expected %t, got %t.", bit, bbit)
				}
			}
		})
	}
}

func TestMemoryBitBufferWriteRead(t *testing.T) {
	td := t.TempDir()

	for _, tt := range []struct {
		name string
		bits []bool
	}{
		{
			name: "8Bits",
			bits: []bool{false, true, false, false, true, true, false, false},
		},
		{
			name: "9Bits",
			bits: []bool{false, true, true, false, false, true, true, true, true},
		},
		{
			name: "ManyBits",
			bits: []bool{
				false, true, true, false, false, true, true,
				true, true, false, true, true, false, false,
				true, false, false, false, true, false, true,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.CreateTemp(td, tt.name)
			var bb BitBuffer = NewMemoryBitBuffer()
			bb.AddBits(tt.bits...)
			if err := bb.WriteTo(f); err != nil {
				t.Fatalf("Unexpected error %s during writing memory bit buffer to file", err)
			}
			_, _ = f.Seek(0, io.SeekStart)

			if rbb, err := ReadNewMemoryBitBuffer(f); err != nil {
				t.Fatalf("Unexpected error %s during reading memory bit buffer from file", err)
			} else {
				for i, b := range tt.bits {
					if rb, _ := rbb.At(i); rb != b {
						t.Errorf("%d-th read bit differs from expected: %t expected, but got %t", i, b, rb)
					}
				}
			}
		})
	}
}
