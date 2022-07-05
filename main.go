package main

import (
	"log"

	"flag"
	"path/filepath"
	"runtime"

	"github.com/tarm/serial"
)

var deviceFlag = flag.String("device", "", "Name of the zwave serial device (/dev/ttyACM0 for example). Use of this flag disables autodiscovery.")

func main() {
	// v := C.ZW_CHIMNEY_FAN_MIN_SPEED_SET_FRAME{}
	var device string

	flag.Parse()

	if len(*deviceFlag) > 0 {
		device = *deviceFlag
	} else {
		devices, err := findDevice()
		if err != nil {
			log.Fatal(err)
		}
		if len(devices) < 1 {
			log.Fatal("Failed to find device, please provide -device flag to specify the device explicitly.")
		}
		device = devices[0]
	}

	c := &serial.Config{Name: device, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	// n, err := s.Write([]byte("test"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// buf := make([]byte, 128)
	// n, err = s.Read(buf)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Print("%q", buf[:n])

	s.Close()
	// log.Printf("DEBUG: %+v\n", v)
}

func findDevice() ([]string, error) {
	var glob string
	switch runtime.GOOS {
	case "darwin":
		glob = "/dev/tty.usbmodem*"
	default:
		log.Fatalf("must provide device argument for platform %s", runtime.GOOS)
	}
	return filepath.Glob(glob)
}
