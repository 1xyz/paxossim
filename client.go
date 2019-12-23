package paxossim

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	*Process                // Incoming process and mailbox
	r            []*Replica // Replicas associated
	interval     time.Duration
	done         chan bool
	commandCount int
	ticker       *time.Ticker
}

func NewClient(clientID string, replicas []*Replica, interval time.Duration) *Client {
	return &Client{
		Process:      NewProcess(clientID),
		r:            replicas,
		interval:     interval,
		done:         make(chan bool),
		commandCount: 1,
		ticker:       nil,
	}
}

func (c *Client) nextCommandID() string {
	result := fmt.Sprintf("%d", c.commandCount)
	c.commandCount++
	return result
}

func (c *Client) Run() {
	c.ticker = time.NewTicker(c.interval)
	ctxLog := log.WithFields(log.Fields{
		"id": c.pid,
	})

	for {
		select {
		case <-c.done:
			ctxLog.Debug("done recvd")
			return
		case t := <-c.ticker.C:
			for _, r := range c.r {
				rm := NewRequestMessage(
					c.pid,
					&BasicCommand{
						ClientID:  c.pid,
						CommandID: c.nextCommandID(),
						Op:        "OP",
					})
				r.SendMessage(rm)
				ctxLog.Debugf("message sent to replica %v at %v ", r, t)
			}
		}
	}
}

func (c *Client) Stop() {
	if c.ticker == nil {
		return
	}
	c.ticker.Stop()
	c.done <- true
}
