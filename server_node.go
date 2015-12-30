package consistenthash

import(
  "math/rand"
  "time"
  bt "github.com/tomdionysus/binarytree"
  "errors"
)

const(
  DISTRIBUTION_MAX = 512
  NETWORK_ID_SIZE_BYTES = 16
  KEY_SIZE_BYTES = 16
)

type NodeDistribution [DISTRIBUTION_MAX]Key

type ServerNetworkNode struct {
  ID Key
  HostAddr string
  Distribution NodeDistribution
}

type ServerNode struct {
  ServerNetworkNode
  NetworkNodes map[Key]*ServerNetworkNode
  Network *bt.Tree
}

func NewServerNode(hostAddr string) *ServerNode {
  node := &ServerNode{ 
    NetworkNodes: map[Key]*ServerNetworkNode{}, 
    Network: bt.NewTree(),
  }
  node.ID = NewRandomKey()
  node.HostAddr = hostAddr
  node.Init()
  return node
}

func (me *ServerNode) Init() {
  me.Distribution = NodeDistribution{}
  for x:=0; x<DISTRIBUTION_MAX; x++ {
    rand.Seed(time.Now().UTC().UnixNano())
    me.Distribution[x] = NewRandomKey()
    me.Network.Set(me.Distribution[x], &me.ServerNetworkNode)
  }
}

func (me *ServerNode) AddToNetwork(server *ServerNetworkNode) error {
  if server.ID == me.ID {
    return errors.New("Cannot register a node with itself")
  }
  if _, found := me.NetworkNodes[server.ID]; found {
    return errors.New("Node is already registered")
  }
  me.NetworkNodes[server.ID] = server
  for _, x := range server.Distribution { me.Network.Set(x, server);  }

  me.Network.Balance()
  return nil
}

func (me *ServerNode) RemoveFromNetwork(server *ServerNetworkNode) error {
  if _, found := me.NetworkNodes[server.ID]; !found {
    return errors.New("Node is not registered")
  }
  delete(me.NetworkNodes, server.ID)
  for _, x := range server.Distribution { me.Network.Clear(x);  }

  me.Network.Balance()
  return nil
}

func (me *ServerNode) GetNodeFor(key Key) *ServerNetworkNode {
  found, _, node := me.Network.Next(key)
  if !found { _, node = me.Network.First() }
  return node.(*ServerNetworkNode)
}