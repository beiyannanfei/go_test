package quick_sdk

import (
	"testing"
	"fmt"
)

func TestYsdkVerifyLogin(t *testing.T) {
	test := struct {
		name  string
		token string
		uid   string
	}{
		"TestVerifyLoginQuick",
		"@178@83@173@158@157@88@108@86@118@98@117@107@105@106@108@99@104@120@108@103@112@123@125@106@96@101@104@110@104@115@105@101@169@187@175@156@163@183@152@164@101@155@134@217@160@158@157@87@115@103@102@105@101@99@105@99@99@110@105@92@85@154@157@152@165@163@158@163@121@151@90@112@90@100@102@87@157@151@219@196@217@215@134@165@121@163@225",
		"D2A864635A709FD302080B508FF98D49",
	}

	t.Run(test.name, func(t *testing.T) {
		resp, err := VerifyLoginQuick(test.token, test.uid)
		fmt.Println(resp, err)
	})
}
