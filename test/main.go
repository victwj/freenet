package main

import (
	"fmt"
	"time"

	"github.com/victwj/freenet"
)

func printNodeStates(nodes []*freenet.Node) {
	fmt.Println("\nFinal node states:")
	for _, n := range nodes {
		n.Print()
	}
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

}

func main() {
	testBasic()
	return
}
