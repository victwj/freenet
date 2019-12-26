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

	FailMsgType = 0

	// Requests
	RequestInsertMsgType   = 10
	RequestDataMsgType     = 11
	RequestContinueMsgType = 12

	// Replies
	ReplyInsertMsgType   = 20
	ReplyNotFoundMsgType = 21
	ReplyRestartMsgType  = 22

	// Sends
	SendDataMsgType   = 30
	SendInsertMsgType = 31

	// Temp
	JoinMsgType = 40
)

func (n *node) route(msg nodeMsg) {
	log.Println(n, "received", msg)

	// Decrement HTL
	msg.htl -= 1
	msg.depth += 1
	msgType := msg.msgType

	// Hops to live too low
	if msg.htl <= 0 {
		// TODO: call routeExpire
		failMsg := n.newNodeMsg(FailMsgType, "")
		n.send(failMsg, msg.from)
	}

	// Act based on message type, call handlers
	if msgType == FailMsgType {

	} else if msgType == JoinMsgType {
		n.joinHandler(msg)
	}
}

func (n *node) routeExpire(msg nodeMsg) {

}

func (n *node) routeFail(msg nodeMsg) {
	// If job has not been seen before or expired
	// Fail message means nothing, drop it
	job := n.getJob(msg)
	if job == nil {
		return
	}

	// If job has been seen and we receive a fail
	// Forward it to the boss of this job
	// If we are the boss of this job, drop it
	if msg.from == n {
		n.deleteJob(msg)
	} else {
		n.send(msg, msg.from)
	}
}
