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

	// Initialize freenet with 5 nodes
	for i := uint32(0); i < 5; i++ {
		nodes = append(nodes, newNode(i))
		nodes[i].start()
	}

	// // Each node has everyone else in its routing table
	// for i := uint32(0); i < uint32(len(nodes)); i++ {
	// 	ctr := 0
	// 	for j := uint32(0); j < uint32(len(nodes)); j++ {
	// 		if i != j {
	// 			nodes[i].table.Add(nodes[j])
	// 			ctr++
	// 		}
	// 	}
	// }

	// Test message handling
	n1 := newNode(5)
	n1.sendJoinRequest(nodes[0])

	// Wait a little to let nodes log
	time.Sleep(2 * time.Second)
}
