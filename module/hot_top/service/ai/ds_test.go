package ai

import "testing"

func Test_DeepSeek_Call(t *testing.T) {
	ds, err := NewAiDeepSeek("sk-")
	if err != nil {
		t.Fatal(err)
	}
	answer, err := ds.Call("你好")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(answer)
}
