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

	var InitialNodeCount uint32 = 20  // 20
	var MaxNodeCount uint32 = 1000    // 1000
	var ActionsPerTimestep int = 10   // ?
	var SimulationDuration int = 2000 // 200,000/5 = 40,000

	var currNodeCount uint32 = InitialNodeCount
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
		time.Sleep(1 * time.Millisecond)
		nodes[i].AddRoutingTableEntry(nodes[(i-2+InitialNodeCount)%InitialNodeCount])
		time.Sleep(1 * time.Millisecond)
		nodes[i].AddRoutingTableEntry(nodes[(i+1+InitialNodeCount)%InitialNodeCount])
		time.Sleep(1 * time.Millisecond)
		nodes[i].AddRoutingTableEntry(nodes[(i+2+InitialNodeCount)%InitialNodeCount])
		time.Sleep(1 * time.Millisecond)
		// fmt.Print("Added ", i+1, " nodes\n")
	}

	// // Print node states
	// time.Sleep(2 * time.Second)
	// for _, n := range nodes {
	// 	n.Print()
	// }

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
				fileDesc := "files/file" + strconv.Itoa(FileCount)
				// fmt.Println("Insert: ", fileDesc)
				nodes[srcNodeID].SendRequestInsert(fileDesc, "Inserted new file")
			} else if FileCount >= 0 {
				// Retrieve file
				freenet.HopsToLiveDefault = 500
				fileID := rand.Intn(FileCount)
				fileDesc := "files/file" + strconv.Itoa(fileID)
				// fmt.Println("Retrieve: ", fileDesc)
				nodes[srcNodeID].SendRequestData(fileDesc)
				freenet.HopsToLiveDefault = 20
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
					fmt.Println("Node count reduced to", currNodeCount)
				}

			} else {
				delID := rand.Intn(int(currNodeCount))
				nodes[delID].Stop()
				currNodeCount--

				fmt.Println("Node count reduced to", currNodeCount)
			}
		}
	}

	// // Print node states
	// time.Sleep(2 * time.Second)
	// for _, n := range nodes {
	// 	n.Print()
	// }

	return
}