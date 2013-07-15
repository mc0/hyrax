package dist

import (
    "hyrax-server/custom"
    "log"
    "net/rpc"
)

// The different types of messages rpc will pull
type msgtype int
const (
    ECHO msgtype = iota //for testing
    MONPUSH
)

type RpcPayload struct {
    MsgType msgtype
    Payload interface{}
}


func Echo(client *rpc.Client,e string) {
    r := RpcPayload{ ECHO, e }
    var b byte
    err := client.Call("Dispatcher.Give",&r,&b)
    if err != nil {
        log.Fatal("give error:", err)
    }
}

func processMessage(pay *RpcPayload) {
    switch pay.MsgType {
        case ECHO:
            log.Println("got echo message:",pay.Payload.(string))
        case MONPUSH:
            //TODO try to make this pointer
            monpay := pay.Payload.(custom.MonPushPayload)
            custom.MonPushAlert(&monpay)
        default:
            log.Println("Unknown rpcPayload.msgtype",pay.MsgType)
    }
}

