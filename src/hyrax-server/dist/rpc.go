package dist

import (
    "hyrax-server/custom"
    "log"
)

// The different types of messages rpc will pull
type msgtype int
const (
    MONPUSH msgtype = iota
)

type rpcPayload struct {
    msgtype msgtype
    payload interface{}
}

func processMessage(pay *rpcPayload) {
    switch pay.msgtype {
        case MONPUSH:
            //TODO try to make this pointer
            monpay := pay.payload.(custom.MonPushPayload)
            custom.MonPushAlert(&monpay)
        default:
            log.Println("Unknown rpcPayload.msgtype",pay.msgtype)
    }
}
