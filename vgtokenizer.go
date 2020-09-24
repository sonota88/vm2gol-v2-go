package main

import (
	"./lib"
	"fmt"
	"regexp"
)

func putsToken(kind string, str string) {
	fmt.Printf("%s:%s\n", kind, str)
}

func tokenize(src string) {
	pos := 0

	spaceRe := regexp.MustCompile(`^([ \n]+)`)
	cmtRe := regexp.MustCompile(`^(//.*)`)
	strRe := regexp.MustCompile(`^"(.*)"`)
	kwRe := regexp.MustCompile(`^(def|end|set|var|call_set|call|return|case|while|_cmt)[^a-z]`)
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

		} else if kwRe.MatchString(rest) {
			m := kwRe.FindStringSubmatch(rest)
			temp := m[1]
			putsToken("kw", temp)
			pos += len(temp)

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
			putsToken("ident", temp)
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
