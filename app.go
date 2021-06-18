package main

import (
	"fmt"
	"github.com/sonota88/vm2gol-v2-go/lib"
	"github.com/sonota88/vm2gol-v2-go/test"
	"os"
)

func puts_e(arg interface{}) {
	lib.Puts_e(arg)
}

func puts_kv_e(k string, v interface{}) {
	lib.Puts_e(fmt.Sprintf("%s (%s)", k, v))
}

func puts_fn(msg string) {
	puts_e("    |-->> " + msg + "()")
}

func main() {
	subcmd := os.Args[1]
	if subcmd == "test_json" {
		test.TestJson()
	} else if subcmd == "tokenize" {
		Tokenize()
	} else if subcmd == "parse" {
		Parse()
	} else if subcmd == "codegen" {
		Codegen()
	} else {
		panic(fmt.Sprintf("invalid sub command (%s)", subcmd))
	}
}
