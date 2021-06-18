package main

import (
	"fmt"
	"github.com/sonota88/vm2gol-v2-go/lib"
	"regexp"
)

func putsToken(kind string, str string) {
	fmt.Printf("%s:%s\n", kind, str)
}

func isKw(str string) bool {
	return str == "func" ||
		str == "set" ||
		str == "var" ||
		str == "call_set" ||
		str == "call" ||
		str == "return" ||
		str == "case" ||
		str == "while" ||
		str == "_cmt"
}

func tokenize(src string) {
	pos := 0

	spaceRe := regexp.MustCompile(`^([ \n]+)`)
	cmtRe := regexp.MustCompile(`^(//.*)`)
	strRe := regexp.MustCompile(`^"(.*)"`)
	symbolRe := regexp.MustCompile(`^(==|!=|[(){}=;+*,])`)
	intRe := regexp.MustCompile(`^(-?[0-9]+)`)
	identRe := regexp.MustCompile(`^([a-z_][a-z0-9_\[\]]*)`)

	for pos < len(src) {
		rest := src[pos:]

		if spaceRe.MatchString(rest) {
			m := spaceRe.FindStringSubmatch(rest)
			temp := m[1]
			pos += len(temp)

		} else if cmtRe.MatchString(rest) {
			m := cmtRe.FindStringSubmatch(rest)
			temp := m[1]
			pos += len(temp)

		} else if strRe.MatchString(rest) {
			m := strRe.FindStringSubmatch(rest)
			temp := m[1]
			putsToken("str", temp)
			pos += len(temp) + 2

		} else if intRe.MatchString(rest) {
			m := intRe.FindStringSubmatch(rest)
			temp := m[1]
			putsToken("int", temp)
			pos += len(temp)

		} else if symbolRe.MatchString(rest) {
			m := symbolRe.FindStringSubmatch(rest)
			temp := m[1]
			putsToken("sym", temp)
			pos += len(temp)

		} else if identRe.MatchString(rest) {
			m := identRe.FindStringSubmatch(rest)
			temp := m[1]
			if isKw(temp) {
				putsToken("kw", temp)
			} else {
				putsToken("ident", temp)
			}
			pos += len(temp)

		} else {
			panic("Unexpected pattern (" + rest + ")")
		}
	}
}

func Tokenize() {
	src := lib.ReadStdinAll()
	tokenize(src)
}
