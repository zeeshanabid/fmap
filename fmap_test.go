package fmap

import "testing"

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
