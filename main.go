package main

import (
	"github.com/spf13/viper"
	"github.com/toasterson/pkg6-go/depotd"
	"os"
	"os/signal"
	"syscall"
)

var repoPath = "file://./sample_data/repo"

func init() {
	viper.SetDefault("sockerPath", "/var/run/depotd.sock")
}

func main() {
	depot := depotd.NewDepotServer(repoPath)
	if err := depot.Load(); err != nil {
		panic(err)
	}
	rpcSRV, lerr := depot.HandleRPC(viper.GetString("socketPath"))
	if lerr != nil {
		panic(lerr)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL:
		sig := <-c
		depot.Logger.Printf("Caught signal %s: shutting down.", sig)
		// Stop listening (and unlink the socket if unix type):
		rpcSRV.Socket.Close()
		// And we're done:
		os.Exit(0)
	}(sigc)
	depot.Start(":8080")
}
