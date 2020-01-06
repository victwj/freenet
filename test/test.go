package main

import (
	"fmt"
	"time"

	"github.com/victwj/freenet"
)

func printNodeStates(nodes []*freenet.Node) {
	fmt.Println("\nNode states:")
	for _, n := range nodes {
		n.Print()
	}
	fmt.Println()
}

func testBasic() {
	// Slice containing all nodes
	var nodes []*freenet.Node

	// Create 5 nodes
	for i := uint32(0); i < 5; i++ {
		nodes = append(nodes, freenet.NewNode(i))
		nodes[i].Start()
	}

	// The first node in freenet
	nodeZero := nodes[0]
	time.Sleep(1 * time.Second)

	nodes[1].SendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[2].SendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[3].SendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[4].SendRequestJoin(nodeZero)
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[1].SendRequestInsert("/my/file/hello.txt", "hello world")
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[4].SendRequestData("/my/file/hello.txt")
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[3].SendRequestData("/my/file/hello") // doesn't exist
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[1].SendRequestInsert("/my/file/hello.txt", "testing") // re-inserting
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	nodes[2].SendRequestInsert("/my/file/hello.txt", "testing") // re-inserting
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)

	time.Sleep(2 * time.Second)
	printNodeStates(nodes) // Jobs must expire now

}

func main() {
	testBasic()
	return
}
