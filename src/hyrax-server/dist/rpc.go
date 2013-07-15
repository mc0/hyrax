package dist

import (
    "net"
    "net/rpc"
    "hyrax-server/config"
)

var distCh = make(chan *RpcPayload,1024)
type Dispatcher byte

// Give is where all data from other nodes. Give simply gives the payload
// to this node to be dealt with
func (d *Dispatcher) Receive(pay RpcPayload, r *byte) error {
    distCh <- &pay
    return nil
}

// Returns the list of nodes connected to this one
func (d *Dispatcher) GetNodes(_ interface{}, r *[]string) error {
    *r = getNodes()
    return nil
}

// Tells this node to connect to remote node
func (d *Dispatcher) AddNode(node string, r *byte) error {
    _,err := connectToNode(node)
    return err
}

// Tells this node to connect to remote node, and then tell the remote
// node to connect back
func (d *Dispatcher) AddNodeTwoWay(node string, r *byte) error {
    client,err := connectToNode(node)
    if err != nil {
        return err
    }

    return client.Call("Dispatcher.AddNode",config.GetStr("dist-addr"),nil)
}

// Setup the dist listener and connect to it, adding it to the pool
func Setup() error {
    disp := new(Dispatcher)
    rpc.Register(disp)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", config.GetStr("dist-listen"))
    if err != nil {
        return err
    }
    go rpc.Accept(l)

    _,err = connectToNode(config.GetStr("dist-addr"))
    return err
}

func init() {
    for i:=0; i<10; i++ {
        go func(){
            for pay := range distCh {
                processMessage(pay)
            }
        }()
    }
}
