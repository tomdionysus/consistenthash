package consistenthash

import(
  "bytes"
  bt "github.com/tomdionysus/binarytree"
  "math/rand"
  "crypto/md5"
)

// Key is [KEY_SIZE]byte, essentially a uint128
type Key [KEY_SIZE]byte

// Return true if this key is less than the supplied Key.
func (me Key) LessThan(other bt.Comparable) bool {
  return bytes.Compare(me[:], getKeySlice(other)) < 0
} 

// Return true if this key is equal to the supplied Key.
func (me Key) EqualTo(other bt.Comparable) bool {
  return bytes.Compare(me[:], getKeySlice(other)) == 0
} 

// Return true if this key is greater than the supplied Key.
func (me Key) GreaterThan(other bt.Comparable) bool {
  return bytes.Compare(me[:], getKeySlice(other)) > 0
} 

// Return an interface{} of the underlying [KEY_SIZE]byte.
func (me Key) ValueOf() interface{} {
  return [KEY_SIZE]byte(me)
}

// Internal function to convert [NODE_ID_SIZE]byte to []byte
func getKeySlice(t bt.Comparable) []byte {
  x := t.ValueOf().([KEY_SIZE]byte)
  return x[:]
}

// Return a new random Key
func NewRandomKey() Key {
  b := [KEY_SIZE]byte{}
  for i:=0; i<KEY_SIZE; i++ {
    b[i] = byte(rand.Intn(256))
  }
  x := Key(b)
  return x
}

// Return a new Key which is the MD5 hash of the supplied string.
func NewMD5Key(keystring string) Key {
  return Key(md5.Sum([]byte(keystring)))
}

