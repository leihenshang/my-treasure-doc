package token

import "testing"

func Test_GetBilibili(t *testing.T) {
	wbi, err := GetBilibiliWbi()
	if err != nil {
		t.Errorf("GetBilibiliWbi failed, err: %v", err)
	}
	t.Logf("wbi: %s", wbi)
}
