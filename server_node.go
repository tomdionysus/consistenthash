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

// An array of Keys representing the distribution
// of a node.
type NodeDistribution [DISTRIBUTION_MAX]Key

type Redistribution struct {
  SourceNodeID Key
  DestinationNodeID Key
  Start Key
  End Key
}

// The details of a node on the network, including its
// ID, host address, and Distrubution.
type ServerNetworkNode struct {
  ID Key
  HostAddr string
  Distribution NodeDistribution
}

// An actual server node, containing the local node's ServerNetworkNode 
// information and additionally the actual network, allowing registration
// and deregistration of other ServerNetworkNodes
type ServerNode struct {
  ServerNetworkNode
  NetworkNodes map[Key]*ServerNetworkNode
  Network *bt.Tree
}

// Return a pointer to a new ServerNode with the specified host address.
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

// Initialise the ServerNode, creating a new random distribution for
// this node and registering it with the network. 
func (me *ServerNode) Init() {
  me.Distribution = NodeDistribution{}
  for x:=0; x<DISTRIBUTION_MAX; x++ {
    rand.Seed(time.Now().UTC().UnixNano())
    me.Distribution[x] = NewRandomKey()
    me.Network.Set(me.Distribution[x], &me.ServerNetworkNode)
  }
}

// Register another node with the network.
func (me *ServerNode) RegisterNode(server *ServerNetworkNode) ([]Redistribution, error) {
  if server.ID == me.ID {
    return nil, errors.New("Cannot register a node with itself")
  }
  _, found := me.NetworkNodes[server.ID] 
  if found {
    return nil, errors.New("Node is already registered")
  }

  redistributions := []Redistribution{}
  
  me.NetworkNodes[server.ID] = server
  for _, insertnodekey := range server.Distribution {

    found, key, pn := me.Network.Previous(insertnodekey)
    if !found { key, pn = me.Network.Last() }

    previousNode := pn.(*ServerNetworkNode)
    if previousNode.ID == server.ID { continue }

    redist := Redistribution{
      SourceNodeID: previousNode.ID,
      DestinationNodeID: server.ID,
      Start: key.(Key),
      End: insertnodekey,
    } 
    redistributions = append(redistributions, redist)
    me.Network.Set(insertnodekey, server)
  }

  me.Network.Balance()
  return redistributions, nil
}

// Deregister (remove) a previously registered node from the network.
func (me *ServerNode) DeregisterNode(server *ServerNetworkNode) error {
  if _, found := me.NetworkNodes[server.ID]; !found {
    return errors.New("Node is not registered")
  }
  delete(me.NetworkNodes, server.ID)
  for _, x := range server.Distribution { me.Network.Clear(x);  }

  me.Network.Balance()
  return nil
}

func (me *ServerNode) NodeRegistered(key Key) bool {
  _, found := me.NetworkNodes[key]
  return found
}

// Return the 'primary' ServerNetworkNode responsible for the supplied data key.
func (me *ServerNode) GetNodeFor(key Key) *ServerNetworkNode {
  found, _, node := me.Network.Next(key)
  if !found { _, node = me.Network.First() }
  return node.(*ServerNetworkNode)
}

// Return at most totalNodes distinct ServerNetworkNodes after the supplied Key.
func (me *ServerNode) GetNodesFor(key Key, totalNodes int) []*ServerNetworkNode {
  var nodes map[Key]*ServerNetworkNode = map[Key]*ServerNetworkNode{}
  if totalNodes > len(me.NetworkNodes) { totalNodes = len(me.NetworkNodes)+1 }
  
  for len(nodes) < totalNodes {
    found, ky, nd := me.Network.Next(key)
    if !found { ky, nd = me.Network.First() }
    key = Key(ky.ValueOf().([16]byte))
    var node *ServerNetworkNode = nd.(*ServerNetworkNode)
    _, found = nodes[node.ID]
    if !found { nodes[node.ID] = node }
  }

  var out []*ServerNetworkNode = []*ServerNetworkNode{}
  for _, node :=range nodes { out = append(out, node) }
  return out
}