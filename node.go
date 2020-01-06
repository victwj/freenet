// Functions implementing a node in freenet
package freenet

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	lru "github.com/hashicorp/golang-lru"
	cache "github.com/patrickmn/go-cache"
)

const (
	// Configuration constants
	nodeChannelCapacity = 5 // TODO: Decide on a value
	nodeTableCapacity   = 5 // 250 // 5.1 pg.12
	nodeFileCapacity    = 5 // 50 // 5.1 pg.12
	nodeJobTimeout      = 5
	nodeJobCapacity     = 10
	hopsToLiveDefault   = 3 // 20  // 5.1 pg.13
)

// Node data structure in freenet
type Node struct {
	id        uint32         // Unique ID per node
	ch        chan nodeMsg   // The "IP/port" of the node
	table     *lru.Cache     // Routing table, string->*node
	disk      *lru.Cache     // Files stored in "disk", string->string
	processor *nodeProcessor // Cache with timeout, stores pending msg IDs
}

// Messages sent by nodes
/* Not implemented:
- Finite forwarding probability of HTL/Depth == 1
- Obfuscating depth by setting it randomly
*/

type nodeMsg struct {
	msgType uint8  // Type of message, see constants
	msgID   uint64 // Unique ID, per transaction
	htl     int    // Hops to live
	depth   int    // To let packets backtrack successfully
	from    *Node  // Pointer to Node which sent this msg
	origin  *Node  // The first Node which started this transaction
	body    string // String body, depends on msg type
}

// Wrapper around the timeout cache storing pending jobs
// Need to wrap since there is no way to limit the size without it
// Safer this way
type nodeProcessor struct {
	jobs     *cache.Cache // Stores msgID->*nodeJob
	capacity int
}

// The data type of a pending job, stored in nodeProcessor
// Save space, instead of storing an entire nodeMsg
type nodeJob struct {
	from     *Node // Who sent this job to us
	origin   *Node // The origin of this job
	routeNum int
	// E.g. if routeNum == 1, we want to use the second match (0-indexed)
	// This means the first match was previously unsuccessful
}

// String representation of a node for logging/debugging
func (n Node) String() string {
	return fmt.Sprint("Node", n.id)
}

// String conversion for logging
func (m nodeMsg) String() string {
	return fmt.Sprintf("(MsgID: %d, From: %d, Type: %d, HTL: %d, Depth: %d, Body: %s)", m.msgID, m.from.id, m.msgType, m.htl, m.depth, m.body)
}

// Returns a pointer to an initialized node with the given ID
func NewNode(id uint32) *Node {
	n := new(Node)
	n.id = id
	n.ch = make(chan nodeMsg, nodeChannelCapacity)
	n.table, _ = lru.New(nodeTableCapacity)
	n.disk, _ = lru.New(nodeFileCapacity)

	// Initialize processor
	n.processor = new(nodeProcessor)
	n.processor.jobs = cache.New(nodeJobTimeout*time.Second, (nodeJobTimeout+1)*time.Second)
	n.processor.capacity = nodeJobCapacity
	return n
}

// Factory for Node messages
// Member function of Node, since we need a reference to sender
// Don't return pointer since we never really work with pointer to msg
func (n *Node) newNodeMsg(msgType uint8, body string) nodeMsg {
	m := new(nodeMsg)
	m.msgType = msgType
	m.msgID = rand.Uint64() // Random number for msg ID
	m.htl = hopsToLiveDefault
	m.from = n
	m.body = body
	m.depth = 0
	m.origin = nil // Don't set if not necessary, safer
	return *m
}

// Core functions of a Node, emulating primitive operations

// Spawn a goroutine which will handle/route messages sent to this node
func (n *Node) Start() {
	go n.listen()
}

// This node will no longer handle/route any message
func (n *Node) Stop() {
	close(n.ch)
}

func (n *Node) listen() {
	// Keep listening until the channel is closed
	for msg := range n.ch {
		n.route(msg)
	}
}

func (n *Node) send(msg nodeMsg, dst *Node) {
	if n == dst {
		panic("Sending a message to self")
	}
	// We never want to forward the wrong from field
	msg.from = n
	dst.ch <- msg
}

// Adds job to process
// The job cache is a map of msgID/xactID -> *nodeJob
// Return true if success, false otherwise
func (n *Node) addJob(msg nodeMsg) bool {
	// Processor is full
	if n.processor.jobs.ItemCount() >= n.processor.capacity {
		return false
	}

	msgID := strconv.FormatUint(msg.msgID, 10)

	// Job is in the processor but re adding it, error
	if n.hasJob(msg) {
		fmt.Println("Error")
		panic("Re-adding a job")
	}

	// Create job
	job := new(nodeJob)
	job.from = msg.from
	job.routeNum = 0

	if msg.origin != nil {
		job.origin = msg.origin
	}

	// Add to processor
	n.processor.jobs.SetDefault(msgID, job)

	return true
}

// Check if this job exists
func (n *Node) hasJob(msg nodeMsg) bool {
	msgID := strconv.FormatUint(msg.msgID, 10)
	_, found := n.processor.jobs.Get(msgID)
	return found
}

// If job exists in processor, return the nodeJob, increment routeNum
// If it doesn't exist, return nil
// !!! Note: This getter changes state
func (n *Node) getJob(msg nodeMsg) *nodeJob {
	msgID := strconv.FormatUint(msg.msgID, 10)

	// Check if this job exists in processor
	val, found := n.processor.jobs.Get(msgID)
	if found {
		job := val.(*nodeJob)
		// Increment the routeNum
		job.routeNum += 1
		n.processor.jobs.SetDefault(msgID, job)
		return job
	}
	return nil
}

func (n *Node) deleteJob(msg nodeMsg) {
	msgID := strconv.FormatUint(msg.msgID, 10)
	_, found := n.processor.jobs.Get(msgID)
	if found {
		n.processor.jobs.Delete(msgID)
	}
}

func (n *Node) Print() {
	fmt.Println(
		"Node ", n.id,
		"\n  Table:", n.table.Keys(),
		"\n  Disk:", n.disk.Keys(),
		"\n  Jobs:", n.processor.jobs.Items())
}
