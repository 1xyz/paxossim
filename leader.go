package paxossim

import (
	log "github.com/sirupsen/logrus"
)

type Leader struct {
	*Process
}

func NewLeader(leaderID string) *Leader {
	return &Leader{
		Process: NewProcess(leaderID),
	}
}

func (l *Leader) Run() {
	ctxLog := log.WithFields(log.Fields{
		"id": l.pid,
	})
	for {
		msg := l.inbox.WaitForItem()
		ctxLog.Debugf("Got a new message %v", msg)
	}
}
