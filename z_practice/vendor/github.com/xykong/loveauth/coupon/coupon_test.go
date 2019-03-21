package coupon

import (
	"math/rand"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", "TJRU-AU0G-F58Q"},
	}

	rand.Seed(0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(); got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	rand.Seed(0)
	for i := 0; i < b.N; i++ {
		Generate()
	}
}

func BenchmarkValidate(b *testing.B) {
	rand.Seed(0)
	for i := 0; i < b.N; i++ {
		var code = "TJRU-AU0G-F58Q"
		Validate(code)
	}
}

func BenchmarkGenerateAndValidate(b *testing.B) {
	rand.Seed(0)
	for i := 0; i < b.N; i++ {
		var code = Generate()
		Validate(code)
	}
}

func Test_randomSymbol(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", "T"},
	}

	rand.Seed(0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomSymbol(); got != tt.want {
				t.Errorf("randomSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDigitAlg1(t *testing.T) {
	type args struct {
		data  string
		check int
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"", args{data: "", check: 0}, 48},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDigitAlg1(tt.args.data, tt.args.check); got != tt.want {
				t.Errorf("checkDigitAlg1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasBadWord(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{code: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasBadWord(tt.args.code); got != tt.want {
				t.Errorf("hasBadWord(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}

	// check each rude word not locatable in the string
	for i := 0; i < len(badWordsList); i++ {

		code := badWordsList[i] + "YYYYZZZZ"
		if got := hasBadWord(code); got != true {
			t.Errorf("hasBadWord(%v) = %v, want %v", code, got, true)
		}

		code = "XXXX" + badWordsList[i] + "ZZZZ"
		if got := hasBadWord(code); got != true {
			t.Errorf("hasBadWord(%v) = %v, want %v", code, got, true)
		}

		code = "YYYYZZZZ" + badWordsList[i]
		if got := hasBadWord(code); got != true {
			t.Errorf("hasBadWord(%v) = %v, want %v", code, got, true)
		}
	}

	code := "XXXX" + strings.ToLower(badWordsList[0]) + "YYYYZZZZ"
	if got := hasBadWord(code); got != true {
		t.Errorf("hasBadWord(%v) = %v, want %v", code, got, true)
	}

	code = "XXXXYYYYZZZZ"
	if got := hasBadWord("XXXXYYYYZZZZ"); got != false {
		t.Errorf("hasBadWord(%v) = %v, want %v", code, got, false)
	}
}

func TestValidate(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{"", args{"1K7Q-CTFM-LMTC"}, true, "1K7Q-CTFM-LMTC"},
		{"", args{"1k7q-ctfm-lmtc"}, true, "1K7Q-CTFM-LMTC"},
		{"", args{"I9oD-V467-8D52"}, true, "190D-V467-8D52"},
		{"", args{"I9OD-V467-8D52"}, true, "190D-V467-8D52"},
		{"", args{"i9OD-V467-8D52"}, true, "190D-V467-8D52"},
		{"", args{"I9OD-V467-8D52"}, true, "190D-V467-8D52"},
		{"", args{"I9OD-V467-8D5z"}, true, "190D-V467-8D52"},
		{"", args{"I9OD-V467-8D5Z"}, true, "190D-V467-8D52"},
		{"", args{"I9OD-V467-8DsZ"}, true, "190D-V467-8D52"},
		{"", args{"I9OD-V467-8DSZ"}, true, "190D-V467-8D52"},
		{"", args{" i9oD V467 8Dsz "}, true, "190D-V467-8D52"},
		{"", args{" i9oD_V467_8Dsz "}, true, "190D-V467-8D52"},
		{"", args{"i9oDV4678Dsz"}, true, "190D-V467-8D52"},
		{"", args{"1K7Q-CTFM-LMTC"}, true, "1K7Q-CTFM-LMTC"},
		{"", args{"1K7Q-CTFM-LMT1"}, false, ""},
		{"", args{"1K7Q-CTFM"}, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Validate(tt.args.code)
			if got != tt.want {
				t.Errorf("Validate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Validate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestEncryptAESCFB(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
		key []byte
		iv  []byte
	}

	var dst = make([]byte, 16)
	var src = []byte("message")
	const key = "1234567890123456"

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{dst, src, []byte(key), []byte(key)}, false},
		{"", args{dst, src, []byte(key), []byte(key)}, false},
		{"", args{dst, src, []byte(key), []byte(key)}, false},
		{"", args{dst, src, []byte(key), []byte(key)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EncryptAESCFB(tt.args.dst, tt.args.src, tt.args.key, tt.args.iv); (err != nil) != tt.wantErr {
				t.Errorf("EncryptAESCFB(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
		})
	}
}

func TestDecryptAESCFB(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
		key []byte
		iv  []byte
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecryptAESCFB(tt.args.dst, tt.args.src, tt.args.key, tt.args.iv); (err != nil) != tt.wantErr {
				t.Errorf("DecryptAESCFB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateEx(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", "RK89-64D8-G6B4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateEx("1234567890123456", 1, 4, 4, 4); got != tt.want {
				t.Errorf("GenerateEx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGenerateEx(b *testing.B) {
	rand.Seed(0)
	for i := 0; i < b.N; i++ {
		GenerateEx("1234567890123456", 1, 4, 4, 4)
	}
}
