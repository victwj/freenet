package main

import (
	"fmt"
	"math/rand"
)

const (
	// Configuration constants
	nodeChannelCapacity = 5 // TODO: Decide on a value
	nodeTableCapacity   = 5 // TODO: Decide on a value
	hopsToLiveDefault   = 5 // TODO: Decide on a value

	// Node message types
	failMsgType = 0
	joinMsgType = 1
)

// Freenet node
type node struct {
	id    uint32
	ch    chan nodeMsg
	table []*node
	files []string
}

// Messages sent by nodes
type nodeMsg struct {
	msgType uint8
	msgID   uint32
	htl     int
	from    *node
	body    string
}

// Factory function, Golang doesn't have constructors
func newNode(id uint32) *node {
	n := new(node)
	n.id = id
	n.ch = make(chan nodeMsg, nodeChannelCapacity)
	n.table = make([]*node, nodeTableCapacity)
	return n
}

// Factory for node messages
// Member function of node, since we need a reference to sender
// Don't return pointer since we never really work with pointer to msg
func (n *node) newNodeMsg(msgType uint8, body string) nodeMsg {
	m := new(nodeMsg)
	m.msgType = msgType
	m.msgID = rand.Uint32() // Random number for msg ID
	m.htl = hopsToLiveDefault
	m.from = n
	m.body = body
	return *m
}

// Core functions of a node, emulating primitive operations

func (n *node) start() {
	fmt.Printf("Node %d started\n", n.id)
	go n.listen()
}

func (n *node) stop() {
	close(n.ch)
}

func (n *node) listen() {
	fmt.Printf("Node %d listening\n", n.id)

	// Keep listening until the channel is closed
	for msg := range n.ch {
		// Hops to live too low
		if msg.htl <= 0 {
			failMsg := n.newNodeMsg(failMsgType, "")
			n.send(failMsg, msg.from.ch)
		}

		// Decrement HTL
		msg.htl -= 1
		msgType := msg.msgType

		// Act based on message type
		if msgType == failMsgType {

		} else if msgType == joinMsgType {
			n.joinHandler(msg)
		}
	}

	fmt.Printf("Node %d done\n", n.id)
}

func (n *node) send(msg nodeMsg, dst chan<- nodeMsg) {
	dst <- msg
}
