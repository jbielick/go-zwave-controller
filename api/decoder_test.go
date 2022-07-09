package api

import (
	"strings"
	"testing"

	. "github.com/stretchr/testify/assert"
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
	CommandType
	CommandID
	Payload  []byte
	Checksum byte
}

func TestOneByteFrameDecoding(t *testing.T) {
	cases := map[string]*DecodeFrameCase{
		"ACK": {
			"\x06",
			&DecodeFrameFixture{ACK},
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
			Equal(t, fixture.FrameType, got.FrameType)
		})
	}
}

func TestDataFrameDecoding(t *testing.T) {
	cases := map[string]*DecodeDataFrameCase{
		"GetLibraryVersion": {
			"\x01\x10\x01\x15Z-Wave 6.07\x00\x01\x97",
			&DecodeDataFrameFixture{SOF, 16, Response, GetLibraryVersion, []byte("Z-Wave 6.07\x00\x01"), byte(0x97)},
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
			Equal(t, fixture.FrameType, got.FrameType)
			Equal(t, fixture.Len, got.Length())
			Equal(t, fixture.CommandType, *got.CommandType)
			Equal(t, fixture.CommandID, *got.CommandID)
		})
	}
}
