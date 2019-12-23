package main

import (
	"github.com/1xyz/paxossim/env"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	NReplicas = 3
	NClients = 2
	NLeaders = 3
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	e := env.NewEnv(NReplicas, NClients, NLeaders)
	log.WithFields(log.Fields{"e": e,}).Debug("Constructed environment")
	e.Run()
	time.Sleep(10 * time.Second)
	e.Stop()
	time.Sleep(1000 * time.Second)
}