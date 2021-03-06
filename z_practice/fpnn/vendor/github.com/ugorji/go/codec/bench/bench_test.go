// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// Sample way to run:
// go test -bi -bv -bd=1 -benchmem -bench=.

const benchUnscientificRes = true
const benchVerify = true

func init() {
	testPreInitFns = append(testPreInitFns, benchPreInit)
	testPostInitFns = append(testPostInitFns, benchPostInit)
}

var (
	benchTs *TestStruc

	approxSize int

	benchCheckers []benchChecker
)

type benchEncFn func(interface{}, []byte) ([]byte, error)
type benchDecFn func([]byte, interface{}) error
type benchIntfFn func() interface{}

type benchChecker struct {
	name     string
	encodefn benchEncFn
	decodefn benchDecFn
}

func benchReinit() {
	benchCheckers = nil
}

func benchPreInit() {
	benchTs = newTestStruc(testDepth, testNumRepeatString, true, !testSkipIntf, testMapStringKeyOnly)
	approxSize = approxDataSize(reflect.ValueOf(benchTs)) * 3 / 2 // multiply by 1.5 to appease msgp, and prevent alloc
	// bytesLen := 1024 * 4 * (testDepth + 1) * (testDepth + 1)
	// if bytesLen < approxSize {
	// 	bytesLen = approxSize
	// }

	benchCheckers = append(benchCheckers,
		benchChecker{"msgpack", fnMsgpackEncodeFn, fnMsgpackDecodeFn},
		benchChecker{"binc", fnBincEncodeFn, fnBincDecodeFn},
		benchChecker{"simple", fnSimpleEncodeFn, fnSimpleDecodeFn},
		benchChecker{"cbor", fnCborEncodeFn, fnCborDecodeFn},
		benchChecker{"json", fnJsonEncodeFn, fnJsonDecodeFn},
		benchChecker{"std-json", fnStdJsonEncodeFn, fnStdJsonDecodeFn},
		benchChecker{"gob", fnGobEncodeFn, fnGobDecodeFn},
		benchChecker{"std-xml", fnStdXmlEncodeFn, fnStdXmlDecodeFn},
	)
}

func benchPostInit() {
	// if benchDoInitBench {
	// 	runBenchInit()
	// }
}

func TestBenchInit(t *testing.T) {
	testOnce.Do(testInitAll)
	// logTv(t, "..............................................")
	logT(t, "BENCHMARK INIT: %v", time.Now())
	// logTv(t, "To run full benchmark comparing encodings, use: \"go test -bench=.\"")
	logT(t, "Benchmark: ")
	logT(t, "\tStruct recursive Depth:             %d", testDepth)
	if approxSize > 0 {
		logT(t, "\tApproxDeepSize Of benchmark Struct: %d bytes", approxSize)
	}
	if benchUnscientificRes {
		logT(t, "Benchmark One-Pass Run (with Unscientific Encode/Decode times): ")
	} else {
		logT(t, "Benchmark One-Pass Run:")
	}
	for _, bc := range benchCheckers {
		doBenchCheck(t, bc.name, bc.encodefn, bc.decodefn)
	}
	logTv(t, "..............................................")
	logTv(t, "<<<<====>>>> depth: %v, ts: %#v\n", testDepth, benchTs)
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
}

var vBenchTs = TestStruc{}

func fnBenchNewTs() interface{} {
	vBenchTs = TestStruc{}
	return &vBenchTs
	// return new(TestStruc)
}

// const benchCheckDoDeepEqual = false

func benchRecoverPanic(t interface{}) {
	if r := recover(); r != nil {
		logT(t, "(recovered) panic: %v\n", r)
	}
}

