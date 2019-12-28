// Functions related to joining freenet
package main

import (
	"log"
)

func (n *node) sendRequestJoin(dst *node) {
	msg := n.newNodeMsg(RequestJoinMsgType, "Test Join")
	n.send(msg, dst)
}

func (n *node) serveRequestJoin(msg nodeMsg) {
	log.Println(n, "handling join", msg)

}
