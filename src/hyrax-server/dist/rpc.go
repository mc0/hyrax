package dist

import (
    "net"
    "net/rpc"
    "hyrax-server/config"
    "errors"
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
func (d *Dispatcher) GetNodes(_ struct{}, r *[]string) error {
    *r = getNodes()
    return nil
}

// Returns the dist-addr string that other nodes should store as this
// nodes' identifier and connect to it on
func (d *Dispatcher) GetDistAddr(_ struct{}, distAddr *string) error {
    *distAddr = config.GetStr("dist-addr")
    return nil
}

// Tells this node to connect to remote node, first getting the dist-addr
// the remote node wants to be identified as
func (d *Dispatcher) AddNode(node string, _ *struct{}) error {
    client,err := rpc.Dial("tcp",node)
    if err != nil {
        return err
    }

    var distAddr string
    err = client.Call("Dispatcher.GetDistAddr",struct{}{},&distAddr)
    if err != nil {
        return err
    }

    _,err = connectToNode(distAddr)
    return err
}

// Tells this node to connect to remote node, and then tell the remote
// node to connect back
func (d *Dispatcher) AddNodeTwoWay(node string, _ *struct{}) error {
    client,err := connectToNode(node)
    if err != nil {
        return err
    }

    return client.Call("Dispatcher.AddNode",config.GetStr("dist-addr"),nil)
}

// Tells this node to connect to remote node and all nodes that THAT node
// is connected to. If this node is already in a pool this won't work
func (d *Dispatcher) AddToPool(node string, _ *struct{}) error {
    if len(getNodes()) > 1 {
        return errors.New("This node is already part of another pool")
    }

    client,err := rpc.Dial("tcp",node)
    if err != nil {
        return err
    }

    var nodes []string
    err = client.Call("Dispatcher.GetNodes",struct{}{},&nodes)
    if err != nil {
        return err
    }

    for i := range nodes {
        err := d.AddNodeTwoWay(nodes[i],&struct{}{})
        if err != nil {
            return err
        }
    }

    return nil
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
