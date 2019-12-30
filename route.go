// Functions related to routing and handling messages
package main

import (
	"container/heap"
	"log"
	"sort"
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

	// Temporary: generic fail message to abort the task
	FailMsgType = 0

	// Requests
	RequestInsertMsgType   = 10
	RequestDataMsgType     = 11
	RequestContinueMsgType = 12
	RequestJoinMsgType     = 13

	// Replies
	ReplyInsertMsgType   = 20
	ReplyNotFoundMsgType = 21
	ReplyRestartMsgType  = 22

	// Sends
	SendDataMsgType   = 30
	SendInsertMsgType = 31
)

func (n *node) route(msg nodeMsg) {
	log.Println(n, "received", msg)

	// Decrement HTL
	msg.htl -= 1
	msg.depth += 1
	msgType := msg.msgType

	// Hops to live too low
	if msg.htl < 0 {
		// Special case
		if msg.msgType == RequestInsertMsgType {
			n.serveRequestInsertExpired(msg)
		}
		return
	}

	// Act based on message type, call handlers
	if msgType == FailMsgType {
		n.serveFail(msg)
	} else if msgType == RequestJoinMsgType {
		n.serveRequestJoin(msg)
	} else if msgType == RequestDataMsgType {
		n.serveRequestData(msg)
	} else if msgType == ReplyNotFoundMsgType {
		n.serveReplyNotFound(msg)
	} else if msgType == SendDataMsgType {
		n.serveSendData(msg)
	} else if msgType == RequestInsertMsgType {
		n.serveRequestInsert(msg)
	} else if msgType == ReplyInsertMsgType {
		n.serveReplyInsert(msg)
	}
}

// Add entry to the routing table
func (n *node) addRoutingTableEntry(key string, nodeEntry *node) {
	if nodeEntry == n {
		panic("Error: adding self to routing table")
	}
	n.table.Add(key, nodeEntry)
}

// Get the n-th match of the routing table, given a string to match
func (n *node) getRoutingTableEntry(match string, routeNum int) *node {

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
	tableKeys := n.table.Keys()

	// TODO: Sort will make this operation more stable but inefficient
	// E.g. [A B] both have the same similarity score
	// We match the first one, get A
	// Next time, we get [B A], match the second one (want B),
	// but we get A again since the keys don't have stable order
	// Maybe not worth it, since the routing table can change
	// anytime between jobs anyways and we cannot avoid the
	// above behavior unless we explicitly track teh used keys per job
	sort.Slice(tableKeys, func(i, j int) bool {
		return tableKeys[i].(string) < tableKeys[j].(string)
	})

	for i, key := range tableKeys {
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

	// Return it
	nodeResult, _ := n.table.Get(keyResult)
	return nodeResult.(*node)
}

func (n *node) serveFail(msg nodeMsg) {
	// Get job associated with this message
	// If job has not been seen before or expired
	// Fail message means nothing, drop it
	job := n.getJob(msg)
	if job == nil {
		return
	}

	// If job has been seen and we receive a fail
	// Forward it to the boss of this job
	// If we are the boss of this job, drop it
	if job.from == n {
		log.Print("Deleted job")
	} else {
		n.send(msg, job.from)
	}
	n.deleteJob(msg)
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
