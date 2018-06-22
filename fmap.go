package fmap

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
)

type Hash interface {
	Put(k interface{}, v interface{}) (interface{}, error)
	Get(k interface{}) (interface{}, error)
}

const (
	defaultSize = 4
)
type fmap struct {
	size     uint
	capacity uint64
}

func (m *fmap) Put(k interface{}, v interface{}) (interface{}, error) {
	return nil, nil
}

func (m *fmap) Get(k interface{}) (interface{}, error) {
	return nil, nil
}

func New() Hash {
	return &fmap{
		size:     defaultSize,
		capacity: 1 << defaultSize,
	}
}

func (m *fmap) fibonacciIndex(k uint64) uint {
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
