package api

//go:generate stringer -type=CommandID
type CommandID byte

const (
	NOP                                CommandID = 0xE9
	GET_NODE_INFORMATION_PROTOCOL_DATA CommandID = 0x41
	SEND_NODE_INFORMATION              CommandID = 0x12
	REQUEST_NODE_INFORMATION           CommandID = 0x60
	SET_LEARN_MODE                     CommandID = 0x50
)

//go:generate stringer -type=BasicDeviceClass
type BasicDeviceClass byte

const (
	BASIC_TYPE_CONTROLLER        BasicDeviceClass = 0x01
	BASIC_TYPE_STATIC_CONTROLLER BasicDeviceClass = 0x02
	BASIC_TYPE_END_NODE          BasicDeviceClass = 0x03
	BASIC_TYPE_ROUTING_END_NODE  BasicDeviceClass = 0x04
)

//go:generate stringer -type=TransmissionOption
type TransmissionOption int

const (
	TRANSMIT_OPTION_ACK        TransmissionOption = 0
	TRANSMIT_OPTION_LOW_POWER  TransmissionOption = 1
	TRANSMIT_OPTION_AUTO_ROUTE TransmissionOption = 2
	TRANSMIT_OPTION_NO_ROUTE   TransmissionOption = 4
	TRANSMIT_OPTION_EXPLORE    TransmissionOption = 5
)
