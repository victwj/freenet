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

func (n *node) serveRequestData(msg nodeMsg) {
	// TODO: For now, send a failure right back
	msg.msgType = FailMsgType
	n.send(msg, msg.from)
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
