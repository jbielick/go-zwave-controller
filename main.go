package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jbielick/zwgo/api"
	log "github.com/sirupsen/logrus"
)

var deviceFlag = flag.String("device", "", "Name of the zwave serial device (/dev/ttyACM0 for example). Use of this flag disables autodiscovery.")
var verbosityFlag = flag.String("log-level", "INFO", "Logging verbosity")

func init() {
	flag.Parse()
	level, err := log.ParseLevel(*verbosityFlag)
	if err != nil {
		log.Fatal(fmt.Errorf("'%s' is not a valid log level, please use one of %v", *verbosityFlag, log.AllLevels))
	}
	log.SetLevel(level)
}

func main() {
	// v := C.ZW_CHIMNEY_FAN_MIN_SPEED_SET_FRAME{}
	var device string

	if len(*deviceFlag) > 0 {
		device = *deviceFlag
	} else {
		log.Info("Discovering device...")
		devices, err := discoverDevices()
		if err != nil {
			log.Fatal(err)
		}
		if len(devices) < 1 {
			log.Fatal("Failed to find device, please provide -device flag to specify the device explicitly.")
		}
		device = devices[0]
	}
	log.Infof("Using device %s", device)

	c := api.NewController(device)
	err := c.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for {
		r, err := c.GetLibraryVersion()
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("Controller Library Version: %+v", r)
		time.Sleep(1 * time.Second)
	}

	// frame2, err := c.SendAndReceive(api.NewRequest(api.GET_PROTOCOL_VERSION, []byte{}))
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// log.Infof("Controller Protocol Version: %s", frame2.PayloadString())

	// frame3, err := c.SendAndReceive(api.NewRequest(api.GET_INIT_DATA, []byte{}))
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// r := new(api.GetInitDataResponse)
	// r.UnmarshalBinary(frame3.Payload)
	// log.Infof("Init Data: %s", r.APIVersion.String())
}

func discoverDevices() ([]string, error) {
	var glob string
	switch runtime.GOOS {
	case "darwin":
		glob = "/dev/tty.usbmodem*"
	case "linux":
		glob = "/dev/ttyACM*"
	default:
		log.Fatalf("must provide device argument for platform %s", runtime.GOOS)
	}
	return filepath.Glob(glob)
}
