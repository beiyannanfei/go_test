package coupon

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestNewEncoding(t *testing.T) {
	type args struct {
		encoder           string
		key               string
		dataByteSize      int
		digitCheckingSize int
		partLens          []int
	}

	var tests []struct {
		name string
		args args
		want *Encoding
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEncoding(tt.args.key, tt.args.dataByteSize, tt.args.digitCheckingSize, tt.args.partLens...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_Encode(t *testing.T) {
	type fields struct {
		encode              string
		key                 string
		partLens            []int
		dataByteSize        int
		availableCharLength int
		availableBitLength  int
		offsetBits          int
		msgLength           int
		offsetByte          int
		encodedLength       int
	}

	type args struct {
		dst []byte
		src uint64
	}

	var tests []struct {
		name   string
		fields fields
		args   args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &Encoding{
				key:                 tt.fields.key,
				partLens:            tt.fields.partLens,
				dataByteSize:        tt.fields.dataByteSize,
				availableCharLength: tt.fields.availableCharLength,
				availableBitLength:  tt.fields.availableBitLength,
				offsetBits:          tt.fields.offsetBits,
				msgLength:           tt.fields.msgLength,
				offsetByte:          tt.fields.offsetByte,
				encodedLength:       tt.fields.encodedLength,
			}
			enc.encode(tt.args.dst, tt.args.src)
		})
	}
}

func TestEncoding_EncodeToString(t *testing.T) {

	type args struct {
		src uint64
	}

	//noinspection SpellCheckingInspection
	var tests = []struct {
		name string
		args args
		want string
	}{
		{"", args{0xFFFF}, "5CB6-WED7-87X6"},
	}

	enc := NewEncoding("1234567890123456", 4, 1, 4, 4, 4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rand.Seed(0)
			if got := enc.EncodeToString(tt.args.src); got != tt.want {
				t.Errorf("Encoding.EncodeToString(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func BenchmarkEncoding_EncodeToString(b *testing.B) {

	rand.Seed(0)
	enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 1, 4, 4, 4)

	for i := 0; i < b.N; i++ {
		enc.EncodeToString(uint64(i))
		//var result = enc.EncodeToString(uint64(i))
		//fmt.Printf("EncodeToString %v -> %v\n", i, result)
	}
}

func TestEncoding_EncodedLen(t *testing.T) {

	type fields struct {
		encode              string
		key                 string
		partLens            []int
		dataByteSize        int
		availableCharLength int
		availableBitLength  int
		offsetBits          int
		msgLength           int
		offsetByte          int
		encodedLength       int
	}

	var tests []struct {
		name   string
		fields fields
		want   int
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &Encoding{
				key:                 tt.fields.key,
				partLens:            tt.fields.partLens,
				dataByteSize:        tt.fields.dataByteSize,
				availableCharLength: tt.fields.availableCharLength,
				availableBitLength:  tt.fields.availableBitLength,
				offsetBits:          tt.fields.offsetBits,
				msgLength:           tt.fields.msgLength,
				offsetByte:          tt.fields.offsetByte,
				encodedLength:       tt.fields.encodedLength,
			}
			if got := enc.EncodedLen(); got != tt.want {
				t.Errorf("Encoding.EncodedLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_decode(t *testing.T) {

	type fields struct {
		encode              string
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
	}

	type args struct {
		code []byte
	}

	var tests []struct {
		name    string
		fields  fields
		args    args
		wantN   uint64
		wantErr bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &Encoding{
				key:                 tt.fields.key,
				partLens:            tt.fields.partLens,
				dataByteSize:        tt.fields.dataByteSize,
				availableCharLength: tt.fields.availableCharLength,
				availableBitLength:  tt.fields.availableBitLength,
				offsetBits:          tt.fields.offsetBits,
				msgLength:           tt.fields.msgLength,
				offsetByte:          tt.fields.offsetByte,
				encodedLength:       tt.fields.encodedLength,
				decodedLength:       tt.fields.decodedLength,
			}
			gotN, err := enc.decode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoding.decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Encoding.decode() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestEncoding_DecodeString(t *testing.T) {

	type args struct {
		code string
	}

	//noinspection SpellCheckingInspection
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"", args{"5CB6-WED7-87X6"}, 0xFFFF, false},
		//{"", args{"T94M-CR4N-T4PP"}, 0xFFFF, false},
	}

	enc := NewEncoding("1234567890123456", 4, 1, 4, 4, 4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := enc.DecodeString(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoding.DecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encoding.DecodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_Encode_Decode(t *testing.T) {

	type args struct {
		code   uint64
		coupon string
	}

	//noinspection SpellCheckingInspection
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{0xFFFF, "Y3Y7Y-33YE-CC6B4"}, false},
		{"", args{17592186044373, "QTL26-W69D-9KYMR"}, false},
		{"", args{17592186044335, "2TU6E-A8WL-QRUW6"}, false},
	}

	rand.Seed(0)
	enc := NewEncoding("1234567890123456", 4, 0, 5, 4, 5)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := enc.EncodeToString(tt.args.code)
			if got != tt.args.coupon {
				t.Errorf("Encoding.EncodeToString(%v) = %v, want %v", tt.args.code, got, tt.args.coupon)
			}

			code, err := enc.DecodeString(tt.args.coupon)
			if code != tt.args.code {
				t.Errorf("Encoding.EncodeToString(%v) = %v, want %v", tt.args.coupon, code, tt.args.code)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Encoding.DecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestEncoding_DecodedLen(t *testing.T) {
	type fields struct {
		encode              string
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
	}

	var tests []struct {
		name   string
		fields fields
		want   int
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &Encoding{
				key:                 tt.fields.key,
				partLens:            tt.fields.partLens,
				dataByteSize:        tt.fields.dataByteSize,
				availableCharLength: tt.fields.availableCharLength,
				availableBitLength:  tt.fields.availableBitLength,
				offsetBits:          tt.fields.offsetBits,
				msgLength:           tt.fields.msgLength,
				offsetByte:          tt.fields.offsetByte,
				encodedLength:       tt.fields.encodedLength,
				decodedLength:       tt.fields.decodedLength,
			}
			if got := enc.DecodedLen(); got != tt.want {
				t.Errorf("Encoding.DecodedLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_EncodeString_DecodeString(t *testing.T) {

	//enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 1, 4, 4, 4)
	enc := NewEncoding("BenchmarkEncoding_EncodeToString", 6, 0, 5, 4, 5)

	//n := 10000
	//n := 268435456
	//n := 4294967296
	var max uint64 = 0x100000000000
	var min uint64 = 0x0FFFFFFFFF00

	for i := max; i > min; i-- {

		got := enc.EncodeToString(i)

		dec, err := enc.DecodeString(got)

		t.Logf("Encoding %d EncodeToString: %s", i, got)

		if i != dec {
			t.Errorf("TestEncoding_EncodeString_DecodeString failed. code = %v, got = %v, want %v", i, got, dec)
		}

		if err != nil {

			t.Errorf("Encoding.DecodeString() = %v, want %v", got, uint64(i))
		}
	}
}

func BenchmarkEncoding_EncodeString_DecodeString(b *testing.B) {

	rand.Seed(0)
	enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 1, 4, 4, 4)

	for i := 0; i < b.N; i++ {

		got := enc.EncodeToString(uint64(i))

		dec, err := enc.DecodeString(got)
		if uint64(i) != dec {
			b.Errorf("TestEncoding_EncodeString_DecodeString failed. code = %v, got = %v, want %v", i, got, dec)
		}

		if err != nil {

			b.Errorf("Encoding.DecodeString() = %v, want %v", got, uint64(i))
		}
	}
}

func BenchmarkParallelEncoding_EncodeString_DecodeString(b *testing.B) {
	rand.Seed(0)
	enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 1, 4, 4, 4)

	i := 4

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			got := enc.EncodeToString(uint64(i))

			dec, err := enc.DecodeString(got)
			if uint64(i) != dec {
				b.Errorf("TestEncoding_EncodeString_DecodeString failed. code = %v, got = %v, want %v", i, got, dec)
			}

			if err != nil {

				b.Errorf("Encoding.DecodeString() = %v, want %v", got, uint64(i))
			}
		}
	})
}

func BenchmarkParallelEncoding_EncodeString_DecodeString_Long(b *testing.B) {

	rand.Seed(0)
	enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 0, 5, 4, 5)

	var i uint64 = 0x100000000000

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			func(code uint64) {

				var got = enc.EncodeToString(code)

				var dec, err = enc.DecodeString(got)

				if code != dec {
					b.Errorf("TestEncoding_EncodeString_DecodeString failed. code = %v encode = %v, decode = %v", code, got, dec)
				}

				if err != nil {

					b.Errorf("Encoding.DecodeString() = %v, want %v", got, code)
				}

			}(i)

			i -= 1
		}
	})
}

func TestEncoding_EncodeToStringTable(t *testing.T) {

	type args struct {
		src uint64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{268435455}, ""},
	}

	enc := NewEncoding("1234567890123456", 4, 1, 4, 4, 4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)

			got := enc.EncodeToString(tt.args.src)

			dec, err := enc.DecodeString(got)

			if uint64(tt.args.src) != dec {
				t.Errorf("TestEncoding_EncodeString_DecodeString failed. code = %v, got = %v, want %v", dec, got, tt.args.src)
			}

			if err != nil {

				t.Errorf("Encoding.DecodeString() = %v, want %v", got, tt.args.src)
			}
		})
	}
}

