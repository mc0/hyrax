package dist

import (
    "hyrax-server/custom"
    "log"
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

func processMessage(pay *RpcPayload) {
    switch pay.MsgType {
        case MONPUSH:
            //TODO try to make this pointer
            monpay := pay.Payload.(custom.MonPushPayload)
            custom.MonPushAlert(&monpay)
        default:
            log.Println("Unknown rpcPayload.msgtype",pay.MsgType)
    }
}

