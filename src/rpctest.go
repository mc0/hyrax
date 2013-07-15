package main

import (
    "net/rpc"
    "log"
    "hyrax-server/dist"
)

func main() {
    client, err := rpc.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    dist.Echo(client,"OHAI")
}
