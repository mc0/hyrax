package dist

import (
    "sync"
    "net/rpc"
    "fmt"
)

var nodeSetLock = sync.RWMutex{}
var nodeSet = map[string]*rpc.Client{}

// connectToNode connects to a node's dist port and stores that
// information
func connectToNode(node string) (*rpc.Client,error) {
    nodeSetLock.Lock()
    defer nodeSetLock.Unlock()

    _,ok := nodeSet[node]
    if ok {
        return nil,fmt.Errorf("Node '%s' already registered")
    }

    client, err := rpc.Dial("tcp", node)
    if err != nil {
        return nil,err
    }

    nodeSet[node] = client
    return client,nil
}

// disconnectFromNode disconnects from a node's dist port and deletes
// that information. Returns whether or not the node was in the set
// previously
func disconnectFromNode(node string) bool {
    nodeSetLock.Lock()
    defer nodeSetLock.Unlock()

    client,ok := nodeSet[node]
    if !ok {
        return false
    }

    client.Close()
    delete(nodeSet,node)
    return true
}

// getNode gets the *rpc.Client object for a given node name. The bool
// portion of the return will be false if the node isn't found
func getNode(node string) (*rpc.Client,bool) {
    nodeSetLock.RLock()
    defer nodeSetLock.RUnlock()

    client,ok := nodeSet[node]
    if !ok {
        return nil,false
    }

    return client,true
}

// getNodes returns the list of nodes this node is currently connected to
func getNodes() []string {
    nodeSetLock.RLock()
    defer nodeSetLock.RUnlock()
    ret := make([]string,0,len(nodeSet))

    for k,_ := range nodeSet {
        ret = append(ret,k)
    }

    return ret
}
