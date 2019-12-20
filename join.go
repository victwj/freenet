// Functions related to joining freenet
package main

import (
	"fmt"
)

func (n *node) joinHandler(msg nodeMsg) {
	fmt.Printf("Node %d received join from node %d\n", n.id, msg.from.id)
}

func (n *node) sendJoinRequest(dst *node) {
	msg := n.newNodeMsg(joinMsgType, "Test Message")
	n.send(msg, dst)
}
