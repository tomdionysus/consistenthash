package consistenthash

import(
  bt "github.com/tomdionysus/binarytree"
  "errors"
)

// An array of Keys representing the distribution of a node.
type NodeDistribution [DISTRIBUTION_MAX]Key

type Redistribution struct {
  SourceNodeID NodeId
  DestinationNodeID NodeId
  Start Key
  End Key
}

// The details of a node on the network, including its
// ID, host address, and Distrubution.
type ServerNetworkNode struct {
  ID NodeId
  HostAddr string
  Distribution NodeDistribution
}

// An actual server node, containing the local node's ServerNetworkNode 
// information and additionally the actual network, allowing registration
// and deregistration of other ServerNetworkNodes
type ServerNode struct {
  ServerNetworkNode
  NetworkNodes map[NodeId]*ServerNetworkNode
  Network *bt.Tree
}

// Return a pointer to a new ServerNode with the specified host address.
func NewServerNode(hostAddr string) *ServerNode {
  node := &ServerNode{ 
    NetworkNodes: map[NodeId]*ServerNetworkNode{}, 
    Network: bt.NewTree(),
  }
  node.ID = NewRandomNodeId()
  node.HostAddr = hostAddr
  node.Init()
  return node
}

// Initialise the ServerNode, creating a new random distribution for
// this node and registering it with the network. 
func (me *ServerNode) Init() {
  me.Distribution = NodeDistribution{}
  for x:=0; x<DISTRIBUTION_MAX; x++ {
    me.Distribution[x] = NewRandomKey()
    me.Network.Set(me.Distribution[x], &me.ServerNetworkNode)
  }
}

// Register another node with the network.
func (me *ServerNode) RegisterNode(server *ServerNetworkNode) error {
  if server.ID == me.ID {
    return errors.New("Cannot register a node with itself")
  }
  _, found := me.NetworkNodes[server.ID] 
  if found {
    return errors.New("Node is already registered")
  }
  
  me.NetworkNodes[server.ID] = server
  for _, insertnodekey := range server.Distribution {
    me.Network.Set(insertnodekey, server)
  }

  me.Network.Balance()
  return nil
}

// Deregister (remove) a previously registered node from the network.
func (me *ServerNode) DeregisterNode(server *ServerNetworkNode) error {
  if _, found := me.NetworkNodes[server.ID]; !found {
    return errors.New("Node is not registered")
  }
  delete(me.NetworkNodes, server.ID)
  for _, x := range server.Distribution { 
    me.Network.Clear(x)
  }

  me.Network.Balance()
  return nil
}

func (me *ServerNode) NodeRegistered(nodeid NodeId) bool {
  _, found := me.NetworkNodes[nodeid]
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
  var nodes map[NodeId]*ServerNetworkNode = map[NodeId]*ServerNetworkNode{}
  if totalNodes > len(me.NetworkNodes) { totalNodes = len(me.NetworkNodes)+1 }
  
  for len(nodes) < totalNodes {
    found, ky, nd := me.Network.Next(key)
    if !found { ky, nd = me.Network.First() }
    key = Key(ky.ValueOf().([NODE_ID_SIZE]byte))
    var node *ServerNetworkNode = nd.(*ServerNetworkNode)
    _, found = nodes[node.ID]
    if !found { nodes[node.ID] = node }
  }

  var out []*ServerNetworkNode = []*ServerNetworkNode{}
  for _, node :=range nodes { out = append(out, node) }
  return out
}