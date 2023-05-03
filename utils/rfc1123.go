package utils

import (
	"math/rand"
	"strings"
)

func randomString(length int, str string) string {
	var letters []rune
	var ret []rune

	letters = []rune(str)
	ret = make([]rune, length)
	for i := range ret {
		ret[i] = letters[rand.Intn(len(letters))]
	}
	return (string(ret))
}

func RFC1123Gen(n uint) string {
	var letters string
	var str string

	letters = "abcdefghijklmnopqrstuvwxyz0123456789-"
	if n >= 1 {
		str = str + randomString(1, letters[:26])
	}
	if n >= 3 {
		str = str + randomString(int(n-2), letters[:37])
	}
	if n >= 2 {
		str = str + randomString(1, letters[:36])
	}
	return (str)
}

func RFC1123Check(str string) bool {
	var flag bool
	var letters string

	flag = true
	letters = "abcdefghijklmnopqrstuvwxyz0123456789-"
	if !strings.Contains(letters[:26], str[:1]) {
		flag = false
	}
	for n := range str[1 : len(str)-1] {
		if !strings.Contains(letters[:37], str[1 : len(str)-1][n:n+1]) {
			flag = false
		}
	}
	if !strings.Contains(letters[:36], str[len(str)-1:]) {
		flag = false
	}
	return (flag)
}
