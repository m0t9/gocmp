package internal

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewFrequencyArray(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
	}{
		{
			name:  "ValidString",
			input: "astringjustseveralletters",
		},
		{
			name:  "EmptyString",
			input: "",
		},
		{
			name: "AllPossibleBytesOnce",
			input: "\u0000\u0001\u0002\u0003\u0004\u0005\u0006\t\n\u000E\u000F\u0010\u0011\u0012" +
				"\u0013\u0014\u0015\u0017\u0018\u0019\u001A\u001B\u001C\u001D\u001E\u001F !\"#$%" +
				"&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrs" +
				"tuvwxyz{|}~\u007F\u0080\u0081\u0082\u0083\u0084\u0085\u0086\u0087\u0088\u0089" +
				"\u008A\u008B\u008C\u008D\u008E\u008F\u0090\u0091\u0092\u0093\u0094\u0095\u0096" +
				"\u0097\u0098\u0099\u009A\u009B\u009C\u009D\u009E\u009F ¡¢£¤¥¦§¨©ª«¬­®¯°±²" +
				"³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fa, err := NewFrequencyArray(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			for bi := 0; bi < BytesCount; bi++ {
				b := byte(bi)
				c := int(fa.FrequencyOf(b))
				ec := bytes.Count([]byte(tt.input), []byte{b})
				if c != ec {
					t.Errorf("Invalid count of byte `%c`. Expected %d, got %d.",
						b, ec, c)
				}
			}
		})
	}
}

func TestNewForest(t *testing.T) {
	for _, tt := range []struct {
		name           string
		input          string
		expectedForest Forest
	}{
		{
			name:  "ValidString",
			input: "justvalidstring",
			expectedForest: Forest{trees: []ForestTree{
				{
					frequency: 1,
					root:      0,
					char:      'a',
				},
				{
					frequency: 1,
					root:      1,
					char:      'd',
				},
				{
					frequency: 1,
					root:      2,
					char:      'g',
				},
				{
					frequency: 2,
					root:      3,
					char:      'i',
				},
				{
					frequency: 1,
					root:      4,
					char:      'j',
				},
				{
					frequency: 1,
					root:      5,
					char:      'l',
				},
				{
					frequency: 1,
					root:      6,
					char:      'n',
				},
				{
					frequency: 1,
					root:      7,
					char:      'r',
				},
				{
					frequency: 2,
					root:      8,
					char:      's',
				},
				{
					frequency: 2,
					root:      9,
					char:      't',
				},
				{
					frequency: 1,
					root:      10,
					char:      'u',
				},
				{
					frequency: 1,
					root:      11,
					char:      'v',
				},
			}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fc, err := NewFrequencyArray(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Unexpected error at FA building: %s", err)
			}
			forest := NewForest(fc)
			if forest.Size() != tt.expectedForest.Size() {
				t.Errorf("Expected and got forest sizes differ - %d and %d",
					tt.expectedForest.Size(), forest.Size())
			}
			for i := 0; i < forest.Size(); i++ {
				if forest.trees[i] != tt.expectedForest.trees[i] {
					t.Errorf("At position %d expected forest tree differs from got one - %+v and %+v",
						i, tt.expectedForest.trees[i], forest.trees[i])
				}
			}
		})
	}
}

func TestForest_FindTwoWithMinFrequency(t *testing.T) {
	for _, tt := range []struct {
		name   string
		input  string
		m1, m2 *struct {
			frequency uint64
			char      byte
		}
	}{
		{
			name:  "UsualString",
			input: "abacaba",
			m1: &struct {
				frequency uint64
				char      byte
			}{frequency: 1, char: 'c'},
			m2: &struct {
				frequency uint64
				char      byte
			}{frequency: 2, char: 'b'},
		},
		{
			name:  "EquallyMinimal",
			input: "bbaaccc",
			m1: &struct {
				frequency uint64
				char      byte
			}{frequency: 2, char: 'b'},
			m2: &struct {
				frequency uint64
				char      byte
			}{frequency: 2, char: 'a'},
		},
		{
			name:  "NoSecondMin",
			input: "aaaaa",
			m1: &struct {
				frequency uint64
				char      byte
			}{frequency: 5, char: 'a'},
		},
		{
			name:  "NoMin",
			input: "",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fa, err := NewFrequencyArray(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Unexpected error at FA building: %s", err)
			}
			f := NewForest(fa)
			m1, m2 := f.FindTwoWithMinFrequency()
			m1notExistence := (tt.m1 == nil && m1 != nil) || (tt.m1 != nil && m1 == nil)
			m2notExistence := (tt.m2 == nil && m2 != nil) || (tt.m2 != nil && m2 == nil)
			if m1notExistence || m2notExistence {
				t.Fatalf(
					"Problems with minimal values pointers existence: m1 %p, m2 %p, but expected m1 %p, m2 %p",
					m1, m2, tt.m1, tt.m2)
			}
			if m1 != nil && (m1.frequency != tt.m1.frequency || m1.char != tt.m1.char) {
				t.Errorf("First minimum expected value differs.\n"+
					"Frequency = %d, Expected frequency = %d\n"+
					"Character = %c, Expected character = %c",
					m1.frequency, tt.m1.frequency, m1.char, tt.m1.char)
			}
			if m2 != nil && (m2.frequency != tt.m2.frequency || m2.char != tt.m2.char) {
				t.Errorf("Second minimum expected value differs.\n"+
					"Frequency = %d, Expected frequency = %d\n"+
					"Character = %c, Expected character = %c",
					m2.frequency, tt.m2.frequency, m2.char, tt.m2.char)
			}
		})
	}
}
