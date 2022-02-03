package log

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type username string

const (
	key = "my-key"
)

func (n username) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("username", string(n))
	return nil
}

func TestFields(t *testing.T) {
	now := time.Now()
	address := net.ParseIP("1.2.3.4")
	name := username("phil")

	testCases := []struct {
		name          string
		field         Field
		expectedField Field
	}{
		{
			name:          "Skip",
			field:         Skip(),
			expectedField: zap.Skip(),
		},
		{
			name:          "Binary",
			field:         Binary(key, []byte("ab12")),
			expectedField: zap.Binary(key, []byte("ab12")),
		},
		{
			name:          "Bool",
			field:         Bool(key, true),
			expectedField: zap.Bool(key, true),
		},
		{
			name:          "ByteString",
			field:         ByteString(key, []byte("ab12")),
			expectedField: zap.ByteString(key, []byte("ab12")),
		},
		{
			name:          "Complex128",
			field:         Complex128(key, 1+2i),
			expectedField: zap.Complex128(key, 1+2i),
		},
		{
			name:          "Complex64",
			field:         Complex64(key, 1+2i),
			expectedField: zap.Complex64(key, 1+2i),
		},
		{
			name:          "Float64",
			field:         Float64(key, 3.14),
			expectedField: zap.Float64(key, 3.14),
		},
		{
			name:          "Float32",
			field:         Float32(key, 3.14),
			expectedField: zap.Float32(key, 3.14),
		},
		{
			name:          "Int",
			field:         Int(key, 1),
			expectedField: zap.Int(key, 1),
		},
		{
			name:          "Int8",
			field:         Int8(key, 1),
			expectedField: zap.Int8(key, 1),
		},
		{
			name:          "Int16",
			field:         Int16(key, 1),
			expectedField: zap.Int16(key, 1),
		},
		{
			name:          "Int32",
			field:         Int32(key, 1),
			expectedField: zap.Int32(key, 1),
		},
		{
			name:          "Int64",
			field:         Int64(key, 1),
			expectedField: zap.Int64(key, 1),
		},
		{
			name:          "String",
			field:         String(key, "value"),
			expectedField: zap.String(key, "value"),
		},
		{
			name:          "Uint",
			field:         Uint(key, 1),
			expectedField: zap.Uint(key, 1),
		},
		{
			name:          "Uint64",
			field:         Uint64(key, 1),
			expectedField: zap.Uint64(key, 1),
		},
		{
			name:          "Uint32",
			field:         Uint32(key, 1),
			expectedField: zap.Uint32(key, 1),
		},
		{
			name:          "Uint16",
			field:         Uint16(key, 1),
			expectedField: zap.Uint16(key, 1),
		},
		{
			name:          "Uint8",
			field:         Uint8(key, 1),
			expectedField: zap.Uint8(key, 1),
		},
		{
			name:          "UintPtr",
			field:         Uintptr(key, 10),
			expectedField: zap.Uintptr(key, 0xa),
		},
		{
			name:          "Reflect",
			field:         Reflect(key, []int{5, 6}),
			expectedField: zap.Reflect(key, []int{5, 6}),
		},
		{
			name:          "Namespace",
			field:         Namespace(key),
			expectedField: zap.Namespace(key),
		},
		{
			name:          "Stringer",
			field:         Stringer(key, address),
			expectedField: zap.Stringer(key, address),
		},
		{
			name:          "Time",
			field:         Time(key, time.Unix(0, 1000).In(time.UTC)),
			expectedField: zap.Time(key, time.Unix(0, 1000).In(time.UTC)),
		},
		{
			name:          "Duration",
			field:         Duration(key, 1),
			expectedField: zap.Duration(key, 1),
		},
		{
			name:          "Any:ObjectMarshaler",
			field:         Any(key, name),
			expectedField: zap.Any(key, name),
		},
		{
			name:          "Any:ArrayMarshaler",
			field:         Any(key, []bool{true}),
			expectedField: zap.Any(key, []bool{true}),
		},
		{
			name:          "Any:Stringer",
			field:         Any(key, address),
			expectedField: zap.Any(key, address),
		},
		{
			name:          "Any:Bool",
			field:         Any(key, true),
			expectedField: zap.Any(key, true),
		},
		{
			name:          "Any:Bools",
			field:         Any(key, []bool{true}),
			expectedField: zap.Any(key, []bool{true}),
		},
		{
			name:          "Any:Byte",
			field:         Any(key, byte(1)),
			expectedField: zap.Any(key, byte(1)),
		},
		{
			name:          "Any:Bytes",
			field:         Any(key, []byte{1}),
			expectedField: zap.Any(key, []byte{1}),
		},
		{
			name:          "Any:Complex128",
			field:         Any(key, 1+2i),
			expectedField: zap.Any(key, 1+2i),
		},
		{
			name:          "Any:Complex128s",
			field:         Any(key, []complex128{1 + 2i}),
			expectedField: zap.Any(key, []complex128{1 + 2i}),
		},
		{
			name:          "Any:Complex64",
			field:         Any(key, complex64(1+2i)),
			expectedField: zap.Any(key, complex64(1+2i)),
		},
		{
			name:          "Any:Complex64s",
			field:         Any(key, []complex64{1 + 2i}),
			expectedField: zap.Any(key, []complex64{1 + 2i}),
		},
		{
			name:          "Any:Float64",
			field:         Any(key, 3.14),
			expectedField: zap.Any(key, 3.14),
		},
		{
			name:          "Any:Float64s",
			field:         Any(key, []float64{3.14}),
			expectedField: zap.Any(key, []float64{3.14}),
		},
		{
			name:          "Any:Float32",
			field:         Any(key, float32(3.14)),
			expectedField: zap.Any(key, float32(3.14)),
		},
		{
			name:          "Any:Float32s",
			field:         Any(key, []float32{3.14}),
			expectedField: zap.Any(key, []float32{3.14}),
		},
		{
			name:          "Any:Int",
			field:         Any(key, 1),
			expectedField: zap.Any(key, 1),
		},
		{
			name:          "Any:Ints",
			field:         Any(key, []int{1}),
			expectedField: zap.Any(key, []int{1}),
		},
		{
			name:          "Any:Int64",
			field:         Any(key, int64(1)),
			expectedField: zap.Any(key, int64(1)),
		},
		{
			name:          "Any:Int64s",
			field:         Any(key, []int64{1}),
			expectedField: zap.Any(key, []int64{1}),
		},
		{
			name:          "Any:Int32",
			field:         Any(key, int32(1)),
			expectedField: zap.Any(key, int32(1)),
		},
		{
			name:          "Any:Int32s",
			field:         Any(key, []int32{1}),
			expectedField: zap.Any(key, []int32{1}),
		},
		{
			name:          "Any:Int16",
			field:         Any(key, int16(1)),
			expectedField: zap.Any(key, int16(1)),
		},
		{
			name:          "Any:Int16s",
			field:         Any(key, []int16{1}),
			expectedField: zap.Any(key, []int16{1}),
		},
		{
			name:          "Any:Int8",
			field:         Any(key, int8(1)),
			expectedField: zap.Any(key, int8(1)),
		},
		{
			name:          "Any:Int8s",
			field:         Any(key, []int8{1}),
			expectedField: zap.Any(key, []int8{1}),
		},
		{
			name:          "Any:Rune",
			field:         Any(key, rune(1)),
			expectedField: zap.Any(key, rune(1)),
		},
		{
			name:          "Any:Runes",
			field:         Any(key, []rune{1}),
			expectedField: zap.Any(key, []rune{1}),
		},
		{
			name:          "Any:String",
			field:         Any(key, "s"),
			expectedField: zap.Any(key, "s"),
		},
		{
			name:          "Any:Strings",
			field:         Any(key, []string{"s"}),
			expectedField: zap.Any(key, []string{"s"}),
		},
		{
			name:          "Any:Uint",
			field:         Any(key, uint(1)),
			expectedField: zap.Any(key, uint(1)),
		},
		{
			name:          "Any:Uints",
			field:         Any(key, []uint{1}),
			expectedField: zap.Any(key, []uint{1}),
		},
		{
			name:          "Any:Uint64",
			field:         Any(key, uint64(1)),
			expectedField: zap.Any(key, uint64(1)),
		},
		{
			name:          "Any:Uint64s",
			field:         Any(key, []uint64{1}),
			expectedField: zap.Any(key, []uint64{1}),
		},
		{
			name:          "Any:Uint32",
			field:         Any(key, uint32(1)),
			expectedField: zap.Any(key, uint32(1)),
		},
		{
			name:          "Any:Uint32s",
			field:         Any(key, []uint32{1}),
			expectedField: zap.Any(key, []uint32{1}),
		},
		{
			name:          "Any:Uint16",
			field:         Any(key, uint16(1)),
			expectedField: zap.Any(key, uint16(1)),
		},
		{
			name:          "Any:Uint16s",
			field:         Any(key, []uint16{1}),
			expectedField: zap.Any(key, []uint16{1}),
		},
		{
			name:          "Any:Uint8",
			field:         Any(key, uint8(1)),
			expectedField: zap.Any(key, uint8(1)),
		},
		{
			name:          "Any:Uint8s",
			field:         Any(key, []uint8{1}),
			expectedField: zap.Any(key, []uint8{1}),
		},
		{
			name:          "Err",
			field:         Err(fmt.Errorf("custom error")),
			expectedField: zap.Error(fmt.Errorf("custom error")),
		},
		{
			name:          "NamedErr",
			field:         NamedErr(key, fmt.Errorf("custom error")),
			expectedField: zap.NamedError(key, fmt.Errorf("custom error")),
		},
		{
			name:          "Bools",
			field:         Bools(key, []bool{true, false}),
			expectedField: zap.Bools(key, []bool{true, false}),
		},
		{
			name:          "ByteStrings",
			field:         ByteStrings(key, [][]byte{{'a', 'b'}, {'c', 'd'}}),
			expectedField: zap.ByteStrings(key, [][]byte{{'a', 'b'}, {'c', 'd'}}),
		},
		{
			name:          "Complex128s",
			field:         Complex128s(key, []complex128{10i, 5i}),
			expectedField: zap.Complex128s(key, []complex128{10i, 5i}),
		},
		{
			name:          "Complex64s",
			field:         Complex64s(key, []complex64{10i, 5i}),
			expectedField: zap.Complex64s(key, []complex64{10i, 5i}),
		},
		{
			name:          "Durations",
			field:         Durations(key, []time.Duration{time.Second, time.Minute}),
			expectedField: zap.Durations(key, []time.Duration{time.Second, time.Minute}),
		},
		{
			name:          "Float64s",
			field:         Float64s(key, []float64{1.0, 2.0, 3.0}),
			expectedField: zap.Float64s(key, []float64{1.0, 2.0, 3.0}),
		},
		{
			name:          "Float32s",
			field:         Float32s(key, []float32{1.0, 2.0, 3.0}),
			expectedField: zap.Float32s(key, []float32{1.0, 2.0, 3.0}),
		},
		{
			name:          "Ints",
			field:         Ints(key, []int{1, 2, 3}),
			expectedField: zap.Ints(key, []int{1, 2, 3}),
		},
		{
			name:          "Int64s",
			field:         Int64s(key, []int64{1, 2, 3}),
			expectedField: zap.Int64s(key, []int64{1, 2, 3}),
		},
		{
			name:          "Int32s",
			field:         Int32s(key, []int32{1, 2, 3}),
			expectedField: zap.Int32s(key, []int32{1, 2, 3}),
		},
		{
			name:          "Int16s",
			field:         Int16s(key, []int16{1, 2, 3}),
			expectedField: zap.Int16s(key, []int16{1, 2, 3}),
		},
		{
			name:          "Int8s",
			field:         Int8s(key, []int8{1, 2, 3}),
			expectedField: zap.Int8s(key, []int8{1, 2, 3}),
		},
		{
			name:          "Strings",
			field:         Strings(key, []string{"first", "second"}),
			expectedField: zap.Strings(key, []string{"first", "second"}),
		},
		{
			name:          "Times",
			field:         Times(key, []time.Time{now, now}),
			expectedField: zap.Times(key, []time.Time{now, now}),
		},
		{
			name:          "Uints",
			field:         Uints(key, []uint{1, 2, 3}),
			expectedField: zap.Uints(key, []uint{1, 2, 3}),
		},
		{
			name:          "Uint64s",
			field:         Uint64s(key, []uint64{1, 2, 3}),
			expectedField: zap.Uint64s(key, []uint64{1, 2, 3}),
		},
		{
			name:          "Uint32s",
			field:         Uint32s(key, []uint32{1, 2, 3}),
			expectedField: zap.Uint32s(key, []uint32{1, 2, 3}),
		},
		{
			name:          "Uint16s",
			field:         Uint16s(key, []uint16{1, 2, 3}),
			expectedField: zap.Uint16s(key, []uint16{1, 2, 3}),
		},
		{
			name:          "Uint8s",
			field:         Uint8s(key, []uint8{1, 2, 3}),
			expectedField: zap.Uint8s(key, []uint8{1, 2, 3}),
		},
		{
			name:          "Uintptrs",
			field:         Uintptrs(key, []uintptr{uintptr(unsafe.Pointer(&address))}),
			expectedField: zap.Uintptrs(key, []uintptr{uintptr(unsafe.Pointer(&address))}),
		},
		{
			name:          "Errors",
			field:         Errors(key, []error{errors.New("my error")}),
			expectedField: zap.Errors(key, []error{errors.New("my error")}),
		},
	}

	for _, tt := range testCases {
		assert.Equal(t, tt.expectedField, tt.field, "test %s failed", tt.name)
	}
}

func TestStackField(t *testing.T) {
	f := Stack("stacktrace")
	assert.Equal(t, "stacktrace", f.Key, "Unexpected field key.")
}
