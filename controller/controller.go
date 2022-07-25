package controller

import (
	"encoding"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jbielick/zwgo/hostapi/capabilities/v0"
	"github.com/jbielick/zwgo/transport"
	log "github.com/sirupsen/logrus"

	"github.com/tarm/serial"
)

type Config struct {
	Serial serial.Config
}

type Controller struct {
	Config         Config
	Port           *serial.Port
	mu             sync.Mutex
	LibraryVersion capabilities.LibraryVersionReport
	// InitData               capabilities.InitDataReport
	Capabilities capabilities.Report
	// ControllerCapabilities capabilities.ControllerCapabilities
	inbox       chan *transport.Frame
	unsolicited chan *transport.Frame
}

func New(config Config) *Controller {
	return &Controller{
		Config:      config,
		inbox:       make(chan *transport.Frame, 20),
		unsolicited: make(chan *transport.Frame, 20),
	}
}

func NewConfig(port string) Config {
	return Config{
		Serial: serial.Config{
			Name:        port,
			Baud:        115200,
			ReadTimeout: 10 * time.Second,
			Parity:      serial.ParityNone,
			Size:        8,
			StopBits:    1,
		},
	}
}

func (c *Controller) Open() error {
	port, err := serial.OpenPort(&c.Config.Serial)
	if err != nil {
		return err
	}
	c.Port = port
	c.Port.Flush()

	go c.receive()
	go c.handleRequests()

	if err := c.initialize(); err != nil {
		return err
	}

	return nil
}

func (c *Controller) receive() {
	d := transport.NewDecoder(c.Port)
	for {
		frame, err := d.Next()
		if err == io.EOF {
			log.Error(err)
			return
		}
		log.Debugf("← %s", frame)
		if frame.IsRequest() {
			c.unsolicited <- frame
		} else {
			c.inbox <- frame
		}
		// @TODO shutdown channel
	}
}

func (c *Controller) handleRequests() {
	for {
		req := <-c.unsolicited
		log.Printf("handleRequest: %q", req)
	}
}

func (c *Controller) ack() error {
	_, err := c.Send(transport.NewACK())
	return err
}

func (c *Controller) Send(f *transport.Frame) (int, error) {
	bytes, err := f.MarshalBinary()
	if err != nil {
		return 0, err
	}
	log.Debugf("→ %s", f)
	return c.Port.Write(bytes)
}

func (c *Controller) sendWithAcknowledgementUnlocked(cmd encoding.BinaryMarshaler) (int, error) {
	attempts := 0
	cmdBytes, err := cmd.MarshalBinary()
	if err != nil {
		return 0, err
	}

retry:
	sent, err := c.Send(transport.NewRequest(cmdBytes))
	if err != nil {
		return sent, err
	}

	select {
	case frame := <-c.inbox:
		switch frame.FrameType {
		case transport.ACK:
			return sent, nil
		case transport.NAK:
			err = fmt.Errorf("Messsage was not acknowledged by the device")
		case transport.CAN:
			err = fmt.Errorf("Received CAN while waiting for acknowledgement: %s", frame)
			log.Warn(err)
			time.Sleep(1 * time.Second)
		default:
			return sent, fmt.Errorf("Received unexpected frame waiting for acknowledgement: %s", frame)
		}
	case <-time.After(2 * time.Second):
		err = fmt.Errorf("Timed out waiting for acknowledgement")
	}
	if err == nil || attempts >= 2 {
		return sent, err
	} else {
		attempts++
		time.Sleep(100*time.Millisecond + time.Duration(attempts)*time.Second)
		goto retry
	}
}

func (c *Controller) SendWithAcknowledgement(cmd encoding.BinaryMarshaler) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sendWithAcknowledgementUnlocked(cmd)
}

func (c *Controller) SendAndReceive(cmd encoding.BinaryMarshaler, v encoding.BinaryUnmarshaler) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.sendWithAcknowledgementUnlocked(cmd)
	if err != nil {
		return err
	}
	select {
	case resp := <-c.inbox:
		err = c.ack()
		if err != nil {
			return err
		}
		return v.UnmarshalBinary(resp.Payload)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timed out waiting for response")
	}
}

func (c *Controller) SendWithCallback(cmd encoding.BinaryMarshaler, callback func(v encoding.BinaryUnmarshaler)) error {
	_, err := c.SendWithAcknowledgement(cmd)
	if err != nil {
		return err
	}
	return err
}

func (c *Controller) initialize() error {
	log.Info("Retrieving initialization data...")

	libraryReport, err := capabilities.NewLibraryVersionGet().Send(c)
	if err != nil {
		return err
	}
	c.LibraryVersion = libraryReport

	capabilitiesReport, err := capabilities.NewGet().Send(c)
	if err != nil {
		return err
	}
	c.Capabilities = capabilitiesReport

	// req2 := capabilities.NewGetInitData()
	// res2 := capabilities.InitData{}
	// if err := c.SendAndReceive(req2, &res2); err != nil {
	// 	return err
	// }
	// c.InitData = res2

	// capabilities, err := capabilities.NewGetCapabilities().Send(c)
	// if err != nil {
	// 	return err
	// }
	// c.Capabilities = capabilities

	// req4 := capabilities.NewGetControllerCapabilities()
	// res4 := capabilities.ControllerCapabilities{}
	// if err := c.SendAndReceive(req4, &res4); err != nil {
	// 	return err
	// }
	// c.ControllerCapabilities = res4
	return nil
}

func (c *Controller) Close() error {
	c.Port.Flush()
	return c.Port.Close()
}
