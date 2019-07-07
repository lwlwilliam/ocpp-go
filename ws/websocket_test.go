package ws

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

const (
	serverPort = 8887
	serverPath = "/ws/{id}"
	testPath   = "/ws/testws"
)

func NewWebsocketServer(t *testing.T, onMessage func(data []byte) ([]byte, error)) *Server {
	wsServer := Server{}
	wsServer.SetMessageHandler(func(ws Channel, data []byte) error {
		assert.NotNil(t, ws)
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = wsServer.Write(ws.GetId(), data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return &wsServer
}

func NewWebsocketClient(t *testing.T, onMessage func(data []byte) ([]byte, error)) *Client {
	wsClient := Client{}
	wsClient.SetMessageHandler(func(data []byte) error {
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = wsClient.Write(data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return &wsClient
}

func TestWebsocketEcho(t *testing.T) {
	message := []byte("Hello WebSocket!")
	var wsServer *Server
	wsServer = NewWebsocketServer(t, func(data []byte) ([]byte, error) {
		assert.True(t, bytes.Equal(message, data))
		return data, nil
	})
	go wsServer.Start(serverPort, serverPath)

	// Test message
	wsClient := NewWebsocketClient(t, func(data []byte) ([]byte, error) {
		assert.True(t, bytes.Equal(message, data))
		return nil, nil
	})
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Wait for connection to be established, then send a message
	go func() {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		err := wsClient.Write(message)
		assert.Nil(t, err)
	}()
	done := make(chan bool)
	// Wait for messages to be exchanged, then close connection
	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		wsClient.Stop()
		done <- true
	}()
	err := wsClient.Start(u.String())
	assert.Nil(t, err)
	result := <-done
	assert.True(t, result)
}