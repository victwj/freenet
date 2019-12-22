package main

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	// Configuration constants
	nodeChannelCapacity = 5 // TODO: Decide on a value
	nodeTableCapacity   = 5 // TODO: Decide on a value
	nodeFileCapacity    = 5 // TODO: Decide on a value
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
	n.files = make([]string, nodeFileCapacity)
	return n
}

// String conversion for logging
func (n node) String() string {
	return fmt.Sprintf("Node %d", n.id)
}

// String conversion for logging
func (m nodeMsg) String() string {
	return fmt.Sprintf("(MsgID: %d, From: %d, Type: %d, HTL: %d, Body: %s)", m.msgID, m.from.id, m.msgType, m.htl, m.body)
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
		msgType := msg.msgType

		// Act based on message type
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
