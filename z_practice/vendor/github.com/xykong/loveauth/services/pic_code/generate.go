package pic_code

//生成图片验证
import (
	"time"
	"math/rand"
	"github.com/beiyannanfei/base64Captcha"
	"encoding/base64"
)

func Generate() (string, string) {
	rand.Seed(time.Now().UnixNano())
	var configC = base64Captcha.ConfigCharacter{
		Height: 60,
		Width:  240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:            base64Captcha.CaptchaModeArithmetic,
		IsUseSimpleFont: true,
		//ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		//ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		//IsShowHollowLine:   false,
		//IsShowNoiseDot:     false,
		//IsShowNoiseText:    false,
		//IsShowSlimeLine:    false,
		//IsShowSineLine:     false,
		//CaptchaLen:         4,
	}
	reults, digitCap := base64Captcha.GenerateCaptcha(configC) //注意: 已经对源码做修改
	binaryData := digitCap.BinaryEncodeing()
	base64Png := base64.StdEncoding.EncodeToString(binaryData)
	return reults, base64Png
}
