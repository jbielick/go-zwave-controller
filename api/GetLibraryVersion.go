// STOP
// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.

package api

const GetLibraryVersion CommandID = 0x15

// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.
//go:generate stringer -type=LibraryType
type LibraryType byte


// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.
const (
	StaticControllerLibrary LibraryType = 0x01
	PortableControllerLibrary LibraryType = 0x02
	Enhanced232EndNodeLibrary LibraryType = 0x03
	EndNodeLibrary LibraryType = 0x04
	InstallerLibrary LibraryType = 0x05
	RoutingEndNodeLibrary LibraryType = 0x06
	BridgeControllerLibrary LibraryType = 0x07
)

// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.
type GetLibraryVersionResponse struct {
	LibraryVersion string
	LibraryType LibraryType
}

func (r *GetLibraryVersionResponse) UnmarshalBinary(data []byte) error {
	pos := 0
  r.LibraryVersion = string(data[pos:pos+12])
  pos += 12
	r.LibraryType = LibraryType(data[pos])
	pos++

	return nil
}

func (c *Controller) GetLibraryVersion() (*GetLibraryVersionResponse, error) {
  r := &GetLibraryVersionResponse{}
	frame, err := c.SendAndReceive(NewRequest(GetLibraryVersion, []byte{}))
  if err != nil {
    return r, err
  }
	r.UnmarshalBinary(frame.Payload)
	return r, nil
}
