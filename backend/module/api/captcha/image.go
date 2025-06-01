package captcha

import "github.com/mojocn/base64Captcha"

// 基于内存的store
var store = base64Captcha.DefaultMemStore

func GenerateImageCaptcha() (id string, b64s string, err error) {
	driverString := &base64Captcha.DriverString{
		Height:          60,
		Width:           240,
		NoiseCount:      0,
		ShowLineOptions: 02,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
	}
	driver := driverString // 多态转换

	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err = c.Generate()
	if err != nil {
		return "", "", err
	}
	return
}

func VerifyImageCaptcha(id string, value string) bool {
	return store.Verify(id, value, true)
}
