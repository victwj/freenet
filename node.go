package main

import (
	"fmt"
	"log"
	"math/rand"
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
	nodeJobCapacity     = 5
	hopsToLiveDefault   = 5 // 20  // 5.1 pg.13

	// Node message types
	failMsgType = 0
	joinMsgType = 1

	requestInsertMsgType   = 10
	requestDataMsgType     = 11
	requestContinueMsgType = 12

	replyInsertMsgType   = 20
	replyNotFoundMsgType = 21
	replyRestartMsgType  = 22

	sendDataMsgType   = 30
	sendInsertMsgType = 31

	/* Message types from paper:
	Request.Data = request file
	Reply.Restart = tell nodes to extend timeout
	Send.Data = file found, sending back
	Reply.NotFound = file not found
	Request.Continue = if file not found, but there is HTL remaining
	Request.Insert = file insert
	Reply.Insert = insert can go ahead
	Send.Insert = contains the data
	*/
)

// Freenet node
type node struct {
	id        uint32        // Unique ID per node
	ch        chan nodeMsg  // The "IP/port" of the node
	table     *lru.Cache    // Routing table
	disk      *lru.Cache    // Files stored in "disk"
	processor nodeProcessor // Cache with timeout, stores pending msg IDs
	// If packets can drop, we need pending jobs to time out
	// But packets will never be dropped in current implementation
	// If channel is full, sender will block
	// But, probably deadlock potential here..
	// Deferring for now
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
	jobs     *cache.Cache
	capacity int
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
		log.Println(n, "received", msg)

		// Hops to live too low
		if msg.htl <= 0 {
			failMsg := n.newNodeMsg(failMsgType, "")
			n.send(failMsg, msg.from)
		}

		// Decrement HTL
		msg.htl -= 1
		msg.depth += 1
		msgType := msg.msgType

		// Act based on message type, call handlers
		if msgType == failMsgType {

		} else if msgType == joinMsgType {
			n.joinHandler(msg)
		}
	}

	log.Println(n, "done")
}

func (n *node) send(msg nodeMsg, dst *node) {
	dst.ch <- msg
}

// Adds job to process
// The job cache is a map of msgID/xactID -> routing table number
// e.g. if routing table number is 2, will route to the second match
func (n *node) addJob(msgID int, tableNum int) bool {
	// Processor is full
	if n.processor.jobs.ItemCount() >= n.processor.capacity {
		return false
	}
	n.processor.jobs.Set(string(msgID), tableNum, cache.DefaultExpiration)
	return true
}
