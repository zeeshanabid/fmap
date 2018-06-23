package fmap

import (
	"testing"
)

func TestHashCode(t *testing.T) {
	code := hashCode(nil)
	expected := uint64(14695981039346656037)

	if code != expected {
		t.Errorf("Expected hashcode %d, got %d", expected, code)
	}

	b, err := getBytes("hello!")
	if err != nil {
		t.Errorf("Cannot get bytes: %s", err)
	}
	code = hashCode(b)
	expected = uint64(18402282655386536633)

	if code != expected {
		t.Errorf("Expected hashcode %d, got %d", expected, code)
	}
}

func TestNewMap(t *testing.T) {
	m := New()
	hm, ok := m.(*fmap)

	if !ok {
		t.Errorf("Expected type fmap %t, got %t", true, ok)
	}

	if hm.size != defaultSize {
		t.Errorf("Expected size %d, got %d", defaultSize, hm.size)
	}

	if hm.Length() != 0 {
		t.Errorf("Expected length %d, got %d", 0, hm.Length())
	}

	if len(hm.Keys()) != 0 {
		t.Errorf("Expected keys length %d, got %d", 0, len(hm.Keys()))
	}

	if len(hm.Values()) != 0 {
		t.Errorf("Expected values length %d, got %d", 0, len(hm.Values()))
	}
}