func doBenchCheck(t *testing.T, name string, encfn benchEncFn, decfn benchDecFn) {
	// if benchUnscientificRes {
	// 	logTv(t, "-------------- %s ----------------", name)
	// }
	defer benchRecoverPanic(nil)
	runtime.GC()
	tnow := time.Now()
	buf, err := encfn(benchTs, nil)
	if err != nil {
		logT(t, "\t%10s: **** Error encoding benchTs: %v", name, err)
		return
	}
	encDur := time.Since(tnow)
	encLen := len(buf)
	runtime.GC()
	if !benchUnscientificRes {
		logT(t, "\t%10s: len: %d bytes\n", name, encLen)
		return
	}
	tnow = time.Now()
	var ts2 TestStruc
	if err = decfn(buf, &ts2); err != nil {
		logT(t, "\t%10s: **** Error decoding into new TestStruc: %v", name, err)
		return
	}
	decDur := time.Since(tnow)
	// if benchCheckDoDeepEqual {
	if benchVerify {
		err = deepEqual(benchTs, &ts2)
		if err == nil {
			logT(t, "\t%10s: len: %d bytes,\t encode: %v,\t decode: %v,\tencoded == decoded", name, encLen, encDur, decDur)
		} else {
			logT(t, "\t%10s: len: %d bytes,\t encode: %v,\t decode: %v,\tencoded != decoded: %v", name, encLen, encDur, decDur, err)
			// if strings.Contains(name, "json") {
			// 	println(">>>>>")
			// 	f1, _ := os.Create("1.out")
			// 	f2, _ := os.Create("2.out")
			// 	f3, _ := os.Create("3.json")
			// 	buf3, _ := json.MarshalIndent(&ts2, "", "\t")
			// 	spew.Config.SortKeys = true
			// 	spew.Config.SpewKeys = true
			// 	println("^^^^^^^^^^^^^^")
			// 	spew.Fdump(f1, benchTs)
			// 	println("^^^^^^^^^^^^^^")
			// 	spew.Fdump(f2, &ts2)
			// 	println("^^^^^^^^^^^^^^")
			// 	f3.Write(buf3)
			// 	f1.Close()
			// 	f2.Close()
			// 	f3.Close()
			// }
			// logT(t, "\t: err: %v,\n benchTs: %#v\n\n, ts2: %#v\n\n", err, benchTs, ts2) // TODO: remove
			// logT(t, "BenchVerify: Error comparing en|decoded TestStruc: %v", err)
			// return
			// logT(t, "BenchVerify: Error comparing benchTs: %v\n--------\n%v\n--------\n%v", err, benchTs, ts2)
			// if strings.Contains(name, "json") {
			// 	logT(t, "\n\tDECODED FROM\n--------\n%s", buf)
			// }
		}
	} else {
		logT(t, "\t%10s: len: %d bytes,\t encode: %v,\t decode: %v", name, encLen, encDur, decDur)
	}
	return
}

func fnBenchmarkEncode(b *testing.B, encName string, ts interface{}, encfn benchEncFn) {
	defer benchRecoverPanic(b)
	testOnce.Do(testInitAll)
	var err error
	bs := make([]byte, 0, approxSize)
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = encfn(ts, bs); err != nil {
			break
		}
	}
	if err != nil {
		failT(b, "Error encoding benchTs: %s: %v", encName, err)
	}
}

func fnBenchmarkDecode(b *testing.B, encName string, ts interface{},
	encfn benchEncFn, decfn benchDecFn, newfn benchIntfFn,
) {
	defer benchRecoverPanic(b)
	testOnce.Do(testInitAll)
	bs := make([]byte, 0, approxSize)
	buf, err := encfn(ts, bs)
	if err != nil {
		failT(b, "Error encoding benchTs: %s: %v", encName, err)
	}
	// if false && benchVerify { // do not do benchVerify during decode
	// 	// ts2 := newfn()
	// 	ts1 := ts.(*TestStruc)
	// 	ts2 := new(TestStruc)
	// 	if err = decfn(buf, ts2); err != nil {
	// 		failT(b, "BenchVerify: Error decoding benchTs: %s: %v", encName, err)
	// 	}
	// 	if err = deepEqual(ts1, ts2); err != nil {
	// 		failT(b, "BenchVerify: Error comparing benchTs: %s: %v", encName, err)
	// 	}
	// }
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts = newfn()
		if err = decfn(buf, ts); err != nil {
			break
		}
	}
	if err != nil {
		failT(b, "Error decoding into new TestStruc: %s: %v", encName, err)
	}
}

// ------------ tests below

func fnMsgpackEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return sTestCodecEncode(ts, bsIn, fnBenchmarkByteBuf, testMsgpackH, &testMsgpackH.BasicHandle)
}

func fnMsgpackDecodeFn(buf []byte, ts interface{}) error {
	return sTestCodecDecode(buf, ts, testMsgpackH, &testMsgpackH.BasicHandle)
}

func fnBincEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return sTestCodecEncode(ts, bsIn, fnBenchmarkByteBuf, testBincH, &testBincH.BasicHandle)
}

func fnBincDecodeFn(buf []byte, ts interface{}) error {
	return sTestCodecDecode(buf, ts, testBincH, &testBincH.BasicHandle)
}

func fnSimpleEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return sTestCodecEncode(ts, bsIn, fnBenchmarkByteBuf, testSimpleH, &testSimpleH.BasicHandle)
}

func fnSimpleDecodeFn(buf []byte, ts interface{}) error {
	return sTestCodecDecode(buf, ts, testSimpleH, &testSimpleH.BasicHandle)
}

func fnCborEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return sTestCodecEncode(ts, bsIn, fnBenchmarkByteBuf, testCborH, &testCborH.BasicHandle)
}

func fnCborDecodeFn(buf []byte, ts interface{}) error {
	return sTestCodecDecode(buf, ts, testCborH, &testCborH.BasicHandle)
}

func fnJsonEncodeFn(ts interface{}, bsIn []byte) (bs []byte, err error) {
	return sTestCodecEncode(ts, bsIn, fnBenchmarkByteBuf, testJsonH, &testJsonH.BasicHandle)
}

func fnJsonDecodeFn(buf []byte, ts interface{}) error {
	return sTestCodecDecode(buf, ts, testJsonH, &testJsonH.BasicHandle)
}

func fnGobEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	buf := fnBenchmarkByteBuf(bsIn)
	err := gob.NewEncoder(buf).Encode(ts)
	return buf.Bytes(), err
}

func fnGobDecodeFn(buf []byte, ts interface{}) error {
	return gob.NewDecoder(bytes.NewReader(buf)).Decode(ts)
}

func fnStdXmlEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	buf := fnBenchmarkByteBuf(bsIn)
	err := xml.NewEncoder(buf).Encode(ts)
	return buf.Bytes(), err
}

func fnStdXmlDecodeFn(buf []byte, ts interface{}) error {
	return xml.NewDecoder(bytes.NewReader(buf)).Decode(ts)
}

func fnStdJsonEncodeFn(ts interface{}, bsIn []byte) ([]byte, error) {
	if testUseIoEncDec >= 0 {
		buf := fnBenchmarkByteBuf(bsIn)
		err := json.NewEncoder(buf).Encode(ts)
		return buf.Bytes(), err
	}
	return json.Marshal(ts)
}

func fnStdJsonDecodeFn(buf []byte, ts interface{}) error {
	if testUseIoEncDec >= 0 {
		return json.NewDecoder(bytes.NewReader(buf)).Decode(ts)
	}
	return json.Unmarshal(buf, ts)
}

// ----------- DECODE ------------------

func Benchmark__Msgpack____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "msgpack", benchTs, fnMsgpackEncodeFn)
}

func Benchmark__Binc_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "binc", benchTs, fnBincEncodeFn)
}

func Benchmark__Simple_____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "simple", benchTs, fnSimpleEncodeFn)
}

func Benchmark__Cbor_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "cbor", benchTs, fnCborEncodeFn)
}

func Benchmark__Json_______Encode(b *testing.B) {
	fnBenchmarkEncode(b, "json", benchTs, fnJsonEncodeFn)
}

func Benchmark__Std_Json___Encode(b *testing.B) {
	fnBenchmarkEncode(b, "std-json", benchTs, fnStdJsonEncodeFn)
}

func Benchmark__Gob________Encode(b *testing.B) {
	fnBenchmarkEncode(b, "gob", benchTs, fnGobEncodeFn)
}

func Benchmark__Std_Xml____Encode(b *testing.B) {
	fnBenchmarkEncode(b, "std-xml", benchTs, fnStdXmlEncodeFn)
}

// ----------- DECODE ------------------

func Benchmark__Msgpack____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "msgpack", benchTs, fnMsgpackEncodeFn, fnMsgpackDecodeFn, fnBenchNewTs)
}

func Benchmark__Binc_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "binc", benchTs, fnBincEncodeFn, fnBincDecodeFn, fnBenchNewTs)
}

func Benchmark__Simple_____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "simple", benchTs, fnSimpleEncodeFn, fnSimpleDecodeFn, fnBenchNewTs)
}

func Benchmark__Cbor_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "cbor", benchTs, fnCborEncodeFn, fnCborDecodeFn, fnBenchNewTs)
}

func Benchmark__Json_______Decode(b *testing.B) {
	fnBenchmarkDecode(b, "json", benchTs, fnJsonEncodeFn, fnJsonDecodeFn, fnBenchNewTs)
}

func Benchmark__Std_Json___Decode(b *testing.B) {
	fnBenchmarkDecode(b, "std-json", benchTs, fnStdJsonEncodeFn, fnStdJsonDecodeFn, fnBenchNewTs)
}

func Benchmark__Gob________Decode(b *testing.B) {
	fnBenchmarkDecode(b, "gob", benchTs, fnGobEncodeFn, fnGobDecodeFn, fnBenchNewTs)
}

func Benchmark__Std_Xml____Decode(b *testing.B) {
	fnBenchmarkDecode(b, "std-xml", benchTs, fnStdXmlEncodeFn, fnStdXmlDecodeFn, fnBenchNewTs)
}
