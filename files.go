// Functions related to inserting/removing files
package main

import (
	"strings"
)

func (n *node) sendRequestInsert(descr string, file string) {
	_, _, ksk := genKeywordSignedKey(descr)

	// Check self immediately
	fileFound := n.hasFile(ksk)
	if fileFound {
		return
	}

	msg := n.newNodeMsg(RequestInsertMsgType, ksk+" "+file)

	// Set origin when requesting an insert
	msg.origin = n

	// Add the job, proceed if there is processing space
	if n.addJob(msg) {
		// Get the job we just made to get routeNum
		job := n.getJob(msg)
		// Figure out who to send the job to
		dst := n.getRoutingTableEntry(ksk, job.routeNum)
		// If there is a node to send to, send it
		if dst != nil {
			n.send(msg, dst)
		}
	}
}

// Currently prioritize reaching HTL of 0
// Never return success for request insert if HTL is nonzero
func (n *node) serveRequestInsert(msg nodeMsg) {
	// The file key is in the body
	ksk, _ := parseFileFromMsg(msg)
	fileFound := n.hasFile(ksk)

	// If file is found, return it
	if fileFound {
		file := n.getFile(ksk)
		msg.body = ksk + " " + file
		msg.msgType = SendDataMsgType
		msg.htl = msg.depth
		msg.origin = n
		msg.depth = 0
		n.send(msg, msg.from)
		return
	} else {
		if !n.hasJob(msg) && !n.addJob(msg) {
			// TODO: Not specified in paper
			// Node cannot process the insert
			// For now, send a fail
			msg.msgType = FailMsgType
			msg.htl = msg.depth
			msg.depth = 0
			n.send(msg, msg.from)
			return
		}
		job := n.getJob(msg)
		dst := n.getRoutingTableEntry(ksk, job.routeNum)
		// If there is a node to send to, send it
		if dst != nil {
			n.send(msg, dst)
		} else {
			// TODO: Not specified in paper
			// For now send a fail
			msg.msgType = FailMsgType
			msg.htl = msg.depth
			msg.depth = 0
			n.send(msg, job.from)
			n.deleteJob(msg)
			return
		}
	}
}

func (n *node) serveRequestInsertExpired(msg nodeMsg) {
	// Insert request expired means we are good to insert
	// Store the file
	n.addFileFromMsg(msg)

	// Send reply
	msg.htl = msg.depth
	msg.msgType = ReplyInsertMsgType
	msg.depth = 0
	n.send(msg, msg.from)
}

func (n *node) serveReplyInsert(msg nodeMsg) {
	if n.hasJob(msg) {
		// Store the file
		job := n.getJob(msg)
		n.addFileFromMsg(msg)
		// Forward, if not self
		if job.from != n {
			n.send(msg, job.from)
		}
		n.deleteJob(msg)
	}
}

// Sending a request
func (n *node) sendRequestData(descr string) {

	_, _, ksk := genKeywordSignedKey(descr)
	msg := n.newNodeMsg(RequestDataMsgType, ksk)

	// Add the job, proceed if there is processing space
	if n.addJob(msg) {
		// Create the keyword signed key
		_, _, ksk := genKeywordSignedKey(descr)
		// Get the job we just made to get routeNum
		job := n.getJob(msg)
		// Figure out who to send the job to
		dst := n.getRoutingTableEntry(ksk, job.routeNum)
		// If there is a node to send to, send it
		if dst != nil {
			n.send(msg, dst)
		}
	}
}

func (n *node) serveRequestData(msg nodeMsg) {
	// If we get a requestData that we've already seen, refuse
	// Prevent loops
	if n.hasJob(msg) {
		msg.msgType = ReplyNotFoundMsgType
		n.send(msg, msg.from)
	}

	// The file key is in the body
	ksk := msg.body
	fileFound := n.hasFile(ksk)

	// If file is found, return it
	if fileFound {
		file := n.getFile(ksk)
		msg.body = ksk + " " + file
		msg.msgType = SendDataMsgType
		msg.htl = msg.depth
		msg.depth = 0
		msg.origin = n
		n.send(msg, msg.from)
		return
	}

	// File is not found
	// Create the job, but it might fail
	// If we can't create the job, send a not found
	if !n.addJob(msg) {
		msg.msgType = ReplyNotFoundMsgType
		n.send(msg, msg.from)
		return
	}

	// Forward the request for the file since we don't have it
	job := n.getJob(msg)
	dst := n.getRoutingTableEntry(ksk, job.routeNum)
	// If there is a node to send to, send it
	if dst != nil {
		n.send(msg, dst)
	} else {
		// We can't forward it
		// Delete the job and give up
		msg.msgType = ReplyNotFoundMsgType
		n.send(msg, job.from)
		n.deleteJob(msg)
	}
}

func (n *node) serveReplyNotFound(msg nodeMsg) {
	// We received a file not found
	if n.hasJob(msg) {
		job := n.getJob(msg)
		ksk := msg.body

		// Try again if possible
		dst := n.getRoutingTableEntry(ksk, job.routeNum)
		if dst != nil {
			msg.msgType = RequestDataMsgType
			n.send(msg, dst)
		} else {
			// If we ran out of tries, forward the file not found
			// Only forward if we did not start off this whole request
			if job.from != n {
				n.send(msg, job.from)
			}
			n.deleteJob(msg)
		}
	}
}

// We received a file we wanted
func (n *node) serveSendData(msg nodeMsg) {
	if n.hasJob(msg) {
		job := n.getJob(msg)
		// Only forward if we did not start off this whole request
		if job.from != n {
			n.send(msg, job.from)
		}
		// Cache this file
		n.addFileFromMsg(msg)
		n.deleteJob(msg)
	}
}

func parseFileFromMsg(msg nodeMsg) (string, string) {
	words := strings.Split(msg.body, " ")
	key := words[0]
	file := strings.Join(words[1:], " ")
	return key, file
}

// Add file, given a msg
// Also adds to the routing table
func (n *node) addFileFromMsg(msg nodeMsg) {
	key, file := parseFileFromMsg(msg)
	n.addFile(key, file)
	if msg.origin == nil {
		panic("Adding nil origin to routing table")
	}
	if msg.origin != n {
		n.addRoutingTableEntry(key, msg.origin)
	}
}

// Add file, given the key (KSK)
func (n *node) addFile(fileKey string, file string) {
	n.disk.Add(fileKey, file)
}

// Return the file
func (n *node) getFile(fileKey string) string {
	file, found := n.disk.Get(fileKey)
	if !found {
		panic("Getting file that does not exist")
	}
	return file.(string)

}

// Check if file exists in disk
func (n *node) hasFile(fileKey string) bool {
	_, found := n.disk.Peek(fileKey)
	return found

}

// Add file based on raw descriptor, more useful for tests
func (n *node) addFileDescr(descr string, file string) {
	_, _, fileKey := genKeywordSignedKey(descr)
	n.disk.Add(fileKey, file)
}
