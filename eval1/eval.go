package main

import (
	"fmt"
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

	var NodeCount uint32 = 10         // 1000
	var ActionsPerTimestep int = 10   // ?
	var SimulationDuration int = 2000 // 5000

	// Slice containing all nodes
	var nodes []*freenet.Node

	// Create nodes
	for i := uint32(0); i < NodeCount; i++ {
		nodes = append(nodes, freenet.NewNode(i))
		nodes[i].Start()
	}

	// Create the regular ring-lattice structure
	for i := uint32(0); i < NodeCount; i++ {
		nodes[i].AddRoutingTableEntry(nodes[(i-1+NodeCount)%NodeCount])
		time.Sleep(1 * time.Millisecond)
		nodes[i].AddRoutingTableEntry(nodes[(i-2+NodeCount)%NodeCount])
		time.Sleep(1 * time.Millisecond)
		nodes[i].AddRoutingTableEntry(nodes[(i+1+NodeCount)%NodeCount])
		time.Sleep(1 * time.Millisecond)
		nodes[i].AddRoutingTableEntry(nodes[(i+2+NodeCount)%NodeCount])
		time.Sleep(1 * time.Millisecond)
		// fmt.Print("Added ", i+1, " nodes\n")
	}

	// // Print node states
	// time.Sleep(2 * time.Second)
	// for _, n := range nodes {
	// 	n.Print()
	// }

	var FileCount int = -1

	for i := 0; i < SimulationDuration; i++ {

		// Actions every timestep
		for j := 0; j < ActionsPerTimestep; j++ {
			srcNodeID := rand.Intn(int(NodeCount))
			if rand.Intn(2) == 0 {
				// Insert file
				FileCount++
				fileDesc := "files/file" + strconv.Itoa(FileCount)
				// fmt.Println("Insert: ", fileDesc)
				nodes[srcNodeID].SendRequestInsert(fileDesc, "Inserted new file")
			} else if FileCount >= 0 {
				// Retrieve file
				fileID := rand.Intn(FileCount)
				fileDesc := "files/file" + strconv.Itoa(fileID)
				// fmt.Println("Retrieve: ", fileDesc)
				nodes[srcNodeID].SendRequestData(fileDesc)
			}
		}

		// Snapshots every 100 timesteps
		if (i+1)%100 == 0 {
			freenet.HopsToLiveDefault = 500
			fmt.Println("Start Snapshot")
			for j := 0; j < 300; j++ {
				fileID := rand.Intn(FileCount)
				fileDesc := "files/file" + strconv.Itoa(fileID)
				srcNodeID := rand.Intn(int(NodeCount))
				nodes[srcNodeID].SendRequestData(fileDesc)
			}
			fmt.Println("End Snapshot")
			freenet.HopsToLiveDefault = 20
		}
	}

	// // 100 random inserts
	// // Can be concurrent, no sleeps
	// for i := 0; i < 100; i++ {
	// 	srcID := rand.Intn(100)
	// 	randomFileDescr := strconv.Itoa(srcID) + "salt"
	// 	nodes[srcID].SendRequestInsert(randomFileDescr, "test file")
	// }

	// // Print node states
	// time.Sleep(2 * time.Second)
	// for _, n := range nodes {
	// 	n.Print()
	// }

	// // 100 random gets
	// // Can be concurrent, no sleeps
	// for i := 0; i < 100; i++ {
	// 	srcID := rand.Intn(100)
	// 	randomFileDescr := strconv.Itoa(srcID) + "salt"
	// 	nodes[srcID].SendRequestData(randomFileDescr)
	// }

	// // Print node states
	// time.Sleep(2 * time.Second)
	// for _, n := range nodes {
	// 	n.Print()
	// }

	return
}
