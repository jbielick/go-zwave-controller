package transport

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

//go:generate stringer -type=state
type state int

const (
	stateBeginFrame state = iota
	stateLength
	stateDataFrameType
	stateStartPayload
	stateEndPayload
)

type Decoder struct {
	r         io.Reader
	state     state
	buf       []byte
	pos       int
	have      int
	dataFrame *Frame
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, buf: make([]byte, 256)}
}

func (d *Decoder) More() error {
	nBytes, err := d.r.Read(d.buf)
	if err != nil {
		return err
	}
	if nBytes < 1 {
		log.Warn("got 0 bytes when reading")
	}
	d.have = nBytes
	return nil
}

// func (d *Decoder) Peek() (byte, error) {
// 	return d.buf[d.pos+1], nil
// }

func (d *Decoder) NextByte() (byte, error) {
	if d.pos >= d.have {
		err := d.More()
		if err != nil {
			return 0x00, err
		}
		d.pos = 0
	}
	b := d.buf[d.pos]
	d.pos++
	return b, nil
}

func (d *Decoder) Next() (*Frame, error) {
	b, err := d.NextByte()
	if err != nil {
		return nil, err
	}

	switch d.state {
	case stateBeginFrame:
		return d.stateBeginFrame(b)
	case stateLength:
		return d.stateLength(b)
	case stateDataFrameType:
		return d.stateDataFrameType(b)
	case stateStartPayload:
		return d.stateStartPayload(b)
	case stateEndPayload:
		return d.stateEndPayload(b)
	default:
		return nil, fmt.Errorf(
			"decoder in unexpected state while decoding frame: state=%s frame=%s byte=%q",
			d.state,
			d.dataFrame,
			b,
		)
	}
}

func (d *Decoder) stateBeginFrame(b byte) (*Frame, error) {
	switch FrameType(b) {
	case ACK:
		return NewACK(), nil
	case NAK:
		return NewNAK(), nil
	case CAN:
		return NewCAN(), nil
	case SOF:
		d.dataFrame = &Frame{FrameType: SOF}
		d.state = stateLength
		return d.Next()
	default:
		return nil, fmt.Errorf("unrecognized frame type: %q", b)
	}
}

func (d *Decoder) stateLength(b byte) (*Frame, error) {
	d.dataFrame.Len = int(b)
	d.state = stateDataFrameType
	return d.Next()
}

func (d *Decoder) stateDataFrameType(b byte) (*Frame, error) {
	d.dataFrame.SetDataFrameType(DataFrameType(b))
	d.state = stateStartPayload
	return d.Next()
}

func (d *Decoder) stateStartPayload(b byte) (*Frame, error) {
	d.dataFrame.Payload = append(d.dataFrame.Payload, b)
	if len(d.dataFrame.Payload) == d.dataFrame.Len-2 {
		d.state = stateEndPayload
	}
	return d.Next()
}

func (d *Decoder) stateEndPayload(statedSum byte) (*Frame, error) {
	calculatedSum := d.dataFrame.Checksum()
	if statedSum != calculatedSum {
		return d.dataFrame, fmt.Errorf(
			"checksum did not match. have: %+q, want: %+q",
			calculatedSum,
			statedSum,
		)
	}
	d.state = stateBeginFrame
	dataFrame := d.dataFrame
	d.dataFrame = nil
	log.Printf("%#v %q", dataFrame, dataFrame.Checksum())

	return dataFrame, nil
}
