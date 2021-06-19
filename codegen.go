package main

import (
	"fmt"
	"github.com/sonota88/vm2gol-v2-go/lib"
	"strings"
)

var gLabelId = 0

// --------------------------------

func getLabelId() int {
	gLabelId++
	return gLabelId
}

func head(list *lib.NodeList) *lib.Node {
	return list.Get(0)
}

func rest(list *lib.NodeList) *lib.NodeList {
	newList := lib.NodeList_new()
	for i := 1; i < list.Len(); i++ {
		newList.Add(list.Get(i))
	}
	return newList
}

// --------------------------------

func toFnArgRef(names *lib.Names, name string) string {
	i := names.IndexOf(name)
	if i == -1 {
		panic("fn arg not found")
	}
	return fmt.Sprintf("[bp:%d]", i+2)
}

func toLvarRef(names *lib.Names, name string) string {
	i := names.IndexOf(name)
	if i == -1 {
		panic("lvar not found")
	}
	return fmt.Sprintf("[bp:%d]", -(i + 1))
}

// --------------------------------

func codegenVar(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	fmt.Printf("  sub_sp 1\n")

	if stmtRest.Len() == 2 {
		codegenSet(fnArgNames, lvarNames, stmtRest)
	}
}

func codegenExprAdd() {
	fmt.Printf("  pop reg_b\n")
	fmt.Printf("  pop reg_a\n")
	fmt.Printf("  add_ab\n")
}

func codegenExprMult() {
	fmt.Printf("  pop reg_b\n")
	fmt.Printf("  pop reg_a\n")
	fmt.Printf("  mult_ab\n")
}

func codegenExprEq() {
	labelId := getLabelId()

	thenLabel := fmt.Sprintf("then_%d", labelId)
	endLabel := fmt.Sprintf("end_eq_%d", labelId)

	fmt.Printf("  pop reg_b\n")
	fmt.Printf("  pop reg_a\n")

	fmt.Printf("  compare\n")
	fmt.Printf("  jump_eq %s\n", thenLabel)

	fmt.Printf("  cp 0 reg_a\n")
	fmt.Printf("  jump %s\n", endLabel)

	fmt.Printf("label %s\n", thenLabel)
	fmt.Printf("  cp 1 reg_a\n")
	fmt.Printf("label %s\n", endLabel)
}

func codegenExprNeq() {
	labelId := getLabelId()

	thenLabel := fmt.Sprintf("then_%d", labelId)
	endLabel := fmt.Sprintf("end_neq_%d", labelId)

	fmt.Printf("  pop reg_b\n")
	fmt.Printf("  pop reg_a\n")

	fmt.Printf("  compare\n")
	fmt.Printf("  jump_eq %s\n", thenLabel)

	fmt.Printf("  cp 1 reg_a\n")
	fmt.Printf("  jump %s\n", endLabel)

	fmt.Printf("label %s\n", thenLabel)
	fmt.Printf("  cp 0 reg_a\n")
	fmt.Printf("label %s\n", endLabel)
}

func _codegenExprBinop(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	expr *lib.NodeList,
) {
	puts_fn("_codegenExprBinop")

	op := head(expr).Strval
	args := rest(expr)

	termL := args.Get(0)
	termR := args.Get(1)

	codegenExpr(fnArgNames, lvarNames, termL)
	fmt.Printf("  push reg_a\n")
	codegenExpr(fnArgNames, lvarNames, termR)
	fmt.Printf("  push reg_a\n")

	if op == "+" {
		codegenExprAdd()
	} else if op == "*" {
		codegenExprMult()
	} else if op == "eq" {
		codegenExprEq()
	} else if op == "neq" {
		codegenExprNeq()
	} else {
		puts_kv_e("op", op)
		panic("not_yet_impl")
	}
}

func codegenExpr(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	expr *lib.Node,
) {
	if expr.KindEq("int") {
		fmt.Printf("  cp %d reg_a\n", expr.Intval)
	} else if expr.KindEq("str") {
		str := expr.Strval
		if 0 <= lvarNames.IndexOf(str) {
			cpSrc := toLvarRef(lvarNames, str)
			fmt.Printf("  cp %s reg_a\n", cpSrc)
		} else if 0 <= fnArgNames.IndexOf(str) {
			cpSrc := toFnArgRef(fnArgNames, str)
			fmt.Printf("  cp %s reg_a\n", cpSrc)
		} else {
			panic("not_yet_impl")
		}
	} else if expr.KindEq("list") {
		_codegenExprBinop(fnArgNames, lvarNames, expr.List)
	} else {
		panic("not_yet_impl")
	}
}

