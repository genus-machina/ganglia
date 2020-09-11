package ganglia

import (
	"testing"
)

func TestBinaryNameGroup(t *testing.T) {
	group := BinaryNameGroup([]string{"red", "green", "blue"})

	if value := group.Value([]string{}); value != 0 {
		t.Errorf("expected %d but got %d", 0, value)
	}

	if value := group.Value([]string{"red"}); value != 4 {
		t.Errorf("expected %d but got %d", 4, value)
	}

	if value := group.Value([]string{"red", "blue"}); value != 5 {
		t.Errorf("expected %d but got %d", 5, value)
	}

	if value := group.Value([]string{"red", "white", "blue"}); value != 5 {
		t.Errorf("expected %d but got %d", 5, value)
	}
}
