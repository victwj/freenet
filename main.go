package main

import (
	"fmt"
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

	// For now, testing things by adding things manually

	// Add routing table entries
	nodes[0].addRoutingTableEntry("cbb", nodes[1])
	nodes[0].addRoutingTableEntry("cbbb", nodes[3])
	nodes[1].addRoutingTableEntry("testkey", nodes[2])

	// Add a file
	nodes[2].addFileDescr("/existing/file", "hello world")
	// "/existing/file KSK is cbbb589"

	// Wait a little to let nodes stabilize
	time.Sleep(1 * time.Second)

	// Send a data request
	nodes[0].sendRequestData("/nonexistent/file")
	nodes[0].sendRequestData("/existing/file")

	// Wait a little to let nodes log
	time.Sleep(2 * time.Second)

	// Print final state of nodes
	fmt.Println("\nFinal node states:")
	for _, n := range nodes {
		fmt.Println(n,
			"\n  Table:", n.table.Keys(),
			"\n  Disk:", n.disk.Keys(),
			"\n  Jobs:", n.processor.jobs.Items())
	}
}
