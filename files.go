// Functions related to inserting/removing files
package main

func (n *node) sendFileInsert(descr string, file string) {

}

// Sending a request
func (n *node) sendRequestData(descr string) {

	msg := n.newNodeMsg(RequestDataMsgType, descr)

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

// TODO: function WIP
func (n *node) serveRequestData(msg nodeMsg) {
	// The descr string is in the body
	_, _, ksk := genKeywordSignedKey(msg.body)
	fileFound := n.hasFile(ksk)

	// If file is found, return it //TODO:
	if fileFound {
		// Don't forget to delete job too
		return
	}

	// File is not found
	// Create a job for this request if it doesn't exist
	if !n.hasJob(msg) {
		// Create the job, but it might fail
		// If we can't create the job, send a not found
		if !n.addJob(msg) {
			msg.msgType = ReplyNotFoundMsgType
			n.send(msg, msg.from)
			return
		}
	}

	// Forward the request for the file since we don't have it
	job := n.getJob(msg)
	dst := n.getRoutingTableEntry(ksk, job.routeNum)
	// If there is a node to send to, send it
	if dst != nil {
		n.send(msg, dst)
	} else {
		// We ran out of possible forwarding nodes, send not found
		// Delete the job and give up
		msg.msgType = ReplyNotFoundMsgType
		n.send(msg, job.from)
		n.deleteJob(msg)
	}
}

func (n *node) serveReplyNotFound(msg nodeMsg) {
	// We received a file not found
	// Delete the job associated with this request, and forward it
	if n.hasJob(msg) {
		job := n.getJob(msg)
		// Only forward if we did not start off this whole request
		if job.from != n {
			n.send(msg, job.from)
		}
		n.deleteJob(msg)
	}
}

func (n *node) addFile(fileKey string, file string) {
	n.disk.Add(fileKey, file)
}

func (n *node) getFile(fileKey string) string {
	file, found := n.disk.Get(fileKey)
	if !found {
		return ""
	}
	return file.(string)

}

func (n *node) hasFile(fileKey string) bool {
	_, found := n.disk.Peek(fileKey)
	return found

}
