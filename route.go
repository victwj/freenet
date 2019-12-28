// Functions related to routing and handling messages
package main

import (
	"container/heap"
	"log"
)

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
	} else if msgType == RequestDataMsgType {
		n.serveRequestData(msg)
	}
}

// Add entry to the routing table
func (n *node) addRoutingTableEntry(key string, nodeEntry *node) {
	if nodeEntry == n {
		panic("Error: adding self to routing table")
	}
	n.table.Add(key, nodeEntry)
}

func (n *node) getRoutingTableEntry(key string) *node {
	result, found := n.table.Get(key)
	if found {
		return result.(*node)
	}
	return nil
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

// Get the n-th match of the routing table, given a string to match
func (n *node) getRouteMatch(match string, routeNum int) *node {

	// Sanity check
	if routeNum == 0 {
		panic("routeNum is zero")
	}

	// Match immediately
	if routeNum == 1 && n.table.Contains(match) {
		result, _ := n.table.Get(match)
		return result.(*node)
	}
	// Return nil immediately
	if routeNum > n.table.Len() {
		return nil
	}

	// Calculate all string similarities and put in a PQ
	pq := make(PriorityQueue, n.table.Len())
	for i, key := range n.table.Keys() {
		keyStr := key.(string)
		pq[i] = &Item{
			value:    keyStr,
			priority: stringSimilarity(match, keyStr),
			index:    i,
		}
	}
	heap.Init(&pq)

	// Pop the PQ routeNum number of times
	keyResult := ""
	for routeNum > 0 {
		keyResult = heap.Pop(&pq).(*Item).value
		routeNum--
	}
	return n.getRoutingTableEntry(keyResult)
}

// No need to do fancy things like levenshtein as long as consistent
// Count number of equivalent characters
func stringSimilarity(s1 string, s2 string) int {
	min := len(s1)
	if len(s1) > len(s2) {
		min = len(s2)
	}
	ctr := 0
	for i := 0; i < min; i++ {
		if s1[i] == s2[i] {
			ctr += 1
		}
	}
	return ctr
}