package helpers

import (
	"errors"
	"reflect"
	"unsafe"
)

type eface struct {
	typ, val unsafe.Pointer
}

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
type integer interface {
	signed | unsigned
}

func Sizeof[T any](v T) int {
	return int(reflect.TypeOf(v).Size())
}

func Bytesof(v interface{}) []byte {
	return unsafe.Slice((*byte)((*eface)(unsafe.Pointer(&v)).val), Sizeof(v))
}

func Frombytes[T any](srcBytes []byte, dst *T) {
	dstBytes := make([]byte, Sizeof(*dst))
	copy(dstBytes, srcBytes)
	*dst = *(*T)(unsafe.Pointer(&dstBytes[0]))
}

func Convert[T integer](from interface{}, to *T) T {
	srcSize := Sizeof(from)
	dstSize := Sizeof(*to)

	if srcSize >= dstSize {
		*to = *(*T)((*eface)(unsafe.Pointer(&from)).val)
		return *to
	}

	Frombytes(Bytesof(from), to)
	return *to
}

func SetLen[T any](sl []T, len int) []T {
	if len > cap(sl) {
		panic(errors.New("'len' must be less or equal to cap(sl)"))
	}
	return unsafe.Slice(unsafe.SliceData(sl), len)
}
