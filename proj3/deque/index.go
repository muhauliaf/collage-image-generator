package deque

import "sync/atomic"

// Index represents 64-bit uint which is combination of 32-bit value and 32-bit stamp
type Index struct {
	packed uint64
}

// pack combines value and stamp into an uint64
func pack(value int32, stamp int32) uint64 {
	return (uint64(value) << 32) | uint64(stamp)
}

// unpack extracts value and stamp from an uint64
func unpack(combined uint64) (int32, int32) {
	value := int32(combined >> 32)
	stamp := int32(combined & 0xFFFFFFFF)
	return value, stamp
}

// NewIndex creates a new index from value and stamp
func NewIndex(value int32, stamp int32) *Index {
	return &Index{pack(value, stamp)}
}

// Get extracts value and stamp from the Index
func (i *Index) Get() (int32, int32) {
	packed := atomic.LoadUint64(&i.packed)
	return unpack(packed)
}

// Set assigns value and stamp to the Index
func (i *Index) Set(value int32, stamp int32) {
	packed := pack(value, stamp)
	atomic.StoreUint64(&i.packed, packed)
}

// Get extracts value from the Index
func (i *Index) Value() int32 {
	return int32(atomic.LoadUint64(&i.packed) >> 32)
}

// Get extracts stamp from the Index
func (i *Index) Stamp() int32 {
	return int32(atomic.LoadUint64(&i.packed) & 0xFFFFFFFF)
}

// CompareAndSwap performs CAS for both value and stamp at the same time
func (i *Index) CompareAndSwap(oldValue int32, newValue int32, oldStamp int32, newStamp int32) bool {
	oldVal := pack(oldValue, oldStamp)
	newVal := pack(newValue, newStamp)
	return atomic.CompareAndSwapUint64(&i.packed, oldVal, newVal)
}
