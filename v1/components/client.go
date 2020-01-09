package components

import (
	"fmt"
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	log "github.com/sirupsen/logrus"
	"time"
)

var clientCount = 0

type Client struct {
	v1.Process

	exchange v1.MessageExchange

	interval time.Duration

	done chan bool

	commandCount int

	ticker *time.Ticker
}

func NewClient(exchange v1.MessageExchange, interval time.Duration) *Client {
	processId := v1.ProcessID(clientCount)
	clientCount++

	return &Client{
		Process:      v1.NewProcess(processId, v1.Client),
		exchange:     exchange,
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
	ctxLog := log.WithFields(log.Fields{"id": c.GetAddr()})

	for {
		select {
		case <-c.done:
			ctxLog.Debug("done recvd")
			return

		case <-c.ticker.C:
			requestMessage := messages.NewRequestMessage(c.GetAddr(), types.BasicCommand{
				ClientID:  fmt.Sprintf("%v", c.GetAddr()),
				CommandID: c.nextCommandID(),
				Op:        "OP",
			})
			c.exchange.SendAll(v1.Replica, requestMessage)
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
