package main

import (
	"fmt"
	"time"
)

func testJoin() {

}

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
	nodes[0].addRoutingTableEntry("n1", nodes[1])
	nodes[0].addRoutingTableEntry("n2", nodes[2])
	nodes[1].addRoutingTableEntry("n3", nodes[3])
	nodes[2].addRoutingTableEntry("n4", nodes[4])
	nodes[3].addRoutingTableEntry("n0", nodes[0])
	nodes[4].addRoutingTableEntry("n3", nodes[3])

	// Add a file
	nodes[4].addFileDescr("/existing/file", "hello world")
	// "/existing/file KSK is cbbb589"

	// Wait a little to let nodes stabilize
	time.Sleep(1 * time.Second)

	// Send a data request
	// nodes[0].sendRequestData("/nonexistent/file")
	// nodes[0].sendRequestData("/existing/file")
	// nodes[0].sendRequestInsert("/new/file", "test file")

	// time.Sleep(1 * time.Second)
	// Print final state of nodes
	// fmt.Println("\nFinal node states:")
	// for _, n := range nodes {
	// 	fmt.Println(n,
	// 		"\n  Table:", n.table.Keys(),
	// 		"\n  Disk:", n.disk.Keys(),
	// 		"\n  Jobs:", n.processor.jobs.Items())
	// }
	// nodes[0].sendRequestData("/new/file")
	// nodes[1].sendRequestJoin(nodes[0])

	// Wait a little to let nodes log
	// time.Sleep(2 * time.Second)

	// Print final state of nodes

	nodes[3].sendRequestJoin(nodes[0])

	time.Sleep(3 * time.Second)

	fmt.Println("\nFinal node states:")
	for _, n := range nodes {
		fmt.Println(n,
			"\n  Table:", n.table.Keys(),
			"\n  Disk:", n.disk.Keys(),
			"\n  Jobs:", n.processor.jobs.Items())
	}

}
