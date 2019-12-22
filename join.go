// Functions related to joining freenet
package main

import (
	"log"
)

func (n *node) joinHandler(msg nodeMsg) {
	log.Println(n, "handling join", msg)

}

func (n *node) sendJoinRequest(dst *node) {
	msg := n.newNodeMsg(joinMsgType, "Test Join")
	n.send(msg, dst)
}
