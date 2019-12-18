package consistenthash

import(
  "bytes"
  bt "github.com/tomdionysus/binarytree"
  "math/rand"
  "crypto/md5"
)

// NodeId is [NODE_ID_SIZE]byte, essentially a uint128
type NodeId [NODE_ID_SIZE]byte

// Return true if this key is less than the supplied NodeId.
func (me NodeId) LessThan(other bt.Comparable) bool {
  return bytes.Compare(me[:], getNodeIdSlice(other)) < 0
} 

// Return true if this key is equal to the supplied NodeId.
func (me NodeId) EqualTo(other bt.Comparable) bool {
  return bytes.Compare(me[:], getNodeIdSlice(other)) == 0
} 

// Return true if this key is greater than the supplied NodeId.
func (me NodeId) GreaterThan(other bt.Comparable) bool {
  return bytes.Compare(me[:], getNodeIdSlice(other)) > 0
} 

// Return an interface{} of the underlying [NODE_ID_SIZE]byte.
func (me NodeId) ValueOf() interface{} {
  return [NODE_ID_SIZE]byte(me)
}

// Internal function to convert [NODE_ID_SIZE]byte to []byte
func getNodeIdSlice(t bt.Comparable) []byte {
  x := t.ValueOf().([NODE_ID_SIZE]byte)
  return x[:]
}

// Return a new random NodeId
func NewRandomNodeId() NodeId {
  b := [NODE_ID_SIZE]byte{}
  for i:=0; i<NODE_ID_SIZE; i++ {
    b[i] = byte(rand.Intn(256))
  }
  x := NodeId(b)
  return x
}

// Return a new NodeId which is the MD5 hash of the supplied string.
func NewMD5NodeId(keystring string) NodeId {
  return NodeId(md5.Sum([]byte(keystring)))
}

