package gid

import "testing"

func Test_genSid(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(genSid())
	}
}
