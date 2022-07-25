package transport

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

type CalculatedAssertion struct {
	Len      int
	Checksum byte
}

func TestFrameLength(t *testing.T) {
	testCases := map[interface{}]struct {
		Frame               *Frame
		CalculatedAssertion *CalculatedAssertion
	}{
		"GetLibraryVersionRequest": {
			NewDataFrame(Request, []byte{0x15}),
			&CalculatedAssertion{Len: 3, Checksum: byte(0xe9)},
		},
		"GetLibraryVersionResponse": {
			NewDataFrame(Response, []byte("\x15Z-Wave 6.07\x00\x01")),
			&CalculatedAssertion{Len: 16, Checksum: byte(0x97)},
		},
	}
	for title, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", title), func(t *testing.T) {
			assert.Equal(t, testCase.CalculatedAssertion.Len, testCase.Frame.Length())
		})
	}
}

func TestFrameChecksum(t *testing.T) {
	testCases := map[interface{}]struct {
		Frame               *Frame
		CalculatedAssertion *CalculatedAssertion
	}{
		"GetLibraryVersionRequest": {
			NewDataFrame(Request, []byte{0x15}),
			&CalculatedAssertion{Len: 3, Checksum: byte(0xe9)},
		},
		"GetLibraryVersionResponse": {
			NewDataFrame(Response, []byte("\x15Z-Wave 6.07\x00\x01")),
			&CalculatedAssertion{Len: 16, Checksum: byte(0x97)},
		},
	}
	for title, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", title), func(t *testing.T) {
			assert.Equal(t, testCase.CalculatedAssertion.Checksum, testCase.Frame.Checksum())
		})
	}
}

func TestMarshalBinary(t *testing.T) {
	testCases := map[interface{}]struct {
		Frame *Frame
		Bytes []byte
	}{
		"GetLibraryVersionRequest": {
			NewDataFrame(Request, []byte{0x15}),
			[]byte("\x01\x03\x00\x15\xe9"),
		},
		"GetLibraryVersionResponse": {
			NewDataFrame(Response, []byte("\x15Z-Wave 6.07\x00\x01")),
			[]byte("\x01\x10\x01\x15\x5a\x2d\x57\x61\x76\x65\x20\x36\x2e\x30\x37\x00\x01\x97"),
		},
	}
	for title, testCase := range testCases {
		t.Run(fmt.Sprintf("MarshalBinary/%s", title), func(t *testing.T) {
			data, err := testCase.Frame.MarshalBinary()
			assert.NoError(t, err)
			assert.Equal(t, testCase.Bytes, data)
		})
		t.Run(fmt.Sprintf("UnmarshalBinary/%s", title), func(t *testing.T) {
			got := &Frame{FrameType: SOF}
			err := got.UnmarshalBinary(testCase.Bytes)
			assert.NoError(t, err)
			assert.Equal(t, testCase.Frame, got)
		})
	}
}

func TestFrameIsACK(t *testing.T) {
	f := NewACK()
	assert.Equal(t, true, f.IsACK())
}

func TestFrameIsNAK(t *testing.T) {
	f := NewNAK()
	assert.Equal(t, true, f.IsNAK())
}

func TestFrameIsCAN(t *testing.T) {
	f := NewCAN()
	assert.Equal(t, true, f.IsCAN())
}

func TestFrameIsRequest(t *testing.T) {
	f := NewRequest([]byte{0x02})
	assert.Equal(t, true, f.IsRequest())
}

func TestFrameIsResponse(t *testing.T) {
	f := NewResponse([]byte{0x02})
	assert.Equal(t, true, f.IsResponse())
}

func TestDataFrameTypeStringer(t *testing.T) {
	assert.Equal(t, NewDataFrame(Request, []byte{}).String(), "Request{}")
	assert.Equal(t, NewDataFrame(Response, []byte{}).String(), "Response{}")
	assert.Equal(t, NewDataFrame(0x10, []byte{}).String(), "DataFrameType(16){}")
	assert.Equal(t, NewACK().String(), "ACK{}")
}

func BenchmarkChecksumLongFrame(b *testing.B) {
	frame := NewRequest([]byte{
		0x15,
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
