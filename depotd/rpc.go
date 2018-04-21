package depotd

import (
	"net"
	"net/rpc"
	"strings"
)

type RPCDepot struct {
	Depot      *DepotServer
	Socket     net.Listener
	SockerPath string
	run        bool
}

func (d *DepotServer) HandleRPC(sock string) (r *RPCDepot, err error) {
	r = &RPCDepot{
		Depot:      d,
		SockerPath: sock,
		run:        true,
	}
	socket, err := net.Listen("unix", sock)
	if err != nil {
		return nil, err
	}
	r.Socket = socket
	rpcSrv := rpc.NewServer()
	rpcSrv.Register(r)
	go func() {
		for r.run {
			conn, err := r.Socket.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					r.run = false
				} else {
					d.Logger.Printf("rpc.Serve: accept: %s", err.Error())
				}
				return
			}
			go rpcSrv.ServeConn(conn)
		}
	}()
	return r, nil
}

func (r *RPCDepot) Ping(message string, reply *string) error {
	r.Depot.Logger.Print(message)
	*reply = "Pong"
	return nil
}