func _codegenFuncall(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	fnName string,
	fnArgs *lib.NodeList,
) {
	for i := fnArgs.Len() - 1; i >= 0; i-- {
		fnArg := fnArgs.Get(i)
		codegenExpr(fnArgNames, lvarNames, fnArg)
		fmt.Printf("  push reg_a\n")
	}

	codegenVmComment("call  " + fnName)
	fmt.Printf("  call %s\n", fnName)

	fmt.Printf("  add_sp %d\n", fnArgs.Len())
}

func codegenCall(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	fnName := head(stmtRest).Strval
	fnArgs := rest(stmtRest)
	_codegenFuncall(fnArgNames, lvarNames, fnName, fnArgs)
}

func codegenCallSet(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	puts_fn("codegenCallSet")

	lvarName := stmtRest.Get(0).Strval
	fnTemp := stmtRest.Get(1).List

	fnName := head(fnTemp).Strval
	fnArgs := rest(fnTemp)

	_codegenFuncall(fnArgNames, lvarNames, fnName, fnArgs)

	dest := toLvarRef(lvarNames, lvarName)
	fmt.Printf("  cp reg_a %s\n", dest)
}

func codegenSet(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	rest *lib.NodeList,
) {
	puts_fn("codegenSet")
	dest := rest.Get(0)
	expr := rest.Get(1)

	codegenExpr(fnArgNames, lvarNames, expr)

	if dest.KindEq("str") {
		if 0 <= lvarNames.IndexOf(dest.Strval) {
			cpDest := toLvarRef(lvarNames, dest.Strval)
			fmt.Printf("  cp reg_a %s\n", cpDest)
		} else {
			panic("not_yet_impl")
		}
	} else {
		panic("not_yet_impl")
	}
}

func codegenReturn(
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	retval := head(stmtRest)
	codegenExpr(lib.Names_empty(), lvarNames, retval)
}

func codegenVmComment(cmt string) {
	fmt.Printf("  _cmt %s\n", strings.Replace(cmt, " ", "~", -1))
}

func codegenWhile(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	puts_fn("codegenWhile")

	condExpr := stmtRest.Get(0)
	body := stmtRest.Get(1).List

	labelId := getLabelId()
	labelBegin := fmt.Sprintf("while_%d", labelId)
	labelEnd := fmt.Sprintf("end_while_%d", labelId)
	labelTrue := fmt.Sprintf("true_%d", labelId)

	fmt.Printf("\n")

	fmt.Printf("label %s\n", labelBegin)

	codegenExpr(fnArgNames, lvarNames, condExpr)

	fmt.Printf("  cp 1 reg_b\n")
	fmt.Printf("  compare\n")

	fmt.Printf("  jump_eq %s\n", labelTrue)
	fmt.Printf("  jump %s\n", labelEnd)
	fmt.Printf("label %s\n", labelTrue)

	codegenStmts(fnArgNames, lvarNames, body)

	fmt.Printf("  jump %s\n", labelBegin)

	fmt.Printf("label %s\n", labelEnd)
	fmt.Printf("\n")
}

func codegenCase(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	whenClauses *lib.NodeList,
) {
	puts_fn("codegenCase")

	labelId := getLabelId()
	whenIdx := -1

	labelEnd := fmt.Sprintf("end_case_%d", labelId)
	labelWhenHead := fmt.Sprintf("when_%d", labelId)
	labelEndWhenHead := fmt.Sprintf("end_when_%d", labelId)

	fmt.Printf("\n")
	fmt.Printf("  # -->> case_%d\n", labelId)

	for i := 0; i < whenClauses.Len(); i++ {
		whenClause := whenClauses.Get(i).List
		whenIdx++

		cond := head(whenClause)
		_rest := rest(whenClause)

		fmt.Printf("  # when_%d_%d\n", labelId, whenIdx)

		fmt.Printf("  # -->> expr\n")
		codegenExpr(fnArgNames, lvarNames, cond)
		fmt.Printf("  # <<-- expr\n")

		fmt.Printf("  cp 1 reg_b\n")

		fmt.Printf("  compare\n")
		fmt.Printf("  jump_eq %s_%d\n", labelWhenHead, whenIdx)
		fmt.Printf("  jump %s_%d\n", labelEndWhenHead, whenIdx)

		fmt.Printf("label %s_%d\n", labelWhenHead, whenIdx)

		codegenStmts(fnArgNames, lvarNames, _rest)

		fmt.Printf("  jump %s\n", labelEnd)
		fmt.Printf("label %s_%d\n", labelEndWhenHead, whenIdx)
	}

	fmt.Printf("label end_case_%d\n", labelId)
	fmt.Printf("  # <<-- case_%d\n", labelId)
	fmt.Printf("\n")
}

