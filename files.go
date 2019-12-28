// Functions related to inserting/removing files
package main

func (n *node) sendFileInsert(descr string, file string) {

}

// Sending a request
func (n *node) sendRequestData(descr string, dst *node) {

	msg := n.newNodeMsg(RequestDataMsgType, descr)

	// Add the job, if success, send it
	if n.addJob(msg) {
		n.send(msg, dst)
	}
}

func (n *node) serveRequestData(msg nodeMsg) {
	// For now, send a failure right back
	msg.msgType = FailMsgType
	n.send(msg, msg.from)

}

func (n *node) addFile(fileKey string, file string) {
	n.disk.Add(fileKey, file)
}
