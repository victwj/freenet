// Functions related to routing and handling messages
package main

import "log"

const (

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

	failMsgType = 0

	// Requests
	requestInsertMsgType   = 10
	requestDataMsgType     = 11
	requestContinueMsgType = 12

	// Replies
	replyInsertMsgType   = 20
	replyNotFoundMsgType = 21
	replyRestartMsgType  = 22

	// Sends
	sendDataMsgType   = 30
	sendInsertMsgType = 31

	// Temp
	joinMsgType = 40
)

func (n *node) route(msg nodeMsg) {
	log.Println(n, "received", msg)

	// Decrement HTL
	msg.htl -= 1
	msg.depth += 1
	msgType := msg.msgType

	// Hops to live too low
	if msg.htl <= 0 {
		failMsg := n.newNodeMsg(failMsgType, "")
		n.send(failMsg, msg.from)
	}

	// Act based on message type, call handlers
	if msgType == failMsgType {

	} else if msgType == joinMsgType {
		n.joinHandler(msg)
	}
}
