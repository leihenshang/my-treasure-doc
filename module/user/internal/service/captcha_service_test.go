package service

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestGenCaptchaResponseUsesCamelCase(t *testing.T) {
	resp, err := GenCaptcha(nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.CaptchaId == "" || resp.Captcha == "" {
		t.Fatalf("captcha is empty: %#v", resp)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"captchaId"`, `"captcha"`} {
		if !strings.Contains(string(data), key) {
			t.Fatalf("captcha response missing %s: %s", key, data)
		}
	}
}

func TestVerifyCaptchaWrongAttemptConsumesCaptcha(t *testing.T) {
	resp, err := GenCaptcha(nil)
	if err != nil {
		t.Fatal(err)
	}
	answer := store.Get(resp.CaptchaId, false)
	// 失败尝试同样消耗验证码，避免同一张图被反复猜测。
	if err := verifyCaptcha(resp.CaptchaId, "wrong", true); !errors.Is(err, ErrCaptchaInvalid) {
		t.Fatalf("wrong captcha error = %v, want ErrCaptchaInvalid", err)
	}
	if err := verifyCaptcha(resp.CaptchaId, answer, true); !errors.Is(err, ErrCaptchaInvalid) {
		t.Fatalf("captcha should be consumed by the failed attempt, got %v", err)
	}
}

func TestVerifyCaptchaRejectsEmptyAndWrongValues(t *testing.T) {
	resp, err := GenCaptcha(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ name, id, value string }{
		{name: "empty", id: "", value: ""},
		{name: "missing id", id: "", value: "42"},
		{name: "missing value", id: resp.CaptchaId, value: ""},
		{name: "wrong value", id: resp.CaptchaId, value: "not-the-answer"},
		{name: "unknown id", id: "unknown-captcha-id", value: "42"},
	} {
		if err := verifyCaptcha(test.id, test.value, true); !errors.Is(err, ErrCaptchaInvalid) {
			t.Fatalf("%s: error = %v, want ErrCaptchaInvalid", test.name, err)
		}
	}
}

func TestVerifyCaptchaIsSingleUse(t *testing.T) {
	resp, err := GenCaptcha(nil)
	if err != nil {
		t.Fatal(err)
	}
	answer := store.Get(resp.CaptchaId, false)
	if answer == "" {
		t.Fatal("captcha answer should be stored")
	}
	if err := verifyCaptcha(resp.CaptchaId, answer, true); err != nil {
		t.Fatalf("valid captcha rejected: %v", err)
	}
	// 同一个验证码不能重复用于第二次登录尝试。
	if err := verifyCaptcha(resp.CaptchaId, answer, true); !errors.Is(err, ErrCaptchaInvalid) {
		t.Fatalf("reused captcha error = %v, want ErrCaptchaInvalid", err)
	}
}
