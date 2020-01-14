package main

import (
	"fmt"
	"math/rand"
	"strconv"
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

func init() {
	freenet.NodeChannelCapacity = 100 // ? 5
	freenet.NodeTableCapacity = 250   // 5.1 pg.12
	freenet.NodeFileCapacity = 50     // 5.1 pg.12
	freenet.NodeJobTimeout = 3        // ? seconds
	freenet.NodeJobCapacity = 100     // ? 10
	freenet.HopsToLiveDefault = 20    // 5.1 pg.13
}

func main() {

	var NodeCount uint32 = 1000       // 1000
	var ActionsPerTimestep int = 2    // ? 10
	var SimulationDuration int = 2000 // 5000
	var SnapshotReqCount int = 300    // 300

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
		nodes[i].AddRoutingTableEntry(nodes[(i-2+NodeCount)%NodeCount])
		nodes[i].AddRoutingTableEntry(nodes[(i+1+NodeCount)%NodeCount])
		nodes[i].AddRoutingTableEntry(nodes[(i+2+NodeCount)%NodeCount])
		// fmt.Print("Added ", i+1, " nodes\n")
	}

	// printNodeStates(nodes)

	var FileCount int = -1

	/* // Test code STARTS

	for i := 1; i <= SimulationDuration; i++ {
		srcNodeID := rand.Intn(int(NodeCount))
		// Insert file
		FileCount++
		fileDesc := "/files/file" + strconv.Itoa(FileCount)
		nodes[srcNodeID].SendRequestInsert(fileDesc, "New file")
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)

	for i := 1; i <= SimulationDuration; i++ {
		srcNodeID := rand.Intn(int(NodeCount))
		fileIndex := rand.Intn(int(FileCount))
		// Request file
		fileDesc := "/files/file" + strconv.Itoa(fileIndex)
		nodes[srcNodeID].SendRequestData(fileDesc)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("\n\nStart snapshot")
	freenet.HopsToLiveDefault = 500

	for i := 1; i <= 50; i++ {
		srcNodeID := rand.Intn(int(NodeCount))
		fileIndex := rand.Intn(int(FileCount))
		// Request file
		fileDesc := "/files/file" + strconv.Itoa(fileIndex)
		nodes[srcNodeID].SendRequestData(fileDesc)
		time.Sleep(100 * time.Millisecond)
	}

	// Test code ENDS */

	for i := 1; i <= SimulationDuration; i++ {

		// Actions every timestep
		for j := 0; j < ActionsPerTimestep; j++ {
			srcNodeID := rand.Intn(int(NodeCount))
			if rand.Intn(2) == 0 {
				// Insert file
				FileCount++
				fileDesc := "/files/file" + strconv.Itoa(FileCount)
				// fmt.Println("Insert: ", fileDesc)
				nodes[srcNodeID].SendRequestInsert(fileDesc, "New file")
				time.Sleep(1 * time.Millisecond)
			} else if FileCount >= 0 {
				// Retrieve file
				fileID := rand.Intn(FileCount)
				fileDesc := "/files/file" + strconv.Itoa(fileID)
				// fmt.Println("Retrieve: ", fileDesc)
				nodes[srcNodeID].SendRequestData(fileDesc)
				time.Sleep(1 * time.Millisecond)
			}
		}

		// Snapshots every 100 timesteps
		if i%100 == 0 {
			time.Sleep(100 * time.Millisecond)
			freenet.HopsToLiveDefault = 500
			fmt.Println("Start Snapshot")

			for j := 0; j < SnapshotReqCount; j++ {
				fileID := rand.Intn(FileCount)
				fileDesc := "/files/file" + strconv.Itoa(fileID)
				srcNodeID := rand.Intn(int(NodeCount))
				nodes[srcNodeID].SendRequestData(fileDesc)
				time.Sleep(1 * time.Millisecond)
			}

			fmt.Println("End Snapshot")
			freenet.HopsToLiveDefault = 20
		}
	}

	// // Print node states
	// printNodeStates(nodes)

	return
}
