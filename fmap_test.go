package fmap

import "testing"

func TestHashCode(t *testing.T) {
	code := hashCode(nil)
	expected := uint64(14695981039346656037)

	if code != expected {
		t.Errorf("Expected hashcode %d, got %d", expected, code)
	}

	code = hashCode([]byte("hello!"))
	expected = uint64(12230837384389815902)

	if code != expected {
		t.Errorf("Expected hashcode %d, got %d", expected, code)
	}
}
