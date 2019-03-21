package card_id

import (
	"strconv"
	"regexp"
	"strings"
	"time"
)

var weightFactorList = [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2} //加权因子
var endVeriCodeMap = map[int]string{//末尾验证码
	0: "1",
	1: "0",
	2: "X",
	3: "9",
	4: "8",
	5: "7",
	6: "6",
	7: "5",
	8: "4",
	9: "3",
	10: "2",
}

const regexpReguler = `^[1-9]\d{5}(19|20)\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$` //身份证号正则

func Validate(s string) bool {
	if len(s) != 18 {
		return false
	}

	m, _ := regexp.MatchString(regexpReguler, s)
	if !m {
		return false
	}

	var sum = 0
	for i, k := range s {
		if i >= 17 {
			break
		}

		weightFactor := weightFactorList[i]
		iCode, _ := strconv.Atoi(string(k))
		v := weightFactor * iCode
		sum += v
	}

	endVeriCode := endVeriCodeMap[sum%11]   //最后一位验证码
	idLast := strings.ToUpper(s[len(s)-1:]) //实际最后一位的值

	if endVeriCode == idLast { //校验通过
		return true
	}

	return false
}

//判断是否成年
func IsAdult(id string) bool {
	isCard := Validate(id)
	if !isCard {
		return false
	}

	var birthDate = id[6:14] //出生日期

	t, _ := time.Parse("20060102", birthDate) //将出生日期转换为时间对象

	t18 := t.AddDate(18, 0, 0) //成年的时间

	if t18.Before(time.Now()) { //已满18岁
		return true
	}

	return false
}
