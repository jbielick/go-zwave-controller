package api

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tarm/serial"
)

type Controller struct {
	SerialConfig *serial.Config
	Port         *serial.Port
	wg           sync.WaitGroup
	mu           sync.Mutex
	inbox        chan *Frame
	unsolicited  chan *Frame
}

func NewController(device string) *Controller {
	return &Controller{
		SerialConfig: &serial.Config{
			Name:        device,
			Baud:        115200,
			ReadTimeout: 10 * time.Second,
		},
		inbox:       make(chan *Frame, 10),
		unsolicited: make(chan *Frame, 10),
	}
}

func (c *Controller) Open() error {
	port, err := serial.OpenPort(c.SerialConfig)
	if err != nil {
		return err
	}
	c.Port = port
	c.Port.Flush()
	go c.receive()

	return nil
}

func (c *Controller) receive() {
	d := NewDecoder(c.Port)
	for {
		frame, err := d.Next()
		if err != nil {
			log.Warn(err)
		}
		log.Debugf("← %s", frame)
		c.inbox <- frame
	}
}

func (c *Controller) ack() error {
	_, err := c.Send(NewACK())
	return err
}

func (c *Controller) Send(f *Frame) (int, error) {
	bytes, err := f.MarshalBinary()
	if err != nil {
		return 0, err
	}
	log.Debugf("→ %s", f)

	return c.Port.Write(bytes)
}

func (c *Controller) SendWithAcknowledgement(f *Frame, attempts int, lock bool) (int, error) {
	if lock {
		c.mu.Lock()
		defer c.mu.Unlock()
	}
	sent, err := c.Send(f)
	if err != nil {
		return sent, err
	}

	retry := func(err error) (int, error) {
		if attempts > 2 {
			return sent, err
		} else {
			time.Sleep(100*time.Millisecond + time.Duration(attempts)*time.Second)
			return c.SendWithAcknowledgement(f, attempts+1, false)
		}
	}

	select {
	case frame := <-c.inbox:
		switch frame.FrameType {
		case ACK:
			return sent, nil
		case NAK:
			return retry(fmt.Errorf("Messsage was not acknowledged by the device"))
		case CAN:
			time.Sleep(1 * time.Second)
			return retry(fmt.Errorf("Received CAN while waiting for acknowledgement: %s", frame))
		default:
			return sent, fmt.Errorf("Received unexpected frame waiting for acknowledgement: %s", frame)
		}
	case <-time.After(2 * time.Second):
		return retry(fmt.Errorf("Timed out waiting for acknowledgement"))
	}
}

func (c *Controller) SendAndReceive(f *Frame) (*Frame, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.SendWithAcknowledgement(f, 0, false)
	if err != nil {
		return nil, err
	}
	select {
	case resp := <-c.inbox:
		err = c.ack()
		if err != nil {
			return resp, err
		}
		return resp, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timed out waiting for response")
	}
}

// @TODO
// func (c *Controller) SendWithCallback(f *Frame, fn func(*Frame)) error {
// 	_, err := c.SendWithAcknowledgement(f)
// 	return err
// }

func (c *Controller) Close() error {
	c.Port.Flush()
	return c.Port.Close()
}
