package service

import (
	"errors"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore
var CaptchaFonts []string

const (
	CaptchaHeight          = 40
	CaptchaWidth           = 240
	CaptchaShowLineOptions = 2
	CaptchaNoiseCount      = 0
	CaptchaBgColorR        = 221
	CaptchaBgColorG        = 221
	CaptchaBgColorB        = 3
	CaptchaBgColorA        = 221
)

type CaptchaResp struct {
	CaptchaId string
	Captcha   string
}

func GenCaptcha(ctx *gin.Context) (CaptchaResp, error) {
	resp := CaptchaResp{}
	driver := (&MyDriverMath{
		DriverMath: &base64Captcha.DriverMath{
			Height:          CaptchaHeight,
			Width:           CaptchaWidth,
			ShowLineOptions: CaptchaShowLineOptions,
			NoiseCount:      CaptchaNoiseCount,
			Fonts:           CaptchaFonts,
			BgColor:         &color.RGBA{R: CaptchaBgColorR, B: CaptchaBgColorB, A: CaptchaBgColorA, G: CaptchaBgColorG},
		},
	}).ConvertFonts()

	id, b64s, _, err := base64Captcha.NewCaptcha(driver, store).Generate()
	if err != nil {
		return resp, err
	}
	resp.CaptchaId = id
	resp.Captcha = b64s
	return resp, nil
}

func verifyCaptcha(id, verifyVal string, clear bool) error {
	if !store.Verify(id, verifyVal, clear) {
		return errors.New("captcha verification failed")
	}
	return nil
}

type MyDriverMath struct {
	*base64Captcha.DriverMath
}

func (d *MyDriverMath) GenerateIdQuestionAnswer() (id, question, answer string) {
	id = base64Captcha.RandomId()
	operators := []string{"+", "-"}
	var mathResult int32
	switch operators[rand.Int31n(2)] {
	case "+":
		a := rand.Int31n(10)
		b := rand.Int31n(10)
		question = fmt.Sprintf("%d+%d=?", a, b)
		mathResult = a + b
	default:
		a := rand.Int31n(8) + rand.Int31n(2)
		b := rand.Int31n(10)
		if a < b {
			a, b = b, a
		}

		question = fmt.Sprintf("%d-%d=?", a, b)
		mathResult = a - b

	}
	answer = fmt.Sprintf("%d", mathResult)
	return
}

// ConvertFonts loads fonts from names
func (d *MyDriverMath) ConvertFonts() *MyDriverMath {
	d.DriverMath = d.DriverMath.ConvertFonts()
	return d
}
