package dist

import (
    "net"
    //"net/http"
    "net/rpc"
    "log"
)

var distCh = make(chan *RpcPayload,1024)
type Dispatcher byte

// Give is where all data from other nodes. Give simply gives the payload
// to this node to be dealt with
func (b *Dispatcher) Give(pay RpcPayload, r *byte) error {
    distCh <- &pay
    return nil
}

func Setup() {
    disp := new(Dispatcher)
    rpc.Register(disp)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":1234")
    if e != nil {
        log.Fatal("listen error:", e)
    }
    go rpc.Accept(l)
    //go http.Serve(l, nil)
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
