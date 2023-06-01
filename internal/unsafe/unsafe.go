package unsafe

import (
	"fmt"
	"strconv"
	"unsafe"
)

func ToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func ToInt(bytes []byte) (int, bool) {
	s := ToString(bytes)

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}

	return n, true
}

func ToFloat64(bytes []byte) (float64, bool) {
	s := ToString(bytes)

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}

	return f, true
}

func example() {
	bs := []byte{'1', '1', '1'}
	s := ToString(bs)

	fmt.Println(s)

	bs[0] = '0'

	fmt.Println(s)
}
