package main

import (
	"fmt"
)

const (
	// Configuration constants
	nodeChannelCapacity = 5 // TODO: Decide on a value
	nodeTableCapacity   = 5 // TODO: Decide on a value

	// Node message type codes
	joinMsgType = 0
)

// Freenet node
type node struct {
	id    uint32
	ch    chan nodeMsg
	table []*node
}

// Messages sent by nodes
type nodeMsg struct {
	msgType uint8
	msgID   uint32
	from    uint32
	to      uint32
	htl     int
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
		fmt.Printf("Node %d received: %s\n", n.id, msg.body)
		// TODO: parse the message, act based on message type
		// if messageType == joinMessage, ...
	}

	fmt.Printf("Node %d done\n", n.id)
}

func (n *node) send(msg nodeMsg, dst chan<- nodeMsg) {
	dst <- msg
}

// Freenet-specific functions, built on top of primitive ops

func (n *node) sendJoinRequest(dst chan<- nodeMsg) {
	var msg nodeMsg
	msg.body = fmt.Sprintf("Test join message from node %d", n.id)
	n.send(msg, dst)
}
