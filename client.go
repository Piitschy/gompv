package gompv

import (
	"fmt"

	"github.com/dexterlb/mpvipc"
)

const defaultSocketPath = "/tmp/mpv_socket"

// MPVClient extends the mpvipc.Connection to provide a client for interacting with MPV.
// You can use all methods from mpvipc.Connection, and it adds some convenience methods.
type MPVClient struct {
	*mpvipc.Connection
	socketPath string
}

func NewMPVClient() *MPVClient {
	client := &MPVClient{socketPath: defaultSocketPath}
	client.Connection = client.NewConnection()
	return client
}

func NewMPVClientWithSocketPath(socketPath string) *MPVClient {
	client := &MPVClient{socketPath: socketPath}
	client.Connection = client.NewConnection()
	return client
}

func NewMPVClientWithConnection(conn *mpvipc.Connection) *MPVClient {
	client := &MPVClient{Connection: conn}
	client.socketPath = ""
	return client
}

func (c *MPVClient) NewConnection() *mpvipc.Connection {
	return mpvipc.NewConnection(c.socketPath)
}

func (c MPVClient) SocketPath() string {
	if c.socketPath == "" {
		return defaultSocketPath
	}
	return c.socketPath
}

func (c *MPVClient) GetPath() (string, error) {
	path, err := c.Get("path")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", path), nil
}

func (c *MPVClient) Pause() error {
	return c.Set("pause", true)
}
