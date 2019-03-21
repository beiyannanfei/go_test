package coupon

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base32"
	"encoding/binary"
	"github.com/sigurn/crc8"
	"github.com/sirupsen/logrus"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

//noinspection SpellCheckingInspection
var badWordsList = [...]string{
	"SHPX", "PHAG", "JNAX", "JNAT", "CVFF", "PBPX", "FUVG", "GJNG",
	"GVGF", "SNEG", "URYY", "ZHSS", "QVPX", "XABO", "NEFR", "FUNT",
	"GBFF", "FYHG", "GHEQ", "FYNT", "PENC", "CBBC", "OHGG", "SRPX",
	"OBBO", "WVFZ", "WVMM", "CUNG",
}

var symbolsArr = "0123456789ABCDEFGHJKLMNPQRTUVWXY"

var symbolsObj = map[int32]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7,
	'8': 8, '9': 9, 'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15,
	'G': 16, 'H': 17, 'J': 18, 'K': 19, 'L': 20, 'M': 21, 'N': 22, 'P': 23,
	'Q': 24, 'R': 25, 'T': 26, 'U': 27, 'V': 28, 'W': 29, 'X': 30, 'Y': 31,
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Generate() string {

	var parts []string
	var part string

	var optsParts = 3
	var optsPartLen = 4

	var loopTimes = 0

	// default to a random code
	for {
		parts = []string{}

		for i := 0; i < optsParts; i++ {
			part = ""

			for j := 0; j < optsPartLen-1; j++ {
				part += randomSymbol()
			}

			part = part + string(checkDigitAlg1(part, i+1))
			parts = append(parts, part)
		}

		if !hasBadWord(strings.Join(parts, "")) {
			break
		}

		loopTimes++

		if loopTimes > 100 {
			logrus.WithFields(logrus.Fields{
				"optsParts":   3,
				"optsPartLen": 4,
				"part":        part,
				"parts":       parts,
				"loopTimes":   loopTimes,
			}).Info("Coupon.Generate Loop too many times.")
			return ""
		}
	}

	return strings.Join(parts, "-")
}

func EncryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockEncrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}

func DecryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockDecrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
	return nil
}

//const encodeReadable = "ABCDEFGHJKLMNPQRTUVWXY0123456789"
//const oneByteBitSize = 8
//const base32EncodeBitSize = 5
//const indexByteSize = 4

/*
optsParts
encode |--------|--------|--------|--------|--------|--------|--------|...
       |  CRC8  |<--             index           -->|<- random data  ->
 */
