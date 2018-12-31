package main

import (
	"strings"
	"strconv"
	"fmt"
	"encoding/xml"
	"reflect"
)

func main() {
	str := "@116@119@168@161@165@88@170@153@167@170@161@163@166@113@87@99@94@102@83@85@153@166@154@164@155@157@166@152@114@90@140@135@126@101@104@86@89@171@168@149@163@155@153@160@167@162@154@111@82@164@160@87@115@118@115@168@162@173@165@160@164@166@170@146@165@157@163@167@154@159@153@114@113@164@157@167@171@149@156@151@110@114@154@168@147@172@156@168@171@114@104@109@100@161@170@146@172@157@163@168@119@116@151@156@150@165@166@153@164@114@109@106@104@110@109@100@151@160@152@163@165@153@164@111@113@155@159@148@166@166@149@160@152@173@157@152@115@165@163@147@155@102@174@170@154@108@149@160@106@161@164@161@158@162@177@153@171@158@115@98@155@160@145@162@167@157@160@147@170@160@156@114@116@155@150@159@149@149@160@167@152@157@169@115@153@107@158@151@104@153@108@149@153@155@97@102@109@107@153@149@109@112@113@101@154@150@152@106@149@152@102@153@149@153@152@155@115@99@159@146@162@157@150@162@170@156@149@166@119@116@163@166@153@156@170@147@166@163@115@100@100@103@99@101@101@110@103@108@105@102@105@97@106@106@108@99@113@105@99@101@110@111@109@109@104@115@103@163@170@152@154@164@143@164@160@115@112@168@152@174@150@168@161@158@154@118@105@99@105@110@93@100@112@101@102@102@85@104@104@110@109@102@111@103@96@114@96@165@149@177@150@169@160@161@157@111@113@153@164@162@173@166@164@114@106@102@100@100@113@102@153@161@167@169@163@166@110@114@164@169@149@172@172@168@117@100@116@96@168@172@152@167@173@171@110@112@158@176@168@166@150@170@151@164@153@166@150@159@163@116@172@102@177@151@178@103@180@112@103@150@173@172@169@148@171@151@160@149@171@153@161@167@115@115@103@161@157@167@168@147@151@155@111@113@99@171@162@174@164@163@167@159@168@151@164@152@171@171@145@155@158@118"
	key := "88049844578484520615487574815873"
	xmlStr := decode(str, key)
	parseXML(xmlStr)
}

func parseXML(xmlStr string) {
	fmt.Println(xmlStr)
	type SkymoonsMessage struct {
		Message struct {
			IsTest       string `xml:"is_test"`
			Channel      string `xml:"channel"`
			ChannelUid   string `xml:"channel_uid"`
			GameOrder    string `xml:"game_order"`
			OrderNo      string `xml:"order_no"`
			PayTime      string `xml:"pay_time"`
			Amount       string `xml:"amount"`
			Status       string `xml:"status"`
			ExtrasParams string `xml:"extras_params"`
		} `xml:"message"`
	}

	v := SkymoonsMessage{}
	err := xml.Unmarshal([]byte(xmlStr), &v)
	fmt.Println(err)
	fmt.Println(v)
	fmt.Println(v.Message.Amount)
}

func decode(str string, key string) string {
	list := strings.Split(str, "@")[1:]

	keysByte := []byte(key)
	var dataByte []byte

	/*for i := 0; i < len(list); i++ {
		l, _ := strconv.Atoi(list[i])
		lu := byte(l)
		k := keysByte[i%len(keysByte)]
		cha := lu - k
		//fmt.Println(l, reflect.TypeOf(l), lu, reflect.TypeOf(lu), k, reflect.TypeOf(k), cha, reflect.TypeOf(cha))
		dataByte = append(dataByte, cha)
	}*/

	for i := 0; i < len(list); i++ {
		l, _ := strconv.Atoi(list[i])
		k := keysByte[i%len(keysByte)]
		cha := l - int(k)
		fmt.Println(l, reflect.TypeOf(l), k, reflect.TypeOf(k), cha, reflect.TypeOf(cha))
		dataByte = append(dataByte, byte(cha))
	}

	//fmt.Println(string(dataByte))
	return string(dataByte)
}

//package quick
//
//import (
//	"strconv"
//	"strings"
//)
//
///**
//QuickSDK游戏同步加解密算法描述
//解密方法
//strEncode 密文
//keys 解密密钥 为游戏接入时分配的 callback_key
//*/
//func Decode(str string, keys string) string {
//
//	strs := strings.Split(str, "@")
//	strs = strs[1:]
//
//	keysNum := GetBytes(keys)
//
//	_data := []int{}
//	_len := len(keysNum)
//
//	for i, v := range strs {
//
//		keyVar := keysNum[i%_len]
//		kn, _ := strconv.Atoi(v)
//		_data = append(_data, kn-0xff&keyVar)
//	}
//
//	return ToStr(_data)
//
//}
//
///**
//* 转成字符数据
//*/
//func GetBytes(strs string) []int {
//
//	_keys := []byte(strs)
//	keysNum := []int{}
//	for _, _n := range _keys {
//
//		num := int(_n)
//		keysNum = append(keysNum, num)
//	}
//
//	return keysNum
//}
//
///**
//* 转化字符串
//*/
//func ToStr(keysNum []int) string {
//
//	_b := []string{}
//	for _, v := range keysNum {
//
//		_b = append(_b, string(v))
//	}
//
//	return strings.Join(_b, "")
//}
