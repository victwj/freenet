package main

import (
	"fmt"
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

	// Testing KSK, must all be the same and deterministic
	_, _, ksk := genKeywordSignedKey("/test/test/hello")
	fmt.Println(ksk)
	_, _, ksk = genKeywordSignedKey("/test/test/hello")
	fmt.Println(ksk)
	_, _, ksk = genKeywordSignedKey("/test/test/hello")
	fmt.Println(ksk)

	// Test node processor
	testMsg1 := n1.newNodeMsg(failMsgType, "test msg 1")
	testMsg2 := n1.newNodeMsg(failMsgType, "test msg 2")
	n1.addJob(testMsg1)

	fmt.Println(n1.getJob(testMsg1))
	fmt.Println(n1.getJob(testMsg1))
	fmt.Println(n1.getJob(testMsg1))
	fmt.Println(n1.getJob(testMsg2))

	time.Sleep(8 * time.Second)
	fmt.Println(n1.getJob(testMsg1)) //Should expire

}
