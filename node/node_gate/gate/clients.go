package gate

import (
	"github.com/yddeng/dnet"
	"sync"
)

type Client struct {
	UserID  string
	session dnet.Session
}

var clients sync.Map //map[string]*Client

func addClient(c *Client) {
	clients.Store(c.UserID, c)
}

func removeClient(userID string) {
	clients.Delete(userID)
}
