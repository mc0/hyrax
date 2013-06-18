package dispatch

import (
    "hyrax/storage"
    "errors"
)

func DoCommand(rawJson []byte) ([]byte,error) {
    cmd,err := DecodeCommand(rawJson)
    if err != nil { return nil,err }

    ret,err := doCommandWrap(cmd)
    if err != nil {
        return EncodeError(err.Error())
    }

    return EncodeMessage(cmd.Command,ret)
}

func doCommandWrap(cmd *Command) (interface{},error) {
    if !CommandExists(&cmd.Command) {
        return nil,errors.New("Unsupported command")
    }

    if CommandModifies(&cmd.Command) {
        for i := range(cmd.Payload) {
            if !CheckAuth(&cmd.Payload[i]) {
                return nil,errors.New("cannot auth for domain:"+cmd.Payload[i].Domain+
                                                     " id:"+cmd.Payload[i].Id)
            }
        }
    }

    if len(cmd.Payload) == 0 {
        return nil,errors.New("empty payload")
    }

    numArgs := 0
    for i:=0; i<len(cmd.Payload); i++ {
        numArgs++
        numArgs += len(cmd.Payload[i].Values)
    }

    args := make([]interface{},0,numArgs)
    for i:=0; i<len(cmd.Payload); i++ {
        pay := &cmd.Payload[i]
        strKey := storage.CreateKey(pay.Domain,pay.Id)
        args = append(args,strKey)
        for j:=0; j<len(pay.Values); j++ {
            args = append(args,pay.Values[j])
        }
    }

    r,err := storage.Cmd(cmd.Command,args)
    if err != nil { return nil,err }

    if CommandReturnsMap(&cmd.Command) {
        return storage.RedisListToMap(r.([]string))
    }

    return r,nil
}
