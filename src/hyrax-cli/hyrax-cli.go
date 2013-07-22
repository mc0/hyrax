package main

import (
    "flag"
    "fmt"
    "net/rpc"
)

var client *rpc.Client

func main() {
    var distAddr string
    flag.StringVar(&distAddr,"dist-addr","127.0.0.1:4400","hostname:port of the hyrax instance's dist port")

    var getPool bool
    flag.BoolVar(&getPool,"get-pool",false,"Return the list of nodes currently in the pool")

    var addToPool string
    flag.StringVar(&addToPool,"add-to-pool","remote-addr","Adds this node to the pool the node at remote-addr is part of")

    flag.Parse()

    client,err := rpc.Dial("tcp",distAddr)
    if err != nil {
        fmt.Println(err)
        return
    }

    switch {
        case getPool:
            var ret []string
            err = client.Call("Dispatcher.GetNodes",struct{}{},&ret)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Println(ret)
        case addToPool != "add-to-pool":
            err = client.Call("Dispatcher.AddToPool",addToPool,nil)
            if err != nil {
                fmt.Println(err)
                return
            }
        default:
            fmt.Println("Unknown command, try -h")
    }

}
