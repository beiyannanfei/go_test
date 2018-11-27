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
		Mode:               base64Captcha.CaptchaModeArithmetic,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         4,
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
}
