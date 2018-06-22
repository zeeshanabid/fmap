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

type fmap struct {
}

func (m *fmap) Put(k interface{}, v interface{}) (interface{}, error) {
	return nil, nil
}

func (m *fmap) Get(k interface{}) (interface{}, error) {
	return nil, nil
}

func New() Hash {
	return &fmap{}
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
