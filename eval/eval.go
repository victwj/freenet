package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/victwj/freenet"
)

func init() {
	freenet.NodeChannelCapacity = 5 // ?
	freenet.NodeTableCapacity = 250 // 5.1 pg.12
	freenet.NodeFileCapacity = 50   // 5.1 pg.12
	freenet.NodeJobTimeout = 5      // ? seconds
	freenet.NodeJobCapacity = 10    // ?
	freenet.HopsToLiveDefault = 20  // 5.1 pg.13
}

func main() {

	// Slice containing all nodes
	var nodes []*freenet.Node

	// Create 100 nodes
	for i := uint32(0); i < 100; i++ {
		nodes = append(nodes, freenet.NewNode(i))
		nodes[i].Start()
	}

	// The first node in freenet
	nodeZero := nodes[0]

	// The next 10 nodes join using nodeZero
	for i := uint32(1); i < 11; i++ {
		nodes[i].SendRequestJoin(nodeZero)
		time.Sleep(100 * time.Millisecond)
	}

	// The remaining nodes join a random prev node
	for i := uint32(11); i < 100; i++ {
		dstID := rand.Intn(int(i))
		nodes[i].SendRequestJoin(nodes[dstID])
		time.Sleep(100 * time.Millisecond)
	}

	// Print node states
	time.Sleep(2 * time.Second)
	for _, n := range nodes {
		n.Print()
	}

	// 100 random inserts
	// Can be concurrent, no sleeps
	for i := 0; i < 100; i++ {
		srcID := rand.Intn(100)
		randomFileDescr := strconv.Itoa(srcID) + "salt"
		nodes[srcID].SendRequestInsert(randomFileDescr, "test file")
	}

	// Print node states
	time.Sleep(2 * time.Second)
	for _, n := range nodes {
		n.Print()
	}

	// 100 random gets
	// Can be concurrent, no sleeps
	for i := 0; i < 100; i++ {
		srcID := rand.Intn(100)
		randomFileDescr := strconv.Itoa(srcID) + "salt"
		nodes[srcID].SendRequestData(randomFileDescr)
	}

	// Print node states
	time.Sleep(2 * time.Second)
	for _, n := range nodes {
		n.Print()
	}

	return
}
