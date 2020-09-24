package lib

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func isSpace(c byte) bool {
	return c == '\n' || c == ' '
}

func parseJson(json string) (*NodeList, int) {
	var list NodeList
	var pos = 1

	var intRe = regexp.MustCompile(`^(-?[0-9]+)`)
	var strRe = regexp.MustCompile(`^"(.*?)"`)

	for pos < len(json) {
		rest := json[pos:]

		if rest[0] == ']' {
			pos++
			break

		} else if isSpace(rest[0]) || rest[0] == ',' {
			pos++

		} else if rest[0] == '"' {
			group := strRe.FindStringSubmatch(rest)
			list.Add(Node_newStr(group[1]))
			pos += len(group[1]) + 2

		} else if intRe.MatchString(rest) {
			group := intRe.FindStringSubmatch(rest)
			n, _ := strconv.Atoi(group[1])
			list.Add(Node_newInt(n))
			pos += len(group[1])

		} else if rest[0] == '[' {
			tempList, size := parseJson(rest)
			list.Add(Node_newList(tempList))
			pos += size

		} else {
			if len(rest) < 10 {
				fmt.Fprintf(os.Stderr, "pos (%d) rest (%s)\n", pos, rest)
			} else {
				fmt.Fprintf(os.Stderr, "pos (%d) rest[0:10] (%s)\n", pos, rest[0:10])
			}
			panic("Unexpected pattern")
		}
	}

	return &list, pos
}

func ParseJson(json string) *NodeList {
	list, _ := parseJson(json)
	return list
}

func printIndent(lv int) {
	for i := 0; i < lv; i++ {
		fmt.Print("  ")
	}
}

func printNode(node *Node, lv int) {
	if node.KindEq("int") {
		n := node.Intval
		printIndent(lv + 1)
		fmt.Print(n)
	} else if node.KindEq("str") {
		s := node.Strval
		printIndent(lv + 1)
		fmt.Print("\"")
		fmt.Print(s)
		fmt.Print("\"")
	} else if node.KindEq("list") {
		printNodeList(node.List, lv+1)
	} else {
		panic("must not happen")
	}
}

func printNodeList(nodeList *NodeList, lv int) {
	printIndent(lv)
	fmt.Print("[")
	fmt.Print("\n")

	for i := 0; i < nodeList.Len(); i++ {
		if i >= 1 {
			fmt.Print(",")
			fmt.Print("\n")
		}
		node := nodeList.Get(i)
		printNode(node, lv)
	}
	fmt.Print("\n")
	printIndent(lv)
	fmt.Print("]")
}

func PrintAsJson(nodeList *NodeList) {
	printNodeList(nodeList, 0)
}
