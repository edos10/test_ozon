package main

import (
	"fmt"
)

var numToLetter = make(map[int]byte)
var letterToNum = make(map[byte]int)
var arraySymb = make([]byte, 0)

var startSmallLts = 97
var endSmallLts = 122
var startBigLts = 65
var endBigLts = 90
var startNums = 48
var downcase = 95

const sizeUrl, sizeAlphabet = 10, 63

/*
	 Порядок следования всех символов в map:
	a-z = 0-25
	A-Z = 26-51
	0-9 = 52-61
	_ = 62
*/

func makeMaps() {
	startSym := byte('a')
	for i := 0; i < 26; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	startSym = byte('A')
	for i := 26; i < 52; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	startSym = byte('0')
	for i := 52; i < 62; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	numToLetter[62] = '_'

	startSym = byte('a')
	numForMap := 0
	for startSym <= 'z' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	startSym = byte('A')
	for startSym <= 'Z' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	startSym = byte('0')
	for startSym <= '9' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	numToLetter['_'] = 63
}

func NextUrlString(current string) string {
	newUrl := make([]int, 0)
	for _, elem := range []byte(current) {
		newUrl = append(newUrl, letterToNum[elem])
	}
	n := sizeUrl
	k := sizeAlphabet
	i := n - 1
	for i >= 0 {
		if newUrl[i] < k-1 {
			break
		}
		i--
	}
	if i == -1 {
		return ""
	}
	newUrl[i] += 1
	for j := i + 1; j < n; j++ {
		newUrl[j] = 0
	}
	res := make([]byte, 0)
	for _, elem := range newUrl {
		res = append(res, numToLetter[elem])
	}
	return string(res)
}

func makeArraySym() {
	startSym := byte('a')
	for startSym <= 'z' {
		arraySymb = append(arraySymb, startSym)
		startSym++
	}

	startSym = byte('A')
	for startSym <= 'Z' {
		arraySymb = append(arraySymb, startSym)
		startSym++
	}

	startSym = byte('0')
	for startSym <= '9' {
		arraySymb = append(arraySymb, startSym)
		startSym++
	}

	arraySymb = append(arraySymb, '_')
}

func getRandomUrl() string {
	res := make([]byte, 0)
	fmt.Println(arraySymb)
	for i := 0; i < 10; i++ {
		res = append(res, arraySymb[2])
	}
	return string(res)
}
