package dist

import (
    "net"
    "net/rpc"
    "hyrax-server/config"
    "errors"
    "log"
)

func getRemoteDistAddr(node string) (string,error) {
    client,err := rpc.Dial("tcp",node)
    defer client.Close()
    if err != nil {
        return "",err
    }

    var distAddr string
    err = client.Call("Dispatcher.GetDistAddr",struct{}{},&distAddr)
    return distAddr,err
}

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
    distAddr,err := getRemoteDistAddr(node)
    if err != nil {
        return err
    }

    log.Println("Generating a dist connection to:",distAddr)
    _,err = connectToNode(distAddr)
    if err != nil {
        log.Println("Error generating connection to",distAddr,":",err.Error())
    }
    return err
}

// Tells this node to connect to remote node, and then tell the remote
// node to connect back
func (d *Dispatcher) AddNodeTwoWay(node string, _ *struct{}) error {
    err := d.AddNode(node,nil)
    if err != nil {
        return err
    }

    distAddr,err := getRemoteDistAddr(node)
    if err != nil {
        return err
    }

    client,_ := getNode(distAddr)

    return client.Call("Dispatcher.AddNode",config.GetStr("dist-addr"),nil)
}

// Tells this node to connect to remote node and all nodes that THAT node
// is connected to. If this node is already in a pool this won't work
func (d *Dispatcher) AddToPool(node string, _ *struct{}) error {
    if len(getNodes()) > 1 {
        return errors.New("This node is already part of another pool")
    }

    client,err := rpc.Dial("tcp",node)
    defer client.Close()
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

// Tells this node to disconnect from remote node, first getting the dist-addr
// the remote node is identified as
func (d *Dispatcher) RemoveNode(node string, _ *struct{}) error {
    distAddr,err := getRemoteDistAddr(node)
    if err != nil {
        return err
    }

    log.Println("Disconnecting from:",distAddr)
    ok := disconnectFromNode(distAddr)
    if !ok {
        log.Println("No connection to",distAddr,"found")
    }
    return nil
}

// Tells this node to tell the remote node to disconnect from it, then
// disconnects from the remote node
func (d *Dispatcher) RemoveNodeTwoWay(node string, _ *struct{}) error {
    client,ok := getNode(node)
    if !ok {
        return errors.New("Remote node not currently connected")
    }

    err := client.Call("Dispatcher.RemoveNode",config.GetStr("dist-addr"),nil)
    if err != nil {
        return err
    }

    return d.RemoveNode(node,nil)
}

// Tells this node to disconnect from the pool it's in
func (d *Dispatcher) RemoveFromPool(_ struct{}, _ *struct{}) error {
    nodes :=  getNodes()
    if len(nodes) <= 1 {
        return errors.New("This node is not part of a pool")
    }

    thisDistAddr := config.GetStr("dist-addr")
    for i := range nodes {
        if nodes[i] == thisDistAddr { continue; }
        err := d.RemoveNodeTwoWay(nodes[i],&struct{}{})

        //If there's an error just disconnect locally and continue, nothing we can do
        if err != nil {
            log.Println("Got error disconnecting from",nodes[i],":",err.Error())
            disconnectFromNode(nodes[i])
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
