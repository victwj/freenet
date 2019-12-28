package main

import (
	"testing"
	"time"
)

func TestKSK(t *testing.T) {
	_, _, a := genKeywordSignedKey("/test/test/hello")
	_, _, b := genKeywordSignedKey("/test/test/hello")
	_, _, c := genKeywordSignedKey("/test/test/hello")

	if !(a == b && a == c && b == c) {
		t.Error("KSK not deterministic")
	}
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

	n1.addJob(testMsg2)
	time.Sleep(6 * time.Second)
	b = n1.getJob(testMsg2)
	if b != nil {
		t.Error("Expired job still found")
	}
}

func TestStringSimilarity(t *testing.T) {
	s1 := "abc"
	s2 := "def"
	s3 := "aef"
	s4 := "abdc"

	if stringSimilarity(s1, s2) != 0 {
		t.Error("Error with comparing", s1, s2)
	}
	if stringSimilarity(s1, s3) != 1 {
		t.Error("Error with comparing", s1, s3)
	}
	if stringSimilarity(s2, s3) != 2 {
		t.Error("Error with comparing", s2, s3)
	}
	if stringSimilarity(s1, s1) != 3 {
		t.Error("Error with comparing", s1, s1)
	}
	if stringSimilarity(s1, s4) != 2 {
		t.Error("Error with comparing", s1, s4)
	}
}
