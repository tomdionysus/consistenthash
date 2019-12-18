package consistenthash

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestNodeIdLessThan(t *testing.T) {
  a := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 }
  b := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1 }
  
  x := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0x20,0 }
  y := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0x30,0x20,0 }

  c := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0xAA,0xEE,0xFF,0xFF }
  d := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0xAA,0xEE }

  assert.True(t, a.LessThan(b))
  assert.False(t, b.LessThan(a))

  assert.True(t, a.LessThan(c))
  assert.False(t, c.LessThan(a))

  assert.True(t, b.LessThan(x))
  assert.False(t, x.LessThan(b))

  assert.True(t, x.LessThan(y))
  assert.False(t, y.LessThan(x))

  assert.True(t, a.LessThan(d))
  assert.False(t, d.LessThan(a))

  assert.True(t, d.LessThan(c))
  assert.False(t, c.LessThan(d))

  assert.False(t, c.LessThan(c))
}

func TestNodeIdEqualTo(t *testing.T) {
  a := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 }
  b := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 }
  
  x := NodeId{ 0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15 }
  y := NodeId{ 15,14,13,12,11,10,9,8,7,6,5,4,3,2,1,0 }
  z := NodeId{ 0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15 }

  assert.True(t, a.EqualTo(b))
  assert.False(t, x.EqualTo(a))
  assert.False(t, a.EqualTo(x))
  assert.False(t, y.EqualTo(a))
  assert.False(t, a.EqualTo(y))
  assert.False(t, z.EqualTo(a))
  assert.False(t, a.EqualTo(z))

  assert.True(t, x.EqualTo(z))
  assert.True(t, z.EqualTo(x))
}

func TestNodeIdGreaterThan(t *testing.T) {
  a := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 }
  b := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1 }
  
  x := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0x20,0 }
  y := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0x30,0x20,0 }

  c := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0xAA,0xEE,0xFF,0xFF }
  d := NodeId{ 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0xAA,0xEE }

  assert.False(t, a.GreaterThan(b))
  assert.True(t, b.GreaterThan(a))

  assert.False(t, a.GreaterThan(c))
  assert.True(t, c.GreaterThan(a))

  assert.False(t, b.GreaterThan(x))
  assert.True(t, x.GreaterThan(b))

  assert.False(t, x.GreaterThan(y))
  assert.True(t, y.GreaterThan(x))

  assert.False(t, a.GreaterThan(d))
  assert.True(t, d.GreaterThan(a))

  assert.False(t, d.GreaterThan(c))
  assert.True(t, c.GreaterThan(d))

  assert.False(t, c.GreaterThan(c))
}

func TestNodeIdNewRandomNodeId(t *testing.T) {
  x := NewRandomNodeId()
  y := NewRandomNodeId()
  z := NewRandomNodeId()

  assert.False(t, x.EqualTo(y))
  assert.False(t, y.EqualTo(z))
  assert.False(t, z.EqualTo(x))
}

func TestNodeIdNewMD5NodeId(t *testing.T) {
  x := NewMD5NodeId("")
  y := NewMD5NodeId("helloworld")
  z := NewMD5NodeId("1) First Solve The Problem. 2) IT is All About People. 3) Do not Produce, or Tolerate, Shit.")

  a := NodeId{ 0xD4,0x1D,0x8C,0xD9,0x8F,0x00,0xB2,0x04,0xE9,0x80,0x09,0x98,0xEC,0xF8,0x42,0x7E }
  b := NodeId{ 0xFC,0x5E,0x03,0x8D,0x38,0xA5,0x70,0x32,0x08,0x54,0x41,0xE7,0xFE,0x70,0x10,0xB0 }
  c := NodeId{ 0x95,0x83,0x6E,0xE0,0xD0,0xBA,0x6C,0x85,0xB4,0xCE,0xD2,0x86,0x9A,0x02,0x34,0x10 }

  assert.True(t, x.EqualTo(a))
  assert.True(t, y.EqualTo(b))
  assert.True(t, z.EqualTo(c))
}





