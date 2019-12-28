// Functions related to joining freenet
package main

import (
	"log"
)

func (n *node) sendJoinRequest(dst *node) {
	msg := n.newNodeMsg(JoinMsgType, "Test Join")
	n.send(msg, dst)
}

func (n *node) serveJoinRequest(msg nodeMsg) {
	log.Println(n, "handling join", msg)

}
