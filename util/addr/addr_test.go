package addr

import "testing"

func TestExtract(t *testing.T) {
	ip := "10.1.160.78"
	extractId, _ := Extract("")
	if ip != extractId {
		t.Fatal("extract id fail")
	}
}
