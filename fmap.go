package fmap

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
)

type Hash interface {
	Put(key interface{}, value interface{}) error
	Get(key interface{}) (interface{}, bool, error)
	Delete(key interface{}) error
	Pop(key interface{}) (interface{}, error)
	Keys() []interface{}
	Values() []interface{}
	Length() uint64
}

const (
	defaultSize          uint    = 4
	defaultMaxLoadFactor float32 = 0.65
	defaultMinLoadFactor float32 = 0.30
)

type d struct{}

var deleted = d{}

type fmap struct {
	size          uint
	capacity      uint64
	maxIndex      uint64
	keyCount      uint64
	maxLoadFactor float32
	minLoadFactor float32
	maxThreshold  uint64
	minThreshold  uint64
	keys          []interface{}
	values        []interface{}
}

func (m *fmap) Put(key interface{}, value interface{}) error {
	if m.keyCount >= m.maxThreshold {
		m.increaseSize(1)
	}
	i, err := m.getIndex(key)
	if err != nil {
		return err
	}

	delFound := false
	delIndex := uint(0)
	for m.keys[i] != nil && m.keys[i] != key {
		if m.keys[i] == deleted && !delFound {
			delFound = true
			delIndex = i
		}
		i = uint(uint64(i+1) & m.maxIndex)
	}
	if m.keys[i] != key {
		m.keyCount++
	}
	if m.keys[i] == key && delFound {
		m.keys[i] = nil
		m.values[i] = nil
	}
	if delFound {
		i = delIndex
	}
	m.keys[i] = key
	m.values[i] = value

	return nil
}

func (m *fmap) Get(key interface{}) (interface{}, bool, error) {
	i, err := m.getIndex(key)
	if err != nil {
		return nil, false, err
	}
	for m.keys[i] != nil && m.keys[i] != key {
		i = uint(uint64(i+1) & m.maxIndex)
	}
	if m.keys[i] != key {
		return nil, false, nil
	}
	return m.values[i], true, nil
}

func (m *fmap) Delete(key interface{}) error {
	_, err := m.remove(key)
	return err
}

func (m *fmap) Pop(key interface{}) (interface{}, error) {
	return m.remove(key)
}

func (m *fmap) remove(key interface{}) (interface{}, error) {
	if m.keyCount <= m.minThreshold {
		m.decreaseSize(1)
	}
	i, err := m.getIndex(key)
	if err != nil {
		return nil, err
	}
	for m.keys[i] != nil && m.keys[i] != key {
		i = uint(uint64(i+1) & m.maxIndex)
	}
	if m.keys[i] != key {
		return nil, nil
	}

	m.keyCount--
	old := m.values[i]
	m.keys[i] = deleted
	m.values[i] = nil
	return old, nil
}

func (m *fmap) Keys() []interface{} {
	var keys []interface{}
	for _, k := range m.keys {
		if k != nil && k != deleted {
			keys = append(keys, k)
		}
	}
	return keys
}

func (m *fmap) Values() []interface{} {
	var values []interface{}
	for _, v := range m.values {
		if v != nil {
			values = append(values, v)
		}
	}
	return values
}

func (m *fmap) Length() uint64 {
	return m.keyCount
}

func New() Hash {
	m := fmap{}
	m.setValues(defaultSize, defaultMaxLoadFactor, defaultMinLoadFactor)
	return &m
}

func (m *fmap) setValues(size uint, maxLoadFactor, minLoadFactor float32) {
	m.size = size
	m.capacity = 1 << m.size
	m.maxIndex = m.capacity - 1
	m.maxLoadFactor = maxLoadFactor
	m.minLoadFactor = minLoadFactor
	m.maxThreshold = uint64(float32(m.capacity) * m.maxLoadFactor)
	m.minThreshold = uint64(float32(m.capacity) * m.minLoadFactor)
	m.keyCount = 0
	m.keys = make([]interface{}, m.capacity, m.capacity)
	m.values = make([]interface{}, m.capacity, m.capacity)
}

func (m *fmap) increaseSize(size uint) {
	m.resize(m.size + size)
}

func (m *fmap) decreaseSize(size uint) {
	newSize := maxInt(defaultSize, m.size-size)
	if newSize != m.size {
		m.resize(newSize)
	}
}

func (m *fmap) resize(newSize uint) {
	keys, values := m.keys, m.values
	m.setValues(newSize, m.maxLoadFactor, m.minLoadFactor)

	for i := 0; i < len(keys); i++ {
		if keys[i] != nil {
			m.Put(keys[i], values[i])
		}
	}
}

func (m *fmap) getIndex(key interface{}) (uint, error) {
	c, err := getCode(key)
	if err != nil {
		return 0, err
	}
	return m.fibonacciHash(c), nil
}

func (m *fmap) fibonacciHash(k uint64) uint {
	s := 64 - m.size
	k ^= k >> s
	return uint(uint64(11400714819323198485*uint64(k)) >> s)
}

func hashCode(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}

func getBytes(d interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getCode(k interface{}) (code uint64, err error) {
	b, err := getBytes(k)
	if err != nil {
		return 0, err
	}
	return hashCode(b), nil
}

func maxInt(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}
