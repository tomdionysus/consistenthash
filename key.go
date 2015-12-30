package consistenthash

import(
  "bytes"
  bt "github.com/tomdionysus/binarytree"
  "time"
  "math/rand"
  "crypto/md5"
)

// Key is [16]byte, essentially a uint128
type Key [16]byte

// Return true if this key is less than the supplied Key.
func (me Key) LessThan(other bt.Comparable) bool {
  return bytes.Compare(me[:], getSlice(other)) < 0
} 

// Return true if this key is equal to the supplied Key.
func (me Key) EqualTo(other bt.Comparable) bool {
  return bytes.Compare(me[:], getSlice(other)) == 0
} 

// Return true if this key is greater than the supplied Key.
func (me Key) GreaterThan(other bt.Comparable) bool {
  return bytes.Compare(me[:], getSlice(other)) > 0
} 

// Return an interface{} of the underlying [16]byte.
func (me Key) ValueOf() interface{} {
  return [16]byte(me)
}

// Internal function to convert [16]byte to []byte
func getSlice(t bt.Comparable) []byte {
  x := t.ValueOf().([16]byte)
  return x[:]
}

// Return a new random Key
func NewRandomKey() Key {
  rand.Seed(time.Now().UTC().UnixNano())
  b := [16]byte{}
  for i:=0; i<16; i++ {
    b[i] = byte(rand.Intn(256))
  }
  x := Key(b)
  return x
}

// Return a new Key which is the MD5 hash of the supplied string.
func NewMD5Key(keystring string) Key {
  return Key(md5.Sum([]byte(keystring)))
}

