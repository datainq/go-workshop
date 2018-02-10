package topic

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Post struct {
	Author  string
	Message string
}

type Socket interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	io.Closer
}

type Client struct {
	CloseCallback func()
	ToTopic       chan<- Post
	FromTopic     <-chan Post
	Socket        Socket
}

func (c *Client) ReadPipe() {
	defer c.Socket.Close()
	defer c.CloseCallback()
	for {
		var msg Post
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			logrus.Errorf("reading: %s", err)
			break
		}
		c.ToTopic <- msg
		logrus.Infof("%#v", msg)
	}
}

func NewClient(s Socket, t Topic) *Client {
	b := t.GetPipe()
	l, canc := t.AddListener()
	return &Client{
		Socket:        s,
		CloseCallback: canc,
		ToTopic:       b,
		FromTopic:     l,
	}
}

func (c *Client) WritePipe() {
	defer c.Socket.Close()
	for m := range c.FromTopic {
		if err := c.Socket.WriteJSON(m); err != nil {
			break
		}
	}
}
