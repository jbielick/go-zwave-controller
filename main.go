package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/jbielick/zwgo/controller"
	log "github.com/sirupsen/logrus"
)

var portFlag = flag.String("port", "", "Name of the zwave serial port (/dev/ttyACM0 for example). Use of this flag disables autodiscovery.")
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
	var port string

	if len(*portFlag) > 0 {
		port = *portFlag
	} else {
		log.Info("Discovering port...")
		ports, err := discoverPorts()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) < 1 {
			log.Fatal("Failed to find port, please provide -port flag to specify the port explicitly.")
		}
		port = ports[0]
	}
	log.Infof("Using port %s", port)

	c := controller.New(controller.NewConfig(port))
	err := c.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	log.Printf("DEBUG: %+v\n", c)

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-s
		signal.Stop(s)
		cancel()
	}()

	<-ctx.Done()
}

func discoverPorts() ([]string, error) {
	var glob string
	switch runtime.GOOS {
	case "darwin":
		glob = "/dev/tty.usbmodem*"
	case "linux":
		glob = "/dev/ttyACM*"
	default:
		log.Fatalf("must provide port argument for platform %s", runtime.GOOS)
	}
	return filepath.Glob(glob)
}
