// Code generated by "stringer -type=state"; DO NOT EDIT.

package transport

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[stateBeginFrame-0]
	_ = x[stateLength-1]
	_ = x[stateDataFrameType-2]
	_ = x[stateStartPayload-3]
	_ = x[stateEndPayload-4]
}

const _state_name = "stateBeginFramestateLengthstateDataFrameTypestateStartPayloadstateEndPayload"

var _state_index = [...]uint8{0, 15, 26, 44, 61, 76}

func (i state) String() string {
	if i < 0 || i >= state(len(_state_index)-1) {
		return "state(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _state_name[_state_index[i]:_state_index[i+1]]
}
