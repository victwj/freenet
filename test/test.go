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

	freenet.HopsToLiveDefault = 5
	freenet.NodeJobTimeout = 30

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
	time.Sleep(10 * time.Second)
	printNodeStates(nodes)

	nodes[4].SendRequestData("/my/file/hello.txt")
	time.Sleep(10 * time.Second)
	printNodeStates(nodes)

	nodes[3].SendRequestData("/my/file/hello") // doesn't exist
	time.Sleep(10 * time.Second)
	printNodeStates(nodes)

	nodes[1].SendRequestInsert("/my/file/hello.txt", "testing") // re-inserting
	time.Sleep(10 * time.Second)
	printNodeStates(nodes)

	nodes[2].SendRequestInsert("/my/file/hello.txt", "testing") // re-inserting
	time.Sleep(10 * time.Second)
	printNodeStates(nodes)

	time.Sleep(20 * time.Second)
	printNodeStates(nodes) // Jobs must expire now

}

// Test according to 3.2 pg.7 fig.1
func testPaper() {
	freenet.HopsToLiveDefault = 10
	// Slice containing all nodes
	var nodes []*freenet.Node

	// Create 6 nodes
	for i := uint32(0); i < 6; i++ {
		nodes = append(nodes, freenet.NewNode(i))
		nodes[i].Start()
	}

	nodes[0].AddRoutingTableEntry(nodes[1])
	nodes[1].AddRoutingTableEntry(nodes[2])
	nodes[1].AddRoutingTableEntry(nodes[4])
	nodes[4].AddRoutingTableEntry(nodes[3])
	nodes[4].AddRoutingTableEntry(nodes[5])
	nodes[5].AddRoutingTableEntry(nodes[1])
	nodes[3].AddFile("/test/myfile.txt", "Classified document")

	time.Sleep(1 * time.Second)

	nodes[0].SendRequestData("/test/myfile.txt")
	time.Sleep(1 * time.Second)
	printNodeStates(nodes)
}

func main() {
	testBasic()
	// testPaper()
	return
}
