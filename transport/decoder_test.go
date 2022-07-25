package transport

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DecodeFrameFixture struct {
	FrameType
}

type DecodeFrameCase struct {
	String string
	Frame  *DecodeFrameFixture
}

type DecodeDataFrameCase struct {
	String string
	Frame  *DecodeDataFrameFixture
}

type DecodeDataFrameFixture struct {
	FrameType
	Len int
	DataFrameType
	Payload  []byte
	Checksum byte
}

func TestOneByteFrameDecoding(t *testing.T) {
	cases := map[string]*DecodeFrameCase{
		"ACK": {
			"\x06",
			&DecodeFrameFixture{ACK},
		},
		"NAK": {
			"\x15",
			&DecodeFrameFixture{NAK},
		},
		"CAN": {
			"\x18",
			&DecodeFrameFixture{CAN},
		},
	}
	for title, testCase := range cases {
		t.Run(title, func(t *testing.T) {
			d := NewDecoder(strings.NewReader(testCase.String))
			got, err := d.Next()
			if err != nil {
				t.Fatal("could not load location")
			}
			fixture := testCase.Frame
			assert.Equal(t, fixture.FrameType, got.FrameType)
		})
	}
}

func TestDataFrameDecoding(t *testing.T) {
	cases := map[string]*DecodeDataFrameCase{
		"GetLibraryVersion": {
			"\x01\x10\x01\x15Z-Wave 6.07\x00\x01\x97",
			&DecodeDataFrameFixture{SOF, 16, Response, []byte("\x15Z-Wave 6.07\x00\x01"), byte(0x97)},
		},
	}
	for title, testCase := range cases {
		t.Run(title, func(t *testing.T) {
			d := NewDecoder(strings.NewReader(testCase.String))
			got, err := d.Next()
			if err != nil {
				t.Fatal(err)
			}
			fixture := testCase.Frame
			assert.Equal(t, fixture.FrameType, got.FrameType)
			assert.Equal(t, fixture.Len, got.Length())
			assert.Equal(t, fixture.DataFrameType, *got.DataFrameType)
			assert.Equal(t, fixture.Checksum, got.Checksum())
		})
	}
}

func TestBadStartOfFrame(t *testing.T) {
	d := NewDecoder(strings.NewReader("\x09"))
	_, err := d.Next()
	assert.ErrorContains(t, err, "unrecognized frame type")
}

func TestBadDataFrame(t *testing.T) {
	d := NewDecoder(strings.NewReader("\x01"))
	_, err := d.Next()
	assert.True(t, err == io.EOF)
}

func TestBadChecksum(t *testing.T) {
	d := NewDecoder(strings.NewReader("\x01\x10\x01\x15Z-Wave 6.07\x00\x01\x92"))
	_, err := d.Next()
	assert.ErrorContains(t, err, "checksum did not match")
}
