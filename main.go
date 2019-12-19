package main

import (
	"time"
)

/*
How to run:
$ go build
$ ./freenet
*/
func main() {
	// Slice containing all nodes
	var nodes []*node

	// Create and start the first freenet node
	nodes = append(nodes, newNode(0))
	nodes[0].start()

	// Simple testing
	n1 := newNode(1)
	n1.SendJoinRequest(nodes[0].ch)

	// Wait a little to let nodes log
	time.Sleep(2 * time.Second)
}
