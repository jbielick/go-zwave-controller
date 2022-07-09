package api

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func commandType(t CommandType) *CommandType {
	return &t
}

type Assertion struct {
	Len      int
	Checksum byte
}

type FrameCase struct {
	Frame     *Frame
	Assertion *Assertion
}

func TestFrameLength(t *testing.T) {
	testCases := map[interface{}]*FrameCase{
		"GetLibraryVersionRequest": {
			NewDataFrame(Request, GetLibraryVersion, []byte{}),
			&Assertion{Len: 3, Checksum: byte(0xe9)},
		},
		"GetLibraryVersionResponse": {
			NewDataFrame(Response, GetLibraryVersion, []byte("Z-Wave 6.07\x00\x01")),
			&Assertion{Len: 16, Checksum: byte(0x97)},
		},
	}
	for title, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", title), func(t *testing.T) {
			Equal(t, testCase.Assertion.Len, testCase.Frame.Length())
		})
	}
}

func TestFrameChecksum(t *testing.T) {
	testCases := map[interface{}]*FrameCase{
		"GetLibraryVersionRequest": {
			NewDataFrame(Request, GetLibraryVersion, []byte{}),
			&Assertion{Len: 3, Checksum: byte(0xe9)},
		},
		"GetLibraryVersionResponse": {
			NewDataFrame(Response, GetLibraryVersion, []byte("Z-Wave 6.07\x00\x01")),
			&Assertion{Len: 16, Checksum: byte(0x97)},
		},
	}
	for title, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", title), func(t *testing.T) {
			Equal(t, testCase.Assertion.Checksum, testCase.Frame.Checksum())
		})
	}
}

func BenchmarkChecksumLongFrame(b *testing.B) {
	frame := NewRequest(GetLibraryVersion, []byte{
		0x5a,
		0x2d,
		0x57,
		0x61,
		0x76,
		0x65,
		0x20,
		0x34,
		0x2e,
		0x30,
		0x35,
		0x00,
		0x07,
	})

	for n := 0; n < b.N; n++ {
		frame.Checksum()
	}
}
