// Code generated by "stringer -type=FrameType"; DO NOT EDIT.

package transport

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SOF-1]
	_ = x[ACK-6]
	_ = x[NAK-21]
	_ = x[CAN-24]
}

const (
	_FrameType_name_0 = "SOF"
	_FrameType_name_1 = "ACK"
	_FrameType_name_2 = "NAK"
	_FrameType_name_3 = "CAN"
)

func (i FrameType) String() string {
	switch {
	case i == 1:
		return _FrameType_name_0
	case i == 6:
		return _FrameType_name_1
	case i == 21:
		return _FrameType_name_2
	case i == 24:
		return _FrameType_name_3
	default:
		return "FrameType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}