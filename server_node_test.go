package consistenthash

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestNewServerNode(t *testing.T) {
  inst := NewServerNode("HOSTADDR")

  assert.NotNil(t, inst.ID)
  assert.NotNil(t, inst.Distribution)
  assert.Equal(t, "HOSTADDR", inst.HostAddr)
}

func TestRegisterNode(t *testing.T) {
  inst1 := NewServerNode("host1")
  inst2 := NewServerNode("host2")
  inst3 := NewServerNode("host3")

  // Can't add node to itself
  err := inst1.RegisterNode(&inst1.ServerNetworkNode)
  assert.NotNil(t, err)
  assert.Equal(t, "Cannot register a node with itself", err.Error())

  // Should register OK
  err = inst1.RegisterNode(&inst2.ServerNetworkNode)
  assert.Nil(t, err)

  // Can't register twice
  err = inst1.RegisterNode(&inst2.ServerNetworkNode)
  assert.NotNil(t, err)
  assert.Equal(t, "Node is already registered", err.Error())

  // Should register multiple nodes
  err = inst1.RegisterNode(&inst3.ServerNetworkNode)
  assert.Nil(t, err)
}

func TestDeregisterNode(t *testing.T) {
  inst1 := NewServerNode("host1")
  inst2 := NewServerNode("host2")

  // Can't deregister unregistered node
  err := inst1.DeregisterNode(&inst2.ServerNetworkNode)
  assert.NotNil(t, err)
  assert.Equal(t, "Node is not registered", err.Error())

  // Should deregister OK
  err = inst1.RegisterNode(&inst2.ServerNetworkNode)
  assert.Nil(t, err)
  err = inst1.DeregisterNode(&inst2.ServerNetworkNode)
  assert.Nil(t, err)

  // Can't deregister twice
  err = inst1.RegisterNode(&inst2.ServerNetworkNode)
  assert.Nil(t, err)
  err = inst1.DeregisterNode(&inst2.ServerNetworkNode)
  assert.Nil(t, err)
  err = inst1.DeregisterNode(&inst2.ServerNetworkNode)
  assert.NotNil(t, err)
  assert.Equal(t, "Node is not registered", err.Error())
}

func TestGetNodeFor(t *testing.T) {
  inst1 := NewServerNode("host1")
  inst2 := NewServerNode("host2")
  inst3 := NewServerNode("host3")

  inst1.RegisterNode(&inst2.ServerNetworkNode)
  inst1.RegisterNode(&inst3.ServerNetworkNode)
  inst2.RegisterNode(&inst1.ServerNetworkNode)
  inst2.RegisterNode(&inst3.ServerNetworkNode)
  inst3.RegisterNode(&inst1.ServerNetworkNode)
  inst3.RegisterNode(&inst2.ServerNetworkNode)

  // Should agree on key placement
  k1 := NewRandomKey()

  nodeId1 := inst1.GetNodeFor(k1)
  nodeId2 := inst2.GetNodeFor(k1)
  nodeId3 := inst3.GetNodeFor(k1)

  assert.Equal(t, nodeId1, nodeId2)
  assert.Equal(t, nodeId2, nodeId3)

  // And go around the circle

  // There's a 1 in 2^128 chance this will randomly fail.
  // If you see this actually happen you probably just broke your universe.
  k2 := Key{ 0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF }
  nodeId1 = inst1.GetNodeFor(k2)
  nodeId2 = inst2.GetNodeFor(k2)
  nodeId3 = inst3.GetNodeFor(k2)

  assert.Equal(t, nodeId1, nodeId2)
  assert.Equal(t, nodeId2, nodeId3)
}


