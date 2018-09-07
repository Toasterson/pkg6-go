package depotd

type AddMirrorArgs struct {
	Name       string
	BaseURL    string
	Publishers []string
}

func (r *RPCDepot) AddMirror(args AddMirrorArgs, reply *string) error {
	return nil
}
