package main

import (
	"fmt"
)

const (
	// Configuration constants
	nodeChannelCapacity = 5 // TODO: Decide on a value

	// Message codes
	joinMessage = 0
)

// We can define our own message type, i.e. another struct not string
// But string should work for now, seems sufficient and easy enough
// Regardless we should define a message schema
type node struct {
	id           uint32
	ch           chan string
	routingTable []*node
}

// A message struct might be like this, unused for now
type msg struct {
	msgType uint8
	msg     string
}

// Factory function, Golang doesn't have constructors
func newNode(id uint32) *node {
	n := new(node)
	n.id = id
	n.ch = make(chan string, nodeChannelCapacity)
	return n
}

// *** Core functions of a node, emulating primitive operations ***

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
		fmt.Printf("Node %d received: %s\n", n.id, msg)
		// TODO: parse the message, act based on message type
		// if messageType == joinMessage, ...
	}
	fmt.Printf("Node %d done\n", n.id)
}

func (n *node) send(msg string, dst chan<- string) {
	dst <- msg
}

// *** Freenet-specific functions, built on top of primitive ops ***

func (n *node) SendJoinRequest(dst chan<- string) {
	msg := fmt.Sprintf("Test join message from node %d", n.id)
	n.send(msg, dst)
}
