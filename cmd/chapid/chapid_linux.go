package main

// Copyright 2019 Hewlett Packard Enterprise Development LP.
import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hpe-storage/common-host-libs/chapi"
	log "github.com/hpe-storage/common-host-libs/logger"
	"github.com/hpe-storage/common-host-libs/util"
)

var (
	// Version contains the current version added by the build process
	Version = "dev"
	// Commit containers the hg commit added by the build process
	Commit    = "unknown"
	chapidLog = util.GetNltHome() + "log/chapid.log"
)

func main() {
	log.InitLogging(chapidLog, &log.LogParams{Level: "trace"}, false)
	log.Infof("Starting chapi server version %s(%s)...", Version, Commit)
	nimbledChan := make(chan error)
	chapi.RunNimbled(nimbledChan)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		log.Fatalf("Exiting due to signal notification.  Signal was %v.", s.String())
		os.Exit(1)
	}()
	x := <-nimbledChan
	log.Error("error on chapid socket:", x)
}
