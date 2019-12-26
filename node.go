package main

import (
	"fmt"
	"log"
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
	hopsToLiveDefault   = 5 // 20  // 5.1 pg.13
)

// Freenet node
type node struct {
	id        uint32         // Unique ID per node
	ch        chan nodeMsg   // The "IP/port" of the node
	table     *lru.Cache     // Routing table
	disk      *lru.Cache     // Files stored in "disk"
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
	from    *node  // Pointer to node which sent this msg
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
	from     *node // Who sent this job
	routeNum int
	// E.g. if routeNum == 1, we want to use the second match (0-indexed)
	// This means the first match was previously unsuccessful
}

// String conversion for logging
func (n node) String() string {
	return fmt.Sprintf("Node %d", n.id)
}

// String conversion for logging
func (m nodeMsg) String() string {
	return fmt.Sprintf("(MsgID: %d, From: %d, Type: %d, HTL: %d, Depth: %d, Body: %s)", m.msgID, m.from.id, m.msgType, m.htl, m.depth, m.body)
}

// Factory function, Golang doesn't have constructors
func newNode(id uint32) *node {
	n := new(node)
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

// Factory for node messages
// Member function of node, since we need a reference to sender
// Don't return pointer since we never really work with pointer to msg
func (n *node) newNodeMsg(msgType uint8, body string) nodeMsg {
	m := new(nodeMsg)
	m.msgType = msgType
	m.msgID = rand.Uint64() // Random number for msg ID
	m.htl = hopsToLiveDefault
	m.from = n
	m.body = body
	m.depth = 0
	return *m
}

// Core functions of a node, emulating primitive operations

func (n *node) start() {
	log.Println(n, "started")
	go n.listen()
}

func (n *node) stop() {
	close(n.ch)
}

func (n *node) listen() {
	log.Println(n, "listening")

	// Keep listening until the channel is closed
	for msg := range n.ch {
		n.route(msg)
	}

	log.Println(n, "done")
}

func (n *node) send(msg nodeMsg, dst *node) {
	dst.ch <- msg
}

// Adds job to process
// The job cache is a map of msgID/xactID -> *nodeJob
func (n *node) addJob(msg nodeMsg) *nodeJob {
	// Processor is full
	if n.processor.jobs.ItemCount() >= n.processor.capacity {
		return nil
	}

	// Create job
	job := new(nodeJob)
	job.from = msg.from
	job.routeNum = 0

	// Add to processor
	msgID := strconv.FormatUint(msg.msgID, 10)
	n.processor.jobs.SetDefault(msgID, job)
	return job
}

// If job exists in processor, return the nodeJob, increment routeNum
// If it doesn't exist, return nil
// Note: This getter changes state
func (n *node) getJob(msg nodeMsg) *nodeJob {
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

func (n *node) deleteJob(msg nodeMsg) {
	msgID := strconv.FormatUint(msg.msgID, 10)
	n.processor.jobs.Delete(msgID)
}
