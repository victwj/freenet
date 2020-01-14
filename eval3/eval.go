package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/victwj/freenet"
)

func init() {
	freenet.NodeChannelCapacity = 100 // ?
	freenet.NodeTableCapacity = 250   // 5.1 pg.12
	freenet.NodeFileCapacity = 50     // 5.1 pg.12
	freenet.NodeJobTimeout = 1        // ? seconds
	freenet.NodeJobCapacity = 100     // ?
	freenet.HopsToLiveDefault = 20    // 5.1 pg.13
}

func main() {

	var InitialNodeCount uint32 = 20  // 20
	var MaxNodeCount uint32 = 1000    // 1000
	var ActionsPerTimestep int = 2    // ?
	var SimulationDuration int = 9000 // 8600
	// var SnapshotFreq int = 250        // ? (every 50 node additions)
	// var SnapshotReqCount int = 50     // ? 300

	var currNodeCount uint32 = uint32(InitialNodeCount)
	var maxReached bool = false

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

		if maxReached && currNodeCount <= uint32(0.25*float32(MaxNodeCount)) {
			break
		}

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

		// New node addition or failure (after max nodes reached) every 5 timesteps
		if i%5 == 0 {

			if !maxReached {
				nodes = append(nodes, freenet.NewNode(currNodeCount))
				nodes[currNodeCount].Start()

				freenet.HopsToLiveDefault = 10
				dstID := rand.Intn(int(currNodeCount))
				nodes[currNodeCount].SendRequestJoin(nodes[dstID])
				time.Sleep(1 * time.Millisecond)

				freenet.HopsToLiveDefault = 20
				currNodeCount++

				if currNodeCount == MaxNodeCount {
					maxReached = true
					currNodeCount--
					// fmt.Println("Node count reduced to", currNodeCount)
				}

			} else {
				// delID := rand.Intn(int(currNodeCount))
				delID := currNodeCount
				nodes[delID].Stop()
				currNodeCount--

				// fmt.Println("Node count reduced to", currNodeCount)
			}
		}

		// // Snapshot at fixed intervals
		// if maxReached && i%SnapshotFreq == 0 {
		// 	time.Sleep(10 * time.Millisecond)
		// 	freenet.HopsToLiveDefault = 100
		// 	fmt.Println("Start Snapshot")

		// 	for j := 0; j < SnapshotReqCount; j++ {
		// 		fileID := rand.Intn(FileCount)
		// 		fileDesc := "/files/file" + strconv.Itoa(fileID)
		// 		srcNodeID := rand.Intn(int(currNodeCount))
		// 		nodes[srcNodeID].SendRequestData(fileDesc)
		// 		time.Sleep(10 * time.Millisecond)
		// 	}

		// 	fmt.Println("End Snapshot")
		// 	freenet.HopsToLiveDefault = 20
		// }
	}

	return
}
