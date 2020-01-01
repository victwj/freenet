// Functions related to joining freenet
package main

import (
	"fmt"
)

func (n *node) sendRequestJoin(dst *node) {
	// Add routing table entry for the node we are sending to
	dstIdStr := fmt.Sprint(dst.id)
	_, _, dstIdKey := genKeywordSignedKey(dstIdStr)
	n.addRoutingTableEntry(dstIdKey, dst)

	// Create request
	nodeIdStr := fmt.Sprint(n.id)
	_, _, nodeIdKey := genKeywordSignedKey(nodeIdStr)
	msg := n.newNodeMsg(RequestJoinMsgType, nodeIdKey)
	msg.origin = n

	// Add the job, proceed if there is processing space
	if n.addJob(msg) {
		// Get the job we just made to get routeNum
		job := n.getJob(msg)
		// Figure out who to send the job to
		dst := n.getRoutingTableEntry(dstIdKey, job.routeNum)
		// If there is a node to send to, send it
		if dst != nil {
			n.send(msg, dst)
		}
	}
}

func (n *node) serveRequestJoin(msg nodeMsg) {

	// If this msg went full circle and came back to us
	// Just send a success back
	if n.hasJob(msg) {
		msg.msgType = ReplyJoinMsgType
		msg.htl = msg.depth
		msg.depth = 0
		n.send(msg, msg.from)
		return
	}

	ksk := msg.body

	// If node doesn't have enough processing capacity to add this job
	// Send fail
	if !n.addJob(msg) {
		msg.msgType = FailMsgType
		msg.htl = msg.depth
		msg.depth = 0
		n.send(msg, msg.from)
		return
	}

	job := n.getJob(msg)

	// If the announcement is about to expire soon
	// Don't forward it, just send a success
	if msg.htl == 0 {
		msg.msgType = ReplyJoinMsgType
		msg.htl = msg.depth
		msg.depth = 0
		n.send(msg, job.from)
		n.deleteJob(msg)

		n.addRoutingTableEntry(msg.body, msg.origin)
		return
	}

	// Else, forward it to a random node
	// Get a random node from the routing table
	dst := n.getRoutingTableEntry(ksk, -1)

	// If there is a node to send to, send it
	if dst != nil {
		n.send(msg, dst)
	} else {
		// If there is no node to send to, send success
		msg.msgType = ReplyJoinMsgType
		msg.htl = msg.depth
		msg.depth = 0
		n.send(msg, job.from)
		n.deleteJob(msg)

		// Add to routing table
		n.addRoutingTableEntry(msg.body, msg.origin)
		return
	}
}

func (n *node) serveReplyJoin(msg nodeMsg) {
	if n.hasJob(msg) {
		job := n.getJob(msg)
		// Only forward if we did not start off this whole request
		if job.from != n {
			n.send(msg, job.from)
			// Add to routing table
			n.addRoutingTableEntry(msg.body, msg.origin)
		}
		n.deleteJob(msg)
	}
}
