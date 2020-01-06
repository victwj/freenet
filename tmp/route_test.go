package freenet

import (
	"testing"
)

func TestStringSimilarity(t *testing.T) {
	s1 := "abc"
	s2 := "def"
	s3 := "aef"
	s4 := "abdc"
	s5 := "abcc"

	if stringSimilarity(s1, s2) != 0 {
		t.Error("Error with comparing", s1, s2)
	}
	if stringSimilarity(s1, s3) != 1 {
		t.Error("Error with comparing", s1, s3)
	}
	if stringSimilarity(s2, s3) != 2 {
		t.Error("Error with comparing", s2, s3)
	}
	if stringSimilarity(s1, s1) != 3 {
		t.Error("Error with comparing", s1, s1)
	}
	if stringSimilarity(s1, s4) != 2 {
		t.Error("Error with comparing", s1, s4)
	}
	if stringSimilarity(s1, s5) != 3 {
		t.Error("Error with comparing", s1, s5)
	}
}

// Don't test situations involving ties
func TestGetRouteMatch(t *testing.T) {
	n := newNode(0)
	n1 := newNode(1)
	n2 := newNode(2)
	n3 := newNode(3)

	n.addRoutingTableEntry("abc", n1)
	n.addRoutingTableEntry("def", n2)
	n.addRoutingTableEntry("acd", n3)

	result := n.getRoutingTableEntry("abc", 1) // should match n1
	if result != n1 {
		t.Error()
	}

	result = n.getRoutingTableEntry("abc", 2) // should match n3
	if result != n3 {
		t.Error()
	}

	result = n.getRoutingTableEntry("abc", 3) // should match n2
	if result != n2 {
		t.Error()
	}

	result = n.getRoutingTableEntry("abc", 4) // should be nil
	if result != nil {
		t.Error(result)
	}

	result = n.getRoutingTableEntry("aec", 1) // should be n1
	if result != n1 {
		t.Error(result)
	}

	result = n.getRoutingTableEntry("accfiller", 1) // should be n1
	if result != n1 {
		t.Error(result)
	}

	result = n.getRoutingTableEntry("acc", 2) // should be n3
	if result != n3 {
		t.Error(result)
	}

	result = n.getRoutingTableEntry("acc", 3) // should be n2
	if result != n2 {
		t.Error(result)
	}

	result = n.getRoutingTableEntry("acc", 100) // should be nil
	if result != nil {
		t.Error(result)
	}
}