func codegenStmt(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmt *lib.NodeList,
) {
	puts_fn("codegenStmt")

	stmtHead := head(stmt).Strval
	stmtRest := rest(stmt)

	if stmtHead == "set" {
		codegenSet(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "call" {
		codegenCall(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "call_set" {
		codegenCallSet(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "return" {
		codegenReturn(lvarNames, stmtRest)
	} else if stmtHead == "while" {
		codegenWhile(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "case" {
		codegenCase(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "_cmt" {
		codegenVmComment(stmtRest.Get(0).Strval)
	} else {
		puts_kv_e("stmtHead", stmtHead)
		panic("Unsupported statement")
	}
}

func codegenStmts(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmts *lib.NodeList,
) {
	for i := 0; i < stmts.Len(); i++ {
		stmt := stmts.Get(i).List
		codegenStmt(fnArgNames, lvarNames, stmt)
	}
}

func codegenFuncDef(topStmt *lib.NodeList) {
	fnName := topStmt.Get(0).Strval
	fnArgVals := topStmt.Get(1).List
	body := topStmt.Get(2).List

	fnArgNames := lib.Names_new()
	for i := 0; i < fnArgVals.Len(); i++ {
		fnArgNames.Add(fnArgVals.Get(i).Strval)
	}

	lvarNames := lib.Names_new()

	fmt.Println("")
	fmt.Printf("label %s\n", fnName)
	fmt.Println(`  push bp`)
	fmt.Println(`  cp sp bp`)

	fmt.Println("")
	fmt.Println("  # 関数の処理本体")

	for i := 0; i < body.Len(); i++ {
		stmt := body.Get(i).List
		stmtHead := head(stmt).Strval
		stmtRest := rest(stmt)

		if stmtHead == "var" {
			varName := stmtRest.Get(0).Strval
			lvarNames.Add(varName)
			codegenVar(fnArgNames, lvarNames, stmtRest)
		} else {
			codegenStmt(fnArgNames, lvarNames, stmt)
		}
	}

	fmt.Println("")

	fmt.Println(`  cp bp sp`)
	fmt.Println(`  pop bp`)
	fmt.Println(`  ret`)
}

func codegenTopStmts(topStmts *lib.NodeList) {
	topStmtsRest := rest(topStmts)
	for i := 0; i < topStmtsRest.Len(); i++ {
		topStmt := topStmtsRest.Get(i).List
		stmtRest := rest(topStmt)
		codegenFuncDef(stmtRest)
	}
}

func codegenBuiltinSetVram() {
	fmt.Println("")
	fmt.Println("label set_vram")
	fmt.Println("  push bp")
	fmt.Println("  cp sp bp")

	fmt.Println("  set_vram [bp:2] [bp:3]") // vram_addr value

	fmt.Println("  cp bp sp")
	fmt.Println("  pop bp")
	fmt.Println("  ret")
}

func codegenBuiltinGetVram() {
	fmt.Println("")
	fmt.Println("label get_vram")
	fmt.Println("  push bp")
	fmt.Println("  cp sp bp")

	fmt.Println("  get_vram [bp:2] reg_a") // vram_addr dest

	fmt.Println("  cp bp sp")
	fmt.Println("  pop bp")
	fmt.Println("  ret")
}

func Codegen() {
	json := lib.ReadStdinAll()
	tree := lib.ParseJson(json)

	fmt.Println(`  call main`)
	fmt.Println(`  exit`)

	codegenTopStmts(tree)

	fmt.Println("#>builtins")
	codegenBuiltinSetVram()
	codegenBuiltinGetVram()
	fmt.Println("#<builtins")
}
