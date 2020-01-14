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
	freenet.NodeChannelCapacity = 100 // ?
	freenet.NodeTableCapacity = 250   // 5.1 pg.12
	freenet.NodeFileCapacity = 50     // 5.1 pg.12
	freenet.NodeJobTimeout = 1        // ? seconds
	freenet.NodeJobCapacity = 100     // ?
	freenet.HopsToLiveDefault = 20    // 5.1 pg.13
}

func main() {

	var InitialNodeCount uint32 = 20   // 20
	var ActionsPerTimestep int = 2     // ?
	var SimulationDuration int = 40000 // 200,000/5 = 40,000
	var SnapshotFreq int = 250         // ? (every 50 node additions)
	var SnapshotReqCount int = 50      // ? 300

	var currNodeCount uint32 = InitialNodeCount

	// Slice containing all nodes
	var nodes []*freenet.Node

	// Create nodes
	for i := uint32(0); i < InitialNodeCount; i++ {
		nodes = append(nodes, freenet.NewNode(i))
		nodes[i].Start()
	}

	// Create the regular ring-lattice structure
	for i := uint32(0); i < InitialNodeCount; i++ {
		nodes[i].AddRoutingTableEntry(nodes[(i-1+InitialNodeCount)%InitialNodeCount])
		nodes[i].AddRoutingTableEntry(nodes[(i-2+InitialNodeCount)%InitialNodeCount])
		nodes[i].AddRoutingTableEntry(nodes[(i+1+InitialNodeCount)%InitialNodeCount])
		nodes[i].AddRoutingTableEntry(nodes[(i+2+InitialNodeCount)%InitialNodeCount])
		// fmt.Print("Added ", i+1, " nodes\n")
	}

	var FileCount int = -1

	for i := 1; i <= SimulationDuration; i++ {

		// Actions every timestep
		for j := 0; j < ActionsPerTimestep; j++ {
			srcNodeID := rand.Intn(int(currNodeCount))
			if rand.Intn(2) == 0 {
				// Insert file
				FileCount++
				fileDesc := "/files/file" + strconv.Itoa(FileCount)
				nodes[srcNodeID].SendRequestInsert(fileDesc, "New file")
				time.Sleep(1 * time.Millisecond)
			} else if FileCount >= 0 {
				// Retrieve file
				// freenet.HopsToLiveDefault = 500
				fileID := rand.Intn(FileCount)
				fileDesc := "/files/file" + strconv.Itoa(fileID)
				nodes[srcNodeID].SendRequestData(fileDesc)
				// freenet.HopsToLiveDefault = 20
				time.Sleep(1 * time.Millisecond)
			}
		}

		// New node addition every 5 timesteps
		if i%5 == 0 {
			nodes = append(nodes, freenet.NewNode(currNodeCount))
			nodes[currNodeCount].Start()

			freenet.HopsToLiveDefault = 10
			dstID := rand.Intn(int(currNodeCount))
			nodes[currNodeCount].SendRequestJoin(nodes[dstID])
			time.Sleep(1 * time.Millisecond)

			freenet.HopsToLiveDefault = 20
			currNodeCount++
		}

		// Snapshot at fixed intervals
		if i%SnapshotFreq == 0 {
			time.Sleep(10 * time.Millisecond)
			freenet.HopsToLiveDefault = 100
			fmt.Println("Start Snapshot")

			for j := 0; j < SnapshotReqCount; j++ {
				fileID := rand.Intn(FileCount)
				fileDesc := "/files/file" + strconv.Itoa(fileID)
				srcNodeID := rand.Intn(int(currNodeCount))
				nodes[srcNodeID].SendRequestData(fileDesc)
				time.Sleep(10 * time.Millisecond)
			}

			fmt.Println("End Snapshot")
			freenet.HopsToLiveDefault = 20
		}
	}

	// // Print node states
	// time.Sleep(2 * time.Second)
	// for _, n := range nodes {
	// 	n.Print()
	// }

	return
}
