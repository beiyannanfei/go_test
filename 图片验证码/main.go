package main

import (
	"github.com/mojocn/base64Captcha"
	"fmt"
	"os"
	"time"
	"math/rand"
)

//ConfigJsonBody json request body.
type ConfigJsonBody struct {
	Id              string
	CaptchaType     string
	VerifyValue     string
	ConfigAudio     base64Captcha.ConfigAudio
	ConfigCharacter base64Captcha.ConfigCharacter
	ConfigDigit     base64Captcha.ConfigDigit
}

//start a net/http server
//启动golang net/http 服务器
func main() {
	/*
		//serve Vuejs+ElementUI+Axios Web Application
		http.Handle("/", http.FileServer(http.Dir("./static")))

		//api for create captcha
		//创建图像验证码api
		http.HandleFunc("/api/getCaptcha", generateCaptchaHandler)

		//api for verify captcha
		http.HandleFunc("/api/verifyCaptcha", captchaVerifyHandle)

		fmt.Println("Server is at localhost:3333")
		if err := http.ListenAndServe("localhost:3333", nil); err != nil {
			log.Fatal(err)
		}*/

	fmt.Println("==========================")
	userFile := "/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/图片验证码/a.png"
	fout, err := os.Create(userFile)
	defer fout.Close()
	if err != nil {
		fmt.Println("create file err:", err)
		return
	}

	/*var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 4,
	}

	vId, digitCap := base64Captcha.GenerateCaptcha("", configD)
	digitCap.WriteTo(fout)
	fmt.Println(vId)*/
	rand.Seed(time.Now().UnixNano())
	var configC = base64Captcha.ConfigCharacter{
		Height: 60,
		Width:  240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode: base64Captcha.CaptchaModeArithmetic,
		//ComplexOfNoiseText: base64Captcha.CaptchaComplexHigh,
		//ComplexOfNoiseDot:  base64Captcha.CaptchaComplexHigh,
		//IsShowHollowLine:   false,
		//IsShowNoiseDot:     false,
		//IsShowNoiseText:    false,
		//IsShowSlimeLine:    false,
		//IsShowSineLine:     false,
		//CaptchaLen:         4,
		IsUseSimpleFont: true,
	}
	vId, digitCap := base64Captcha.GenerateCaptcha(configC)
	//digitCap.WriteTo(fout)
	fmt.Println(vId)

	binaryData := digitCap.BinaryEncodeing()
	fout.Write(binaryData)

	//verifyResult := base64Captcha.VerifyCaptcha("vId", "8")
	//fmt.Println(verifyResult)
	//base64Png := base64Captcha.CaptchaWriteToBase64Encoding(digitCap)
	//fmt.Println(base64Png)

	//s := "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"

}
