// Package coupon implements coupon encoding as specified by xykong.
package coupon

import (
	"crypto/aes"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/sigurn/crc8"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

/*
 * Encodings
 */

// An Encoding is a radix 32 encoding/decoding scheme, defined by a
// 32-character alphabet. The most common is the "base32" encoding
// introduced for SASL GSSAPI and standardized in RFC 4648.
// The alternate "base32hex" encoding is used in DNSSEC.
type Encoding struct {
	key                 string
	partLens            []int
	dataByteSize        int
	availableCharLength int
	availableBitLength  int
	offsetBits          int
	msgLength           int
	offsetByte          int
	encodedLength       int
	decodedLength       int
	digitCheckingSize   int
	crc8Table           *crc8.Table
	readableEncoder     *base32.Encoding
}

const encodeReadable = "ABCDEFGHJKLMNPQRTUVWXY0123456789"
const oneByteBitSize = 8
const base32EncodeBitSize = 5
const indexByteSize = 4

// NewEncoding returns a new Encoding defined by the given alphabet,
// which must be a 32-byte string.
func NewEncoding(key string, dataByteSize int, digitCheckingSize int, partLens ...int) *Encoding {

	e := new(Encoding)

	e.key = key
	e.dataByteSize = dataByteSize
	e.digitCheckingSize = digitCheckingSize
	e.partLens = partLens

	// available char length
	for _, i := range partLens {
		e.availableCharLength += i - digitCheckingSize
	}

	e.availableBitLength = e.availableCharLength * base32EncodeBitSize

	if e.availableBitLength < dataByteSize*oneByteBitSize {

		logrus.WithFields(logrus.Fields{
			"key":                         key,
			"dataByteSize":                dataByteSize,
			"partLens":                    partLens,
			"availableBitLength":          e.availableBitLength,
			"dataByteSize*oneByteBitSize": dataByteSize * oneByteBitSize,
		}).Error("Encoder require more date space.")

		return nil
	}

	e.offsetBits = oneByteBitSize - e.availableBitLength%oneByteBitSize
	e.msgLength = (e.availableBitLength - e.offsetBits) / oneByteBitSize

	if e.offsetBits > 0 {
		e.msgLength += 1
		e.offsetByte = 1
	}

	e.encodedLength = e.availableCharLength + len(partLens) + len(partLens)*digitCheckingSize - 1
	e.decodedLength = dataByteSize

	e.readableEncoder = base32.NewEncoding(encodeReadable)
	e.readableEncoder = e.readableEncoder.WithPadding(base32.NoPadding)

	e.crc8Table = crc8.MakeTable(crc8.CRC8_MAXIM)

	return e
}

func (enc *Encoding) reverse(src []byte, start int, end int) []byte {

	for ; start < end; {
		var temp = src[end]
		src[end] = src[start]
		src[start] = temp

		start += 1
		end -= 1
	}

	return src
}

func (enc *Encoding) shift(src []byte, start int, end int, k int) []byte {

	if end <= start {
		return src
	}

	k = k % (end - start)

	enc.reverse(src[start:end], start, start+k-1)
	enc.reverse(src[start:end], start+k, end-start-1)
	enc.reverse(src[start:end], start, end-start-1)

	return src
}

func (enc *Encoding) xor(src []byte, data byte) []byte {

	for i := 0; i < len(src); i++ {
		src[i] ^= data
	}

	return src
}

/*
 * Encoder
 */

// Encode encodes src using the encoding enc, writing
// EncodedLen(len(src)) bytes to dst.
//
// The encoding pads the output to a multiple of 8 bytes,
// so Encode is not appropriate for use on individual blocks
// of a large data stream. Use NewEncoder() instead.
func (enc *Encoding) encode(dst []byte, src uint64) {

	msg := make([]byte, enc.msgLength)

	//binary.LittleEndian.PutUint32(msg[1:], src)

	var dataLength = binary.PutUvarint(msg[1:], uint64(src))

	rand.Read(msg[dataLength+1 : enc.msgLength])

	//fmt.Printf("Encrypting %b\n", msg[enc.msgLength-1])
	if enc.offsetBits > 0 {
		msg[enc.msgLength-1] <<= uint(enc.offsetBits)
	}
	//fmt.Printf("Encrypting %b\n", msg[enc.msgLength-1])

	crc := crc8.Checksum(msg[1:], enc.crc8Table)
	msg[0] = crc
	//fmt.Printf("Encrypting %v, %b\n", crc, msg[enc.msgLength-1])

	//var msg = "message"
	var iv = []byte(enc.key)[:aes.BlockSize] // Using IV same as key is probably bad
	var err error

	// Encrypt
	encrypted := make([]byte, enc.msgLength)
	err = EncryptAESCFB(encrypted, []byte(msg[0:enc.msgLength-enc.offsetByte]), []byte(enc.key), iv)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("Encrypting %v %s -> %v\n", []byte(msg[0:enc.msgLength]), msg[0:enc.msgLength], encrypted)

	if enc.offsetBits > 0 {
		encrypted[enc.msgLength-1] = msg[enc.msgLength-1]
	}

	if enc.digitCheckingSize == 0 {
		//fmt.Printf("encode shift %d %v\n", int(encrypted[enc.msgLength-1]), []byte(encrypted))

		//enc.shift(encrypted, 0, enc.msgLength-1, int(encrypted[enc.msgLength-1]))
		enc.xor(encrypted[0:enc.msgLength-1], encrypted[enc.msgLength-1])

		//fmt.Printf("       shift %d %v\n", int(encrypted[enc.msgLength-1]), []byte(encrypted))
	}

	str := enc.readableEncoder.EncodeToString(encrypted)

	//fmt.Printf("Encrypting %v %s -> %v %s\n", []byte(msg), msg, encrypted, str)

	var parts []string
	s := 0
	for i, v := range enc.partLens {

		var part = str[s : s+v-enc.digitCheckingSize]
		s += v - enc.digitCheckingSize

		//noinspection GoBoolExpressions
		if enc.digitCheckingSize > 0 {
			part += string(checkDigitAlg1(part, i+enc.digitCheckingSize))
		}

		parts = append(parts, part)
	}

	copy(dst[:], strings.Join(parts, "-"))

	return
}

// EncodeToString returns the base32 encoding of src.
func (enc *Encoding) EncodeToString(src uint64) string {

	buf := make([]byte, enc.EncodedLen())
	enc.encode(buf, src)
	return string(buf)
}

// EncodeToString returns the base32 encoding of src.
//
//encode |--------|--------|--------|--------|--------|--------|--------|...
//       |channel |activityId|<-         index data     ->|
//       |   8    |    10    |<-             26         ->|
func (enc *Encoding) EncodePartsToString(channel uint64, activityId uint64, index uint64) string {

	if channel >= 1<<8 {
		return ""
	}

	if activityId >= 1<<10 {
		return ""
	}

	if index >= 1<<26 {
		return ""
	}

	var src = index
	src += channel << 36
	src += activityId << 26

	return enc.EncodeToString(src)
}

// DecodeString returns the bytes represented by the base32 string s.
func (enc *Encoding) DecodePartsString(coupon string) (uint64, uint64, uint64, error) {

	code, err := enc.DecodeString(coupon)
	if err != nil {
		return 0, 0, 0, err
	}

	channel := code >> 36
	activityId := code << 28 >> 54
	index := code << 38 >> 38

	return channel, activityId, index, nil
}

// EncodedLen returns the length in bytes of the base32 encoding
// of an input buffer of length n.
func (enc *Encoding) EncodedLen() int {

	return enc.encodedLength
}

/*
 * Decoder
 */
// decode is like Decode but returns an additional 'end' value, which
// indicates if end-of-message padding was encountered and thus any
// additional data is an error. This method assumes that src has been
// stripped of all supported whitespace ('\r' and '\n').
func (enc *Encoding) decode(code []byte) (n uint64, err error) {

	var parts []byte
	s := 0
	for i, v := range enc.partLens {

		var part = code[s : s+v]
		s += v

		//noinspection GoBoolExpressions
		if enc.digitCheckingSize > 0 {
			if part[v-enc.digitCheckingSize] != checkDigitAlg1(string(part[0:v-enc.digitCheckingSize]), i+enc.digitCheckingSize) {

				return 0, errors.New("invalid checksum")
			}
		}

		parts = append(parts, part[0:v-enc.digitCheckingSize]...)
	}

	if enc.offsetBits > 0 {
		parts = append(parts, 'A')
	}

	encrypted := make([]byte, enc.msgLength)
	_, err = enc.readableEncoder.Decode(encrypted, parts)

	if err != nil {
		return 0, err
	}

	if enc.digitCheckingSize == 0 {
		//fmt.Printf("decode shift %d %v\n", int(encrypted[enc.msgLength-1]), []byte(encrypted))

		//enc.shift(encrypted, 0, enc.msgLength-1, int(encrypted[enc.msgLength-1]))
		enc.xor(encrypted[0:enc.msgLength-1], encrypted[enc.msgLength-1])

		//fmt.Printf("       shift %d %v\n", int(encrypted[enc.msgLength-1]), []byte(encrypted))
	}

	// Decrypt
	decrypted := make([]byte, len(encrypted))
	var iv = []byte(enc.key)[:aes.BlockSize] // Using IV same as key is probably bad
	err = DecryptAESCFB(decrypted, encrypted[0:enc.msgLength-enc.offsetByte], []byte(enc.key), iv)
	if err != nil {
		return 0, err
	}
	//fmt.Printf("Decrypting %v %s-> %v %s\n", encrypted, parts, decrypted, decrypted)

	if enc.offsetBits > 0 {
		decrypted[enc.msgLength-1] = encrypted[enc.msgLength-1]
	}

	crc := crc8.Checksum(decrypted[1:], enc.crc8Table)
	if crc != decrypted[0] {
		return 0, errors.New("invalid checksum 8")
	}

	var data, _ = binary.Uvarint(decrypted[1:])

	return uint64(data), nil
}

// DecodeString returns the bytes represented by the base32 string s.
func (enc *Encoding) DecodeString(code string) (uint64, error) {

	if len(code) == 0 {
		return 0, errors.New("no string to decode")
	}

	// uppercase the code, take out any random chars and replace OIZS with 0125
	code = strings.ToUpper(code)
	code = strings.Replace(code, "O", "0", -1)
	code = strings.Replace(code, "I", "1", -1)
	code = strings.Replace(code, "Z", "2", -1)
	code = strings.Replace(code, "S", "5", -1)

	re := regexp.MustCompile("[^0-9A-Z]+")
	code = re.ReplaceAllString(code, "")

	return enc.decode([]byte(code))
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base32-encoded data.
func (enc *Encoding) DecodedLen() int {

	return enc.decodedLength
}

func Benchmark() {

	var wg sync.WaitGroup

	var cpu = runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	fmt.Printf("Benchmark called: %d\n", cpu)

	wg.Add(cpu)

	var n = 268435456
	//var n = 1234567

	for i := 0; i < n; i += int(math.Ceil(float64(n) / float64(cpu))) {

		var s = i
		var e = i + int(math.Ceil(float64(n)/float64(cpu)))
		if e > n {
			e = n
		}

		go func(start, end int) {
			defer wg.Done()

			fmt.Printf("Coupon Benchmark from %v to %v\n", start, end)

			rand.Seed(0)
			enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 1, 4, 4, 4)

			for j := start; j < end; j++ {

				got := enc.EncodeToString(uint64(j))

				dec, err := enc.DecodeString(got)
				if uint64(j) != dec {
					_ = fmt.Errorf("TestEncoding_EncodeString_DecodeString failed. code = %v, got = %v, want %v", j, got, dec)
				}

				if err != nil {
					_ = fmt.Errorf("Encoding.DecodeString() = %v, want %v", got, uint64(j))
				}
			}
		}(s, e)
	}

	wg.Wait()
}