func TestEncoding_reverse(t *testing.T) {
	type args struct {
		src   []byte
		start int
		end   int
	}

	enc := NewEncoding("1234567890123456", 4, 1, 4, 4, 4)

	tests := []struct {
		name string
		enc  *Encoding
		args args
		want []byte
	}{
		{
			name: "",
			enc:  enc,
			args: args{
				src:   []byte{1, 2, 3},
				start: 0,
				end:   2,
			},
			want: []byte{3, 2, 1},
		},
		{
			name: "",
			enc:  enc,
			args: args{
				src:   []byte{1, 2, 3, 4, 5, 6, 7},
				start: 0,
				end:   2,
			},
			want: []byte{3, 2, 1, 4, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.reverse(tt.args.src, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoding.reverse(%v, %v, %v) = %v, want %v", tt.args.src, tt.args.start, tt.args.end, got, tt.want)
			}
		})
	}
}

func TestEncoding_encode(t *testing.T) {
	type args struct {
		dst []byte
		src uint64
	}

	var tests []struct {
		name string
		enc  *Encoding
		args args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.enc.encode(tt.args.dst, tt.args.src)
		})
	}
}

func TestEncoding_shift(t *testing.T) {
	type args struct {
		src   []byte
		start int
		end   int
		k     int
	}

	enc := NewEncoding("1234567890123456", 4, 1, 4, 4, 4)

	tests := []struct {
		name string
		enc  *Encoding
		args args
		want []byte
	}{
		{
			name: "",
			enc:  enc,
			args: args{
				src:   []byte{1, 2, 3, 4, 5, 6, 7},
				start: 0,
				end:   6,
				k:     14,
			},
			want: []byte{3, 4, 5, 6, 1, 2, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.shift(tt.args.src, tt.args.start, tt.args.end, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoding.shift(%v, %v, %v, %v) = %v, want %v", tt.args.src, tt.args.start, tt.args.end, tt.args.k, got, tt.want)
			}
		})
	}
}

