package controller

import (
	"context"
	"encoding"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jbielick/zwgo/hostapi/capabilities/v0"
	"github.com/jbielick/zwgo/transport"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tarm/serial"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func waitForFile(ctx context.Context, path string) {
	for {
		_, err := os.Stat(path)
		if err == nil {
			return
		}
		if os.IsNotExist(err) {
			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
}

type StubbedExchange struct {
	Request   encoding.BinaryMarshaler
	Responses []encoding.BinaryMarshaler
}

func stubbedServer(t *testing.T, f func(Config, chan StubbedExchange)) {
	responses := make(chan StubbedExchange, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	i := rand.Intn(50) + 100
	linkServer := fmt.Sprintf("/tmp/zwgos.%d", i)
	linkClient := fmt.Sprintf("/tmp/zwgoc.%d", i)
	cmd := exec.CommandContext(
		ctx,
		"socat",
		"-d",
		"-d",
		fmt.Sprintf("pty,link=%s,raw,echo=0", linkServer),
		fmt.Sprintf("pty,link=%s,raw,echo=0", linkClient),
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	ready := make(chan bool)
	go func() {
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		waitForFile(ctx, linkServer)
		waitForFile(ctx, linkClient)

		serverSocket, err := os.OpenFile(linkServer, os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		ready <- true
		decoder := transport.NewDecoder(serverSocket)
		for {
			frame, err := decoder.Next()
			if err != nil {
				t.Error(err)
				return
			}
			if frame.IsACK() {
				continue
			} else if frame.IsDataFrame() {
				resp := <-responses
				reqBytes, err := resp.Request.MarshalBinary()
				if err != nil {
					t.Error(err)
					return
				}
				assert.EqualValues(t, reqBytes, frame.Payload)
				for _, payload := range resp.Responses {
					payloadBytes, err := payload.MarshalBinary()
					if err != nil {
						t.Error(err)
						return
					}
					resp := transport.NewResponse(payloadBytes)
					respBytes, err := resp.MarshalBinary()
					if err != nil {
						t.Error(err)
						return
					}
					_, err = serverSocket.Write(respBytes)
					if err != nil {
						t.Error(err)
						return
					}
				}
			} else {
				t.Errorf("received unexpected frame: %v", frame)
			}
		}
	}()
	<-ready
	f(NewConfig(linkClient), responses)
	cmd.Process.Kill()
}

func TestNewConfig(t *testing.T) {
	config := NewConfig("/dev/test")
	assert.Equal(t, config.Serial.Name, "/dev/test")
	assert.Equal(t, config.Serial.Baud, 115200)
	assert.Equal(t, config.Serial.Size, uint8(8))
	assert.Equal(t, config.Serial.Parity, serial.ParityNone)
	assert.Equal(t, config.Serial.ReadTimeout, 10*time.Second)
}

func TestControllerOpen(t *testing.T) {
	stubbedServer(t, func(config Config, responses chan StubbedExchange) {
		c := New(config)
		responses <- StubbedExchange{
			Request: capabilities.NewLibraryVersionGet(),
			Responses: []encoding.BinaryMarshaler{
				transport.NewACK(),
				capabilities.LibraryVersionReport{
					Version:     "Z-Wave 6.07\x00",
					LibraryType: capabilities.LibraryType(capabilities.StaticController),
				},
			},
		}
		responses <- StubbedExchange{
			Request: capabilities.NewGet(),
			Responses: []encoding.BinaryMarshaler{
				transport.NewACK(),
				capabilities.Report{},
			},
		}
		err := c.Open()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, c.LibraryVersion.Version, "Z-Wave 6.07\x00")
		assert.Equal(t, c.Capabilities.Version, 0x01)
	})
}