func GenerateEx(key string, index uint32, partLens ...int) string {

	// available char length
	var availableCharLength = 0
	for _, i := range partLens {
		availableCharLength += i - 1
	}

	var availableBitLength = availableCharLength * base32EncodeBitSize

	if availableBitLength < int(unsafe.Sizeof(index))*oneByteBitSize {
		return ""
	}

	//var msgLength = 0
	//for _, i := range partLens {
	//	msgLength += i - 1
	//}
	//
	//msgLength *= 5
	//
	//if msgLength < int(unsafe.Sizeof(index))*8 {
	//	return ""
	//}

	offsetBits := oneByteBitSize - availableBitLength%oneByteBitSize

	msgLength := (availableBitLength - offsetBits) / oneByteBitSize

	offsetByte := 0
	if offsetBits > 0 {
		msgLength += 1
		offsetByte = 1
	}

	msg := make([]byte, msgLength)

	rand.Read(msg[indexByteSize:msgLength])
	binary.LittleEndian.PutUint32(msg[1:], index)

	//fmt.Printf("Encrypting %b\n", msg[msgLength-1])
	if offsetBits > 0 {
		msg[msgLength-1] <<= uint(offsetBits)
	}
	//fmt.Printf("Encrypting %b\n", msg[msgLength-1])

	table := crc8.MakeTable(crc8.CRC8_MAXIM)
	crc := crc8.Checksum(msg[1:], table)
	msg[0] = crc
	//fmt.Printf("Encrypting %v, %b\n", crc, msg[msgLength-1])

	//const key16 = "1234567890123456"
	//const key24 = "123456789012345678901234"
	//const key32 = "12345678901234567890123456789012"
	//var key = key16
	//var msg = "message"
	var iv = []byte(key)[:aes.BlockSize] // Using IV same as key is probably bad
	var err error

	// Encrypt
	encrypted := make([]byte, msgLength)
	err = EncryptAESCFB(encrypted, []byte(msg[0:msgLength-offsetByte]), []byte(key), iv)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Encrypting %v %s -> %v\n", []byte(msg[0:msgLength]), msg[0:msgLength], encrypted)

	if offsetBits > 0 {
		encrypted[msgLength-1] = msg[msgLength-1]
	}

	encoder := base32.NewEncoding(encodeReadable)
	encoder = encoder.WithPadding(base32.NoPadding)
	str := encoder.EncodeToString(encrypted)

	//fmt.Printf("Encrypting %v %s -> %v\n", []byte(msg), msg, encrypted)

	//fmt.Printf("Encrypting %v %s -> %v %s\n", []byte(msg), msg, encrypted, str)

	//decStr, _ := encoder.DecodeString(str)
	//decByte := []byte(decStr)
	//
	//fmt.Printf("Encrypting %v\n", decByte)
	//
	//const key16 = "K234567890123456"
	//
	//// Decrypt
	//decrypted := make([]byte, len(msg))
	//err = DecryptAESCFB(decrypted, encrypted, []byte(key), iv)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Decrypting %v %v-> %v %s\n", encrypted, decStr, decrypted, decrypted)

	var parts []string
	s := 0
	for i, v := range partLens {

		var part = str[s : s+v-1]
		s += v - 1
		part += string(checkDigitAlg1(part, i+1))

		parts = append(parts, part)
	}

	return strings.Join(parts, "-")
	//return str[0:availableCharLength]
}

func randomSymbol() string {
	return string(symbolsArr[rand.Intn(len(symbolsArr))])
}

// returns the checksum character for this (data/part) combination
func checkDigitAlg1(data string, check int) uint8 {
	// check's initial value is the part number (e.g. 3 or above)

	// loop through the data chars
	for _, v := range data {
		var k = symbolsObj[v]
		check = check*19 + k
	}

	return symbolsArr[ check%31 ]
}

func hasBadWord(code string) bool {
	code = strings.ToUpper(code)
	for i := 0; i < len(badWordsList); i++ {

		if strings.Index(code, badWordsList[i]) > -1 {
			return true
		}
	}
	return false
}

func Validate(code string) (bool, string) {
	if len(code) == 0 {
		return false, ""
	}

	// uppercase the code, take out any random chars and replace OIZS with 0125
	code = strings.ToUpper(code)
	code = strings.Replace(code, "O", "0", -1)
	code = strings.Replace(code, "I", "1", -1)
	code = strings.Replace(code, "Z", "2", -1)
	code = strings.Replace(code, "S", "5", -1)

	re := regexp.MustCompile("[^0-9A-Z]+")
	code = re.ReplaceAllString(code, "")

	var optsParts = 3
	var optsPartLen = 4

	// split in the different parts
	var parts []string
	var tmp = code
	for ; len(tmp) > 0; {

		parts = append(parts, tmp[0:optsPartLen])
		tmp = tmp[optsPartLen:]
	}

	// make sure we have been given the same number of parts as we are expecting
	if len(parts) != optsParts {
		return false, ""
	}

	// validate each part
	var part, check, data string
	for i := 0; i < len(parts); i++ {
		part = parts[i]
		// check this part has 4 chars
		if len(part) != optsPartLen {
			return false, ""
		}

		// split out the data and the check
		data = part[0 : optsPartLen-1]
		check = part[optsPartLen-1 : optsPartLen]

		if check != string(checkDigitAlg1(data, i+1)) {
			return false, ""
		}
	}

	// everything looked ok with this code
	//return parts.join('-')
	return true, strings.Join(parts, "-")
}
