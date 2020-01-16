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

	var FinalNodeCount uint32 = 1000 // 1000
	var InitialNodeCount uint32 = 20
	var ActionsPerTimestep int = 2    // ? 10
	var SimulationDuration int = 2000 // 5000
	var SnapshotReqCount int = 300    // 300

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

	for i := uint32(InitialNodeCount); i < FinalNodeCount; i++ {
		nodes = append(nodes, freenet.NewNode(i))
		nodes[i].Start()

		freenet.HopsToLiveDefault = 10
		dstID := rand.Intn(int(i))
		nodes[i].SendRequestJoin(nodes[dstID])
		time.Sleep(1 * time.Millisecond)

		freenet.HopsToLiveDefault = 20
	}

	var FileCount int = -1

	for i := 1; i <= SimulationDuration; i++ {

		// Actions every timestep
		for j := 0; j < ActionsPerTimestep; j++ {

			srcNodeID := rand.Intn(int(FinalNodeCount))

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

		// Snapshot at fixed intervals
		if i%100 == 0 {
			time.Sleep(10 * time.Millisecond)
			freenet.HopsToLiveDefault = 500
			fmt.Println("Start Snapshot")

			for j := 0; j < SnapshotReqCount; j++ {
				fileID := rand.Intn(FileCount)
				fileDesc := "/files/file" + strconv.Itoa(fileID)
				srcNodeID := rand.Intn(int(FinalNodeCount))
				nodes[srcNodeID].SendRequestData(fileDesc)
				time.Sleep(10 * time.Millisecond)
			}

			fmt.Println("End Snapshot")
			freenet.HopsToLiveDefault = 20
		}
	}

	return
}
