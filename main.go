package main

import (
	"fmt"
	"time"
)

func printNodeStates(nodes []*node) {
	fmt.Println("\nFinal node states:")
	for _, n := range nodes {
		fmt.Println(n,
			"\n  Table:", n.table.Keys(),
			"\n  Disk:", n.disk.Keys(),
			"\n  Jobs:", n.processor.jobs.Items())
	}

}

func testBasic() {
	// Slice containing all nodes
	var nodes []*node

	// Create 5 nodes
	for i := uint32(0); i < 5; i++ {
		nodes = append(nodes, newNode(i))
		nodes[i].start()
	}

	// The first node in freenet
	nodeZero := nodes[0]
	time.Sleep(1 * time.Second)

	nodes[1].sendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[2].sendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[3].sendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[4].sendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[0].sendRequestInsert("/my/file/hello.txt", "hello world")
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

}

func main() {
	testBasic()
	return
}
