// Functions related to inserting/removing files
package main

func (n *node) sendFileInsert(descr string, file string) {

}

// TODO: make routing table functions
func (n *node) sendRequestData(descr string, dst *node) {
	msg := n.newNodeMsg(RequestDataMsgType, descr)
	n.send(msg, dst)
}

func (n *node) serveRequestData(msg nodeMsg) {

}

func (n *node) addFile(fileKey string, file string) {
	n.disk.Add(fileKey, file)
}
