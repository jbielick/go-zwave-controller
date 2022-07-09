package api

import (
	"fmt"
	"log"
)

const (
	GET_INIT_DATA                                    CommandID = 0x02
	SET_APPLICATION_NODE_INFORMATION                 CommandID = 0x03
	SET_APPLICATION_NODE_INFORMATION_COMMAND_CLASSES CommandID = 0x0C // END NODE only
	GET_CONTROLLER_CAPABILITIES                      CommandID = 0x05
	GET_CAPABILITIES                                 CommandID = 0x07
	GET_LONG_RANGE_NODES                             CommandID = 0xDA
	GET_ZWAVE_LONG_RANGE_CHANNEL                     CommandID = 0xDB
	SET_ZWAVE_LONG_RANGE_CHANNEL                     CommandID = 0xDC
	GET_PROTOCOL_VERSION                             CommandID = 0x09

	GET_LIBRARY     CommandID = 0xBD
	SOFT_RESET      CommandID = 0x08
	SET_DEFAULT     CommandID = 0x42
	SETUP_ZWAVE_API CommandID = 0x0B // has subcommands, 4.3.15
)

type APIVersion byte

// func (a *APIVersion) String() string {
// 	return strconv.Itoa(int(*a) - 9)
// }

type APICapabilities byte

func (a *APICapabilities) String() string {
	var nodeType string
	if *a&0 == 0 {
		nodeType = "EndNode"
	} else if *a&1 == 1 {
		nodeType = "Controller"
	} else {
		nodeType = "???"
	}
	return string(fmt.Sprintf("%s", nodeType))
}

type Node struct {
	ID int
}

type GetInitDataResponse struct {
	APIVersion
	APICapabilities
	Nodes []Node
}

func (r *GetInitDataResponse) UnmarshalBinary(data []byte) error {
	pos := 0
	log.Printf("%d", data[pos])
	r.APIVersion = APIVersion(data[pos])
	pos++
	r.APICapabilities = APICapabilities(data[pos])
	pos++
	nodeListLen := int(data[pos])
	pos++
	for i := 0; i < nodeListLen; i++ {
		node := data[pos]
		r.Nodes = append(r.Nodes, Node{int(node)})
		pos++
	}

	return nil
}
