package main

import (
	"bufio"
	"fmt"
	"github.com/sonota88/vm2gol-v2-go/lib"
	"io"
	"os"
	"strconv"
	"strings"
)

// --------------------------------

type Token struct {
	kind string
	str  string
}

const NUM_TOKEN_MAX = 1024

var tokens [NUM_TOKEN_MAX]Token
var numTokens int
var pos = 0

// --------------------------------

func (self Token) kindEq(kind string) bool {
	return self.kind == kind
}

func (self Token) strEq(str string) bool {
	return self.str == str
}

func (self *Token) is(kind string, str string) bool {
	return self.kind == kind && self.strEq(str)
}

func addToken(line string, ti int) int {
	stripped := strings.Replace(line, "\n", "", 1)
	parts := strings.Split(stripped, ":")
	tokens[ti] = Token{
		kind: parts[0],
		str:  parts[1],
	}
	return 1
}

func readTokens() {
	var ti = 0
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if line != "" {
			ti += addToken(line, ti)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	numTokens = ti
}

// --------------------------------

func isEnd() bool {
	return numTokens <= pos
}

func peek(offset int) Token {
	return tokens[pos+offset]
}

func assertValue(kind string, str string) {
	t := peek(0)

	if t.kind != kind {
		puts_e(fmt.Sprintf("pos (%d)", pos))
		puts_e(fmt.Sprintf("expected (%s) (%s)", kind, str))
		puts_e(fmt.Sprintf("actual   (%s) (%s)", t.kind, t.str))
		panic("Unexpected kind")
	}

	if t.str != str {
		puts_e(fmt.Sprintf("pos (%d)", pos))
		puts_e(fmt.Sprintf("expected (%s) (%s)", kind, str))
		puts_e(fmt.Sprintf("actual   (%s) (%s)", t.kind, t.str))
		panic("Unexpected str")
	}
}

func consumeKw(str string) {
	assertValue("kw", str)
	pos++
}

func consumeSym(str string) {
	assertValue("sym", str)
	pos++
}

func newlist() *lib.NodeList {
	return lib.NodeList_new()
}

// --------------------------------

func parseArg() *lib.Node {
	// puts_fn("parseArg")

	t := peek(0)

	if t.kindEq("int") {
		pos++
		n, _ := strconv.Atoi(t.str)
		return lib.Node_newInt(n)
	} else if t.kindEq("ident") {
		pos++
		return lib.Node_newStr(t.str)
	} else {
		panic("not yet impl")
	}
}

func parseArgs() *lib.NodeList {
	puts_fn("parseArgs")

	args := newlist()

	if peek(0).str == ")" {
		return args
	} else {
		args.Add(parseArg())
	}

	for peek(0).str == "," {
		consumeSym(",")
		args.Add(parseArg())
	}

	return args
}

func parseFunc() *lib.NodeList {
	puts_fn("parseFunc")

	consumeKw("func")

	fnName := peek(0).str
	pos++

	consumeSym("(")
	args := parseArgs()
	consumeSym(")")

	consumeSym("{")

	stmts := newlist()
	for {
		t := peek(0)
		if t.str == "}" {
			break
		}

		if t.str == "var" {
			stmts.AddList(parseVar())
		} else {
			stmts.AddList(parseStmt())
		}
	}

	consumeSym("}")

	stmt := newlist()
	stmt.AddStr("func")
	stmt.AddStr(fnName)
	stmt.AddList(args)
	stmt.AddList(stmts)
	return stmt
}

func parseVarDeclare() *lib.NodeList {
	t := peek(0)
	pos++
	varName := t.str

	consumeSym(";")

	stmt := newlist()
	stmt.AddStr("var")
	stmt.AddStr(varName)
	return stmt
}

func parseVarInit() *lib.NodeList {
	t := peek(0)
	pos++
	varName := t.str

	consumeSym("=")

	expr := parseExpr()

	consumeSym(";")

	stmt := newlist()
	stmt.AddStr("var")
	stmt.AddStr(varName)
	stmt.Add(expr)
	return stmt
}

func parseVar() *lib.NodeList {
	puts_fn("parseVar")

	consumeKw("var")

	t := peek(1)

	if t.is("sym", ";") {
		return parseVarDeclare()
	} else {
		return parseVarInit()
	}
}

func parseExprRight() *lib.NodeList {
	// puts_fn("parseExprRight")

	t := peek(0)

	if t.is("sym", "+") ||
		t.is("sym", "*") ||
		t.is("sym", "==") ||
		t.is("sym", "!=") {
		// pass
	} else {
		return newlist()
	}

	op := t.str
	consumeSym(op)
	exprR := parseExpr()

	exprEls := newlist()
	exprEls.AddStr(op)
	exprEls.Add(exprR)

	return exprEls
}

func parseExpr() *lib.Node {
	puts_fn("parseExpr")
	var exprL *lib.Node

	tl := peek(0)

	if tl.kindEq("int") {
		pos++
		n, _ := strconv.Atoi(tl.str)
		exprL = lib.Node_newInt(n)
	} else if tl.kindEq("ident") {
		pos++
		s := tl.str
		exprL = lib.Node_newStr(s)
	} else if tl.kindEq("sym") {
		consumeSym("(")
		exprL = parseExpr()
		consumeSym(")")
	} else {
		panic("not_yet_impl")
	}

	opRight := parseExprRight()
	if opRight.Len() == 0 {
		return exprL
	}

	op := opRight.Get(0)
	exprR := opRight.Get(1)

	exprEls := newlist()
	exprEls.Add(op)
	exprEls.Add(exprL)
	exprEls.Add(exprR)
	return lib.Node_newList(exprEls)
}

func parseSet() *lib.NodeList {
	puts_fn("parseSet")

	consumeKw("set")

	t := peek(0)
	pos++
	varName := t.str

	consumeSym("=")

	expr := parseExpr()

	consumeSym(";")

	ret := newlist()
	ret.AddStr("set")
	ret.AddStr(varName)
	ret.Add(expr)
	return ret
}

func parseFuncall() *lib.NodeList {
	// puts_fn("parseFuncall")

	t := peek(0)
	pos++
	fnName := t.str

	consumeSym("(")
	args := parseArgs()
	consumeSym(")")

	list := newlist()
	list.AddStr(fnName)
	list.AddListAll(args)

	return list
}

func parseCall() *lib.NodeList {
	puts_fn("parseCall")

	consumeKw("call")

	funcall := parseFuncall()

	consumeSym(";")

	ret := newlist()
	ret.AddStr("call")
	ret.AddListAll(funcall)

	return ret
}

func parseCallSet() *lib.NodeList {
	puts_fn("parseCallSet")

	consumeKw("call_set")

	t := peek(0)
	pos++
	varName := t.str

	consumeSym("=")

	funcall := parseFuncall()

	consumeSym(";")

	stmt := newlist()
	stmt.AddStr("call_set")
	stmt.AddStr(varName)
	stmt.AddList(funcall)
	return stmt
}

func parseReturn() *lib.NodeList {
	puts_fn("parseReturn")

	consumeKw("return")

	expr := parseExpr()

	consumeSym(";")

	stmt := newlist()
	stmt.AddStr("return")
	stmt.Add(expr)
	return stmt
}

func parseWhile() *lib.NodeList {
	puts_fn("parseWhile")

	consumeKw("while")

	consumeSym("(")
	expr := parseExpr()
	consumeSym(")")

	consumeSym("{")
	stmts := parseStmts()
	consumeSym("}")

	stmt := newlist()
	stmt.AddStr("while")
	stmt.Add(expr)
	stmt.AddList(stmts)
	return stmt
}

func parseWhenClause() *lib.NodeList {
	// puts_fn("parseWhenClause")

	t := peek(0)
	if t.is("sym", "}") {
		return lib.NodeList_empty()
	}

	consumeSym("(")
	expr := parseExpr()
	consumeSym(")")

	consumeSym("{")
	stmts := parseStmts()
	consumeSym("}")

	list := newlist()
	list.Add(expr)
	for i := 0; i < stmts.Len(); i++ {
		stmt := stmts.Get(i).List
		list.AddList(stmt)
	}

	return list
}

func parseCase() *lib.NodeList {
	puts_fn("parseCase")

	consumeKw("case")

	consumeSym("{")

	whenClauses := newlist()

	for {
		whenClause := parseWhenClause()
		if whenClause.Len() == 0 {
			break
		}
		whenClauses.AddList(whenClause)
	}

	consumeSym("}")

	stmt := newlist()
	stmt.AddStr("case")

	for i := 0; i < whenClauses.Len(); i++ {
		whenClause := whenClauses.Get(i).List
		stmt.AddList(whenClause)
	}

	return stmt
}

func parseVmComment() *lib.NodeList {
	puts_fn("parseVmComment")

	consumeKw("_cmt")
	consumeSym("(")

	t := peek(0)
	pos++
	cmt := t.str

	consumeSym(")")
	consumeSym(";")

	ret := newlist()
	ret.AddStr("_cmt")
	ret.AddStr(cmt)
	return ret
}

func parseDebug() *lib.NodeList {
	puts_fn("parseDebug")

	consumeKw("_debug")
	consumeSym("(")
	consumeSym(")")
	consumeSym(";")

	ret := newlist()
	ret.AddStr("_cmt")
	return ret
}

func parseStmt() *lib.NodeList {
	puts_fn("parseStmt")

	t := peek(0)

	if t.is("sym", "}") {
		return nil
	}

	if t.strEq("set") {
		return parseSet()
	} else if t.strEq("call") {
		return parseCall()
	} else if t.strEq("call_set") {
		return parseCallSet()
	} else if t.strEq("return") {
		return parseReturn()
	} else if t.strEq("while") {
		return parseWhile()
	} else if t.strEq("case") {
		return parseCase()
	} else if t.strEq("_cmt") {
		return parseVmComment()
	} else if t.strEq("_debug") {
		return parseDebug()
	} else {
		puts_kv_e("pos", pos)
		puts_kv_e("t", t)
		panic("Unexpected token")
	}
}

func parseStmts() *lib.NodeList {
	stmts := newlist()

	for !isEnd() {
		stmt := parseStmt()
		if stmt == nil {
			break
		}
		stmts.AddList(stmt)
	}

	return stmts
}

func parseTopStmt() *lib.NodeList {
	t := tokens[pos]

	if t.str == "func" {
		return parseFunc()
	} else {
		panic(
			fmt.Sprintf(
				"Unexpected token: pos(%d) kind(%s) str(%s)",
				pos, t.kind, t.str))
	}
}

func parseTopStmts() *lib.NodeList {
	tree := newlist()
	tree.AddStr("top_stmts")

	for {
		if isEnd() {
			break
		}

		tree.AddList(parseTopStmt())
	}

	return tree
}

func Parse() {
	readTokens()
	puts_e(fmt.Sprintf("numTokens(%d)", numTokens))

	topStmts := parseTopStmts()
	lib.PrintAsJson(topStmts)
}