func TestEncoding_xor(t *testing.T) {
	type args struct {
		src  []byte
		data byte
	}

	enc := NewEncoding("1234567890123456", 4, 1, 4, 4, 4)

	tests := []struct {
		name string
		enc  *Encoding
		args args
		want []byte
	}{
		{
			name: "",
			enc:  enc,
			args: args{
				src:  nil,
				data: 0,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.xor(tt.args.src, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoding.xor(%v, %v) = %v, want %v", tt.args.src, tt.args.data, got, tt.want)
			}
		})
	}
}

func TestEncoding_EncodePartsToString(t *testing.T) {
	type args struct {
		channel    uint64
		activityId uint64
		index      uint64
	}

	rand.Seed(0)
	enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 0, 5, 4, 5)

	//noinspection SpellCheckingInspection
	tests := []struct {
		name string
		enc  *Encoding
		args args
		want string
	}{
		{"", enc, args{255, 1023, 0x3FFFFFF}, "YQTQ7-73J9-AV40B"},
		{"", enc, args{1, 1, 1}, "821GC-8EHQ-6XN37"},
		{"", enc, args{128, 512, 0x1FFFFFF}, "TV0QF-84XU-PM4NC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.EncodePartsToString(tt.args.channel, tt.args.activityId, tt.args.index); got != tt.want {
				t.Errorf("Encoding.EncodePartsToString(%v, %v, %v) = %v, want %v",
					tt.args.channel, tt.args.activityId, tt.args.index, got, tt.want)
			}
		})
	}
}

func TestEncoding_DecodePartsString(t *testing.T) {
	type args struct {
		coupon string
	}

	rand.Seed(0)
	enc := NewEncoding("BenchmarkEncoding_EncodeToString", 4, 0, 5, 4, 5)

	//noinspection SpellCheckingInspection
	tests := []struct {
		name    string
		enc     *Encoding
		args    args
		want    uint64
		want1   uint64
		want2   uint64
		wantErr bool
	}{
		{"", enc, args{coupon: "YQTQ7-73J9-AV40B",}, 255, 1023, 67108863, false,},
		{"", enc, args{coupon: "821GC-8EHQ-6XN37",}, 1, 1, 1, false,},
		{"", enc, args{coupon: "TV0QF-84XU-PM4NC",}, 128, 512, 0x1FFFFFF, false,},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := tt.enc.DecodePartsString(tt.args.coupon)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoding.DecodePartsString(%v) error = %v, wantErr %v", tt.args.coupon, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encoding.DecodePartsString(%v) got = %v, want %v", tt.args.coupon, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Encoding.DecodePartsString(%v) got1 = %v, want %v", tt.args.coupon, got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("Encoding.DecodePartsString(%v) got2 = %v, want %v", tt.args.coupon, got2, tt.want2)
			}
		})
	}
}
