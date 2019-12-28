package main

import (
	"testing"
)

func TestGetRouteMatch(t *testing.T) {
	n := newNode(0)
	n1 := newNode(1)
	n2 := newNode(2)
	n3 := newNode(3)

	n.addRoutingTableEntry("abc", n1)
	n.addRoutingTableEntry("def", n2)
	n.addRoutingTableEntry("acd", n3)

	result := n.getRouteMatch("abc", 1) // should match n1
	if result != n1 {
		t.Error()
	}

	result = n.getRouteMatch("abc", 2) // should match n3
	if result != n3 {
		t.Error()
	}

	result = n.getRouteMatch("abc", 3) // should match n2
	if result != n2 {
		t.Error()
	}

	result = n.getRouteMatch("abc", 4) // should be nil
	if result != nil {
		t.Error(result)
	}

	result = n.getRouteMatch("aec", 1) // should be either n1 or n3
	// Match n1 since it was added first
	if result != n1 {
		t.Error(result)
	}

	result = n.getRouteMatch("aec", 2) // should be either n1 or n3
	if result != n3 {
		t.Error(result)
	}

	result = n.getRouteMatch("aec", 3) // should be n2
	if result != n2 {
		t.Error(result)
	}

	result = n.getRouteMatch("aec", 5) // should be nil
	if result != nil {
		t.Error(result)
	}
}