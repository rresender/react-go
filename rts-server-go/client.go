package main

import (
	"log"
	"sync"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

// FindHandler to find a handler
type FindHandler func(string) (Handler, bool)

// Client Object
type Client struct {
	send         chan Message
	socket       *websocket.Conn
	findHandler  FindHandler
	session      *r.Session
	stopChannels sync.Map
	id           string
	userName     string
}

// NewStopChannel stop channel
func (c *Client) NewStopChannel(stopKey int) chan bool {
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels.Store(stopKey, stop)
	return stop
}

// StopForKey stops a channel based on a key
func (c *Client) StopForKey(key int) {
	val, ok := c.stopChannels.Load(key)
	if !ok {
		return
	}
	ch, ok := val.(chan bool)
	if !ok {
		return
	}
	ch <- true
	c.stopChannels.Delete(key)
}

// Read Messages from clisnts
func (c *Client) Read() {
	var message Message
	for {
		if err := c.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := c.findHandler(message.Name); found {
			handler(c, message.Data)
		}
	}
	c.socket.Close()
}

func (c *Client) Write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

// Close closed all the resources
func (c *Client) Close() {
	c.stopChannels.Range(func(k, v interface{}) bool {
		v.(chan bool) <- true
		return true
	})
	close(c.send)
	// delete user
	r.Table("user").Get(c.id).Delete().Exec(c.session)
}

// NewClient - create a new client
func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client {
	var user User
	user.Name = "anonymous"
	result, err := r.Table("user").Insert(user).RunWrite(session)
	if err != nil {
		log.Println(err.Error())
	}
	var id string
	if len(result.GeneratedKeys) > 0 {
		id = result.GeneratedKeys[0]
	}
	return &Client{
		send:        make(chan Message),
		socket:      socket,
		findHandler: findHandler,
		session:     session,
		id:          id,
		userName:    user.Name,
	}
}
