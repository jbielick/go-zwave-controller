package transport

import (
	"encoding"
	"fmt"
	"log"
)

//go:generate stringer -type=FrameType
type FrameType byte

const (
	SOF FrameType = 0x01
	ACK FrameType = 0x06
	NAK FrameType = 0x15
	CAN FrameType = 0x18
)

//go:generate stringer -type=DataFrameType
type DataFrameType byte

const (
	Request  DataFrameType = 0x00
	Response DataFrameType = 0x01
)

type Frame struct {
	FrameType     FrameType
	Len           int
	DataFrameType *DataFrameType
	Payload       encoding.BinaryMarshaler
}

func (f *Frame) Length() int {
	// length (1 byte) + type (1 byte) + command payload
	length := 1
	if f.IsDataFrame() {
		length += 1
		length += len(binary.Marshal(f.Payload))
	}
	return length
}

func (f *Frame) Checksum() byte {
	var sum byte = 0xFF
	sum ^= byte(f.Length())
	sum ^= byte(*f.DataFrameType)
	for i := range f.Payload {
		sum ^= byte(f.Payload[i])
	}
	return sum
}

func (f *Frame) MarshalBinary() ([]byte, error) {
	b := []byte{byte(f.FrameType)}
	if f.IsDataFrame() {
		b = append(b, byte(f.Length()))
		b = append(b, byte(*f.DataFrameType))
		b = append(b, f.Payload...)
		b = append(b, f.Checksum())
	}

	return b, nil
}

func (f *Frame) UnmarshalBinary(data []byte) error {
	f.FrameType = FrameType(data[0])
	if !f.IsDataFrame() {
		return nil
	}
	DataFrameType := DataFrameType(data[2])

	f.DataFrameType = &DataFrameType
	f.Payload = data[3 : len(data)-1]
	log.Print(f)
	log.Print(f.Length())

	checksum := data[len(data)-1]
	calculatedChecksum := f.Checksum()
	if checksum != calculatedChecksum {
		return fmt.Errorf("Failed to unmarshal frame, checksum did not match. Want: %+q, Got: %+q", checksum, calculatedChecksum)
	}

	return nil
}

func (f *Frame) SetDataFrameType(c DataFrameType) {
	f.DataFrameType = &c
}

func (f Frame) IsACK() bool {
	return f.FrameType == ACK
}

func (f Frame) IsNAK() bool {
	return f.FrameType == NAK
}

func (f Frame) IsCAN() bool {
	return f.FrameType == CAN
}

func (f Frame) IsDataFrame() bool {
	return f.FrameType == SOF
}

func (f Frame) IsRequest() bool {
	return f.IsDataFrame() && f.DataFrameType != nil && *f.DataFrameType == Request
}

func (f Frame) IsResponse() bool {
	return f.IsDataFrame() && f.DataFrameType != nil && *f.DataFrameType == Response
}

func (f *Frame) String() string {
	if f.FrameType == SOF {
		return fmt.Sprintf("%s{% x}", f.DataFrameType, f.Payload)
	} else {
		return fmt.Sprintf("%s{}", f.FrameType)
	}
}

func (f *Frame) PayloadString() string {
	return string(f.Payload)
}

func NewRequest(payload []byte) *Frame {
	return NewDataFrame(Request, payload)
}

func NewResponse(payload []byte) *Frame {
	return NewDataFrame(Response, payload)
}

func NewDataFrame(DataFrameType DataFrameType, payload []byte) *Frame {
	return &Frame{
		SOF,
		0,
		&DataFrameType,
		payload,
	}
}

func NewACK() *Frame {
	return &Frame{ACK, 0, nil, nil}
}

func NewNAK() *Frame {
	return &Frame{NAK, 0, nil, nil}
}

func NewCAN() *Frame {
	return &Frame{CAN, 0, nil, nil}
}
