package order

import "testing"

func TestGenerateOrderNo(t *testing.T) {
	count := 100
	collections := func() (res []string) {
		for i := 0; i < count; i++ {
			no, err := GenerateOrderNo()
			if err != nil {
				t.Error(err)
			}
			res = append(res, no)
		}

		return
	}()

	noMap := make(map[string]interface{}, 0)
	for _, v := range collections {
		noMap[v] = new(interface{})
		if v == "" {
			t.Error("no err,is empty")
		}

		t.Logf("no:%+v,len:%+v", v, len(v))
	}

	if len(noMap) != count {
		t.Error("order no is empty.")
	}

	t.Log("success")
}
