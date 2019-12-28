package main

import (
	"time"
)

func main() {
	// Slice containing all nodes
	var nodes []*node

	// Initialize freenet with 5 nodes
	for i := uint32(0); i < 5; i++ {
		nodes = append(nodes, newNode(i))
		nodes[i].start()
	}

	nodes[0].addRoutingTableEntry("testkey", nodes[1])

	// Wait a little to let nodes stabilize
	time.Sleep(1 * time.Second)

	// Send a data request
	nodes[0].sendRequestData("/nonexistent/file")

	// Wait a little to let nodes log
	time.Sleep(1 * time.Second)
}
