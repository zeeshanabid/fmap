package fmap

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
)

type Hash interface {
	Put(key interface{}, value interface{}) (interface{}, error)
	Get(key interface{}) (interface{}, error)
	Delete(key interface{}) (interface{}, error)
	Keys() []interface{}
	Values() []interface{}
	Length() uint64
}

const (
	defaultSize       = 4
	defaultLoadFactor = 0.65
)

type fmap struct {
	size       uint
	capacity   uint64
	keyCount   uint64
	loadFactor float32
	threshold  uint64
	keys       []interface{}
	values     []interface{}
}

func (m *fmap) Put(key interface{}, value interface{}) (interface{}, error) {
	i, err := m.getIndex(key)
	if err != nil {
		return nil, err
	}
	if m.keys[i] == nil {
		m.keyCount++
	}
	if m.keyCount > m.threshold {
		m.resize()
	}
	oldValue := m.values[i]
	m.keys[i] = key
	m.values[i] = value
	return oldValue, nil
}

func (m *fmap) Get(key interface{}) (interface{}, error) {
	i, err := m.getIndex(key)
	if err != nil {
		return nil, err
	}
	return m.values[i], nil
}

func (m *fmap) Delete(key interface{}) (interface{}, error) {
	i, err := m.getIndex(key)
	if err != nil {
		return nil, err
	}
	if m.keys[i] != nil {
		m.keyCount--
	}
	value := m.values[i]
	m.keys[i] = nil
	m.values[i] = nil
	return value, nil
}

func (m *fmap) Keys() []interface{} {
	var keys []interface{}
	for _, k := range m.keys {
		if k != nil {
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
	m.setValues(defaultSize, defaultLoadFactor)
	return &m
}

func (m *fmap) setValues(size uint, loadFactor float32) {
	m.size = size
	m.capacity = 1 << m.size
	m.loadFactor = loadFactor
	m.threshold = uint64(float32(m.capacity) * m.loadFactor)
	m.keyCount = 0
	m.keys = make([]interface{}, m.capacity, m.capacity)
	m.values = make([]interface{}, m.capacity, m.capacity)
}

func (m *fmap) resize() {
	keys, values := m.keys, m.values
	m.setValues(m.size+1, m.loadFactor)

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
