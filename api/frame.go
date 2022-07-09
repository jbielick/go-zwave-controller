package api

import (
	"fmt"
)

//go:generate stringer -type=FrameType
type FrameType byte

const (
	SOF FrameType = 0x01
	ACK FrameType = 0x06
	NAK FrameType = 0x15
	CAN FrameType = 0x18
)

type CommandType byte

//go:generate stringer -type=CommandType
const (
	Request  CommandType = 0x00
	Response CommandType = 0x01
)

type Frame struct {
	FrameType   FrameType
	Len         int
	CommandType *CommandType
	CommandID   *CommandID
	Payload     []byte
}

func (f *Frame) Length() int {
	// length (1 byte) + type (1 byte) + command ID + command payload
	length := 1
	if f.FrameType == SOF {
		length += 2
	}
	length += len(f.Payload)
	return length
}

func (f *Frame) Checksum() byte {
	var sum byte = 0xFF
	sum ^= byte(f.Length())
	sum ^= byte(*f.CommandType)
	sum ^= byte(*f.CommandID)
	for i := range f.Payload {
		if f.Payload[i] == 0x00 {
			// continue
		}
		sum ^= byte(f.Payload[i])
	}
	return sum
}

func (f *Frame) MarshalBinary() ([]byte, error) {
	b := []byte{byte(f.FrameType)}
	if f.FrameType == SOF {
		b = append(b, byte(f.Length()))
		b = append(b, byte(*f.CommandType))
		b = append(b, byte(*f.CommandID))
		b = append(b, f.Payload...)
		b = append(b, f.Checksum())
	}

	return b, nil
}

// func (f *Frame) UnmarshalBinary(data []byte) error {
// 	f.FrameType = FrameType(data[0])
// 	if !f.IsDataFrame() {
// 		return nil
// 	}
// 	commandType := CommandType(data[2])
// 	commandID := CommandID(data[3])

// 	f.CommandType = &commandType
// 	f.CommandID = &commandID
// 	f.Payload = data[4 : len(data)-1]
// 	log.Print(f)
// 	log.Print(f.Length())

// 	checksum := data[len(data)-1]
// 	calculatedChecksum := f.Checksum()
// 	if checksum != calculatedChecksum {
// 		return fmt.Errorf("Failed to unmarshal frame, checksum did not match. Want: %+q, Got: %+q", checksum, calculatedChecksum)
// 	}

// 	return nil
// }

func (f *Frame) SetCommandType(c CommandType) {
	f.CommandType = &c
}

func (f *Frame) SetCommandID(c CommandID) {
	f.CommandID = &c
}

func (f *Frame) IsACK() bool {
	return f.FrameType == ACK
}

func (f *Frame) IsNAK() bool {
	return f.FrameType == NAK
}

func (f *Frame) IsCAN() bool {
	return f.FrameType == CAN
}

func (f *Frame) IsDataFrame() bool {
	return f.FrameType == SOF
}

func (f *Frame) IsRequest() bool {
	return f.IsDataFrame() && *f.CommandType == Request
}

func (f *Frame) IsResponse() bool {
	return f.IsDataFrame() && *f.CommandType == Response
}

func (f *Frame) String() string {
	if f.FrameType == SOF {
		return fmt.Sprintf("%s{Command=%s Payload=%q}", f.CommandType, f.CommandID, string(f.Payload))
	} else {
		return fmt.Sprintf("%s{}", f.FrameType)
	}
}

func (f *Frame) PayloadString() string {
	return string(f.Payload)
}

func NewRequest(commandID CommandID, payload []byte) *Frame {
	return NewDataFrame(Request, commandID, payload)
}

func NewResponse(commandID CommandID, payload []byte) *Frame {
	return NewDataFrame(Response, commandID, payload)
}

func NewDataFrame(commandType CommandType, commandID CommandID, payload []byte) *Frame {
	return &Frame{
		SOF,
		0,
		&commandType,
		&commandID,
		payload,
	}
}

func NewACK() *Frame {
	return &Frame{ACK, 0, nil, nil, nil}
}

func NewNAK() *Frame {
	return &Frame{NAK, 0, nil, nil, nil}
}

func NewCAN() *Frame {
	return &Frame{CAN, 0, nil, nil, nil}
}
