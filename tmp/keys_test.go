package freenet

import "testing"

func TestKSK(t *testing.T) {
	_, _, a := genKeywordSignedKey("/test/test/hello")
	_, _, b := genKeywordSignedKey("/test/test/hello")
	_, _, c := genKeywordSignedKey("/test/test/hello")

	if !(a == b && a == c && b == c) {
		t.Error("KSK not deterministic")
	}
}
