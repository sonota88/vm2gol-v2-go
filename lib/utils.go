package lib

import (
	"bufio"
	// "fmt"
	"io"
    "io/ioutil"
	"os"
)

func Puts_e(arg interface{}) {
	// fmt.Fprint(os.Stderr, arg)
	// fmt.Fprint(os.Stderr, "\n")
}

func ReadStdinAll_v1() string {
	var s = ""
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		s += line
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return s
}

func ReadStdinAll() string {
    bytes, _ := ioutil.ReadAll(os.Stdin)
    return string(bytes)
}
