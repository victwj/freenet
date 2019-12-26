// Functions related to inserting/removing files
package main

func (n *node) sendFileInsert(descr string, file string) {

}

func (n *node) sendFileRequest(descr string) {

}

func (n *node) serveFileRequest(descr string) {

}

func (n *node) addFile(fileKey string, file string) {
	n.disk.Add(fileKey, file)
}
