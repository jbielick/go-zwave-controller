// Code generated by "stringer -type=BasicDeviceClass"; DO NOT EDIT.

package api

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BASIC_TYPE_CONTROLLER-1]
	_ = x[BASIC_TYPE_STATIC_CONTROLLER-2]
	_ = x[BASIC_TYPE_END_NODE-3]
	_ = x[BASIC_TYPE_ROUTING_END_NODE-4]
}

const _BasicDeviceClass_name = "BASIC_TYPE_CONTROLLERBASIC_TYPE_STATIC_CONTROLLERBASIC_TYPE_END_NODEBASIC_TYPE_ROUTING_END_NODE"

var _BasicDeviceClass_index = [...]uint8{0, 21, 49, 68, 95}

func (i BasicDeviceClass) String() string {
	i -= 1
	if i >= BasicDeviceClass(len(_BasicDeviceClass_index)-1) {
		return "BasicDeviceClass(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _BasicDeviceClass_name[_BasicDeviceClass_index[i]:_BasicDeviceClass_index[i+1]]
}
