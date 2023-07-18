package parse

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/jamesTait-jt/go-performance/internal/unsafe"
)

type person struct {
	firstName string
	lastName string
	age int
	height float64
}

func Person(s string) (person, error) {
	split := strings.Split(s, " ")

	age, err := strconv.Atoi(split[2])
	if err != nil {
		return person{}, err
	}

	height, err := strconv.ParseFloat(split[3], 64)
	if err != nil {
		return person{}, err
	}

	p := person{
		firstName: split[0],
		lastName: split[1],
		age: age,
		height: height,
	}

	return p, nil
}

func PersonEfficient(s string) (person, error) {
	first, rest := ParseUntil(s, ' ')
	last, rest := ParseUntil(rest, ' ')
	ageStr, rest := ParseUntil(rest, ' ')
	heightStr, _ := ParseUntil(rest, ' ')

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return person{}, err
	}

	height, err := strconv.ParseFloat(heightStr, 64)
	if err != nil {
		return person{}, err
	}

	p := person{
		firstName: first,
		lastName: last,
		age: age,
		height: height,
	}

	return p, nil
}

func PersonBytes(bs []byte) (person, error) {
	first, rest := ParseUntilBytes(bs, '\x1f')
	last, rest := ParseUntilBytes(rest, '\x1f')
	ageBytes, rest := ParseUntilBytes(rest, '\x1f')
	heightBytes, _ := ParseUntilBytes(rest, '\x1f')

	ageStr := string(ageBytes)
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return person{}, err
	}

	heightStr := string(heightBytes)
	height, err := strconv.ParseFloat(heightStr, 64)
	if err != nil {
		return person{}, err
	}

	p := person{
		firstName: string(first),
		lastName: string(last),
		age: age,
		height: height,
	}

	return p, nil
}

func PersonBytesUnsafe(bs []byte) (person, bool) {
	first, rest := ParseUntilBytes(bs, '\x1f')
	last, rest := ParseUntilBytes(rest, '\x1f')
	ageBytes, rest := ParseUntilBytes(rest, '\x1f')
	heightBytes, _ := ParseUntilBytes(rest, '\x1f')

	age, ok := unsafe.ToInt(ageBytes)
	if !ok {
		return person{}, false
	}

	height, ok := unsafe.ToFloat64(heightBytes)
	if !ok {
		return person{}, false
	}

	p := person{
		firstName: unsafe.ToString(first),
		lastName: unsafe.ToString(last),
		age: age,
		height: height,
	}

	return p, true
}


func ParseUntil(s string, sep rune) (string, string) {
	if len(s) == 0 {
		return "", ""
	}

	indexOfNext := strings.Index(s, " ")
	if indexOfNext == -1 {
		return s, ""
	}

	return s[:indexOfNext], s[indexOfNext + 1:]
}

func ParseUntilBytes(bs []byte, sep byte) ([]byte, []byte) {
	if len(bs) == 0 {
		return nil, nil
	}

	// \x1f is the ASCII control character for a unit separator
	indexOfNext := bytes.IndexByte(bs, '\x1f')
	if indexOfNext == -1 {
		return bs, nil
	}

	return bs[:indexOfNext], bs[indexOfNext + 1:]
}