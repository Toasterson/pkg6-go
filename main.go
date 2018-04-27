package main

import (
	"github.com/toasterson/pkg6-go/depotd"
)

func main() {
	depot := depotd.NewDepotServer()
	/*
		rpcSRV, lerr := depot.HandleRPC(viper.GetString("socket_path"))
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
	*/
	mirror := depotd.NewMirrorConfig("hipster", "http://pkg.openindiana.org/hipster")
	//depot.AddMirror(mirror)
	err := depot.Mirror(mirror.Name)
	if err != nil {
		panic(err)
	}
}
