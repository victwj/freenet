package freenet

import (
	"testing"
)

func TestFreenetBasic(t *testing.T) {

}

func TestJob(t *testing.T) {
	n1 := newNode(5)
	testMsg1 := n1.newNodeMsg(FailMsgType, "test msg 1")
	testMsg2 := n1.newNodeMsg(FailMsgType, "test msg 2")

	n1.addJob(testMsg1)
	a := n1.getJob(testMsg1)
	if a == nil {
		t.Error("Job not found")
	}

	b := n1.getJob(testMsg2)
	if b != nil {
		t.Error("Nonexistent job found")
	}

	n1.deleteJob(testMsg1)
	a = n1.getJob(testMsg1)
	if a != nil {
		t.Error("Deleted job still found")
	}

	// n1.addJob(testMsg2)
	// time.Sleep(6 * time.Second)
	// b = n1.getJob(testMsg2)
	// if b != nil {
	// 	t.Error("Expired job still found")
	// }
}
