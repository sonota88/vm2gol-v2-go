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

func toFnArgDisp(names *lib.Names, name string) int {
	i := names.IndexOf(name)
	if i == -1 {
		panic("fn arg not found")
	}
	return i + 2
}

func toLvarDisp(names *lib.Names, name string) int {
	i := names.IndexOf(name)
	if i == -1 {
		panic("lvar not found")
	}
	return -(i + 1)
}

// --------------------------------

func asmPrologue() {
	fmt.Println("  push bp")
	fmt.Println("  cp sp bp")
}

func asmEpilogue() {
	fmt.Println("  cp bp sp")
	fmt.Println("  pop bp")
}

// --------------------------------

func genExprAdd() {
	fmt.Printf("  pop reg_b\n")
	fmt.Printf("  pop reg_a\n")
	fmt.Printf("  add_ab\n")
}

func genExprMult() {
	fmt.Printf("  pop reg_b\n")
	fmt.Printf("  pop reg_a\n")
	fmt.Printf("  mult_ab\n")
}

func genExprEq() {
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

func genExprNeq() {
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

func _genExprBinop(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	expr *lib.NodeList,
) {
	puts_fn("_genExprBinop")

	op := head(expr).Strval
	args := rest(expr)

	termL := args.Get(0)
	termR := args.Get(1)

	genExpr(fnArgNames, lvarNames, termL)
	fmt.Printf("  push reg_a\n")
	genExpr(fnArgNames, lvarNames, termR)
	fmt.Printf("  push reg_a\n")

	if op == "+" {
		genExprAdd()
	} else if op == "*" {
		genExprMult()
	} else if op == "==" {
		genExprEq()
	} else if op == "!=" {
		genExprNeq()
	} else {
		puts_kv_e("op", op)
		panic("not_yet_impl")
	}
}

func genExpr(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	expr *lib.Node,
) {
	if expr.KindEq("int") {
		fmt.Printf("  cp %d reg_a\n", expr.Intval)
	} else if expr.KindEq("str") {
		str := expr.Strval
		if 0 <= lvarNames.IndexOf(str) {
			disp := toLvarDisp(lvarNames, str)
			fmt.Printf("  cp [bp:%d] reg_a\n", disp)
		} else if 0 <= fnArgNames.IndexOf(str) {
			disp := toFnArgDisp(fnArgNames, str)
			fmt.Printf("  cp [bp:%d] reg_a\n", disp)
		} else {
			panic("not_yet_impl")
		}
	} else if expr.KindEq("list") {
		_genExprBinop(fnArgNames, lvarNames, expr.List)
	} else {
		panic("not_yet_impl")
	}
}

func _genFuncall(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	funcall *lib.NodeList,
) {
	fnName := head(funcall).Strval
	fnArgs := rest(funcall)

	for i := fnArgs.Len() - 1; i >= 0; i-- {
		fnArg := fnArgs.Get(i)
		genExpr(fnArgNames, lvarNames, fnArg)
		fmt.Printf("  push reg_a\n")
	}

	genVmComment("call  " + fnName)
	fmt.Printf("  call %s\n", fnName)

	fmt.Printf("  add_sp %d\n", fnArgs.Len())
}

func genCall(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmt *lib.NodeList,
) {
	funcall := rest(stmt)
	_genFuncall(fnArgNames, lvarNames, funcall)
}

func genCallSet(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	puts_fn("genCallSet")

	lvarName := stmtRest.Get(0).Strval
	funcall := stmtRest.Get(1).List

	_genFuncall(fnArgNames, lvarNames, funcall)

	disp := toLvarDisp(lvarNames, lvarName)
	fmt.Printf("  cp reg_a [bp:%d]\n", disp)
}

func _genSet(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	varName string,
	expr *lib.Node,
) {
	puts_fn("_genSet")

	genExpr(fnArgNames, lvarNames, expr)

	if 0 <= lvarNames.IndexOf(varName) {
		disp := toLvarDisp(lvarNames, varName)
		fmt.Printf("  cp reg_a [bp:%d]\n", disp)
	} else {
		panic("not_yet_impl")
	}
}

func genSet(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmt *lib.NodeList,
) {
	puts_fn("genSet")
	dest := stmt.Get(1)
	expr := stmt.Get(2)

	if dest.KindEq("str") {
		_genSet(fnArgNames, lvarNames, dest.Strval, expr)
	} else {
		panic("not_yet_impl")
	}
}

func genReturn(
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	retval := head(stmtRest)
	genExpr(lib.Names_empty(), lvarNames, retval)
}

func genVmComment(cmt string) {
	fmt.Printf("  _cmt %s\n", strings.Replace(cmt, " ", "~", -1))
}

func genDebug() {
	fmt.Println("  _debug")
}

func genWhile(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmtRest *lib.NodeList,
) {
	puts_fn("genWhile")

	condExpr := stmtRest.Get(0)
	body := stmtRest.Get(1).List

	labelId := getLabelId()
	labelBegin := fmt.Sprintf("while_%d", labelId)
	labelEnd := fmt.Sprintf("end_while_%d", labelId)
	labelTrue := fmt.Sprintf("true_%d", labelId)

	fmt.Printf("\n")

	fmt.Printf("label %s\n", labelBegin)

	genExpr(fnArgNames, lvarNames, condExpr)

	fmt.Printf("  cp 1 reg_b\n")
	fmt.Printf("  compare\n")

	fmt.Printf("  jump_eq %s\n", labelTrue)
	fmt.Printf("  jump %s\n", labelEnd)
	fmt.Printf("label %s\n", labelTrue)

	genStmts(fnArgNames, lvarNames, body)

	fmt.Printf("  jump %s\n", labelBegin)

	fmt.Printf("label %s\n", labelEnd)
	fmt.Printf("\n")
}

func genCase(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	whenClauses *lib.NodeList,
) {
	puts_fn("genCase")

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
		genExpr(fnArgNames, lvarNames, cond)
		fmt.Printf("  # <<-- expr\n")

		fmt.Printf("  cp 1 reg_b\n")

		fmt.Printf("  compare\n")
		fmt.Printf("  jump_eq %s_%d\n", labelWhenHead, whenIdx)
		fmt.Printf("  jump %s_%d\n", labelEndWhenHead, whenIdx)

		fmt.Printf("label %s_%d\n", labelWhenHead, whenIdx)

		genStmts(fnArgNames, lvarNames, _rest)

		fmt.Printf("  jump %s\n", labelEnd)
		fmt.Printf("label %s_%d\n", labelEndWhenHead, whenIdx)
	}

	fmt.Printf("label end_case_%d\n", labelId)
	fmt.Printf("  # <<-- case_%d\n", labelId)
	fmt.Printf("\n")
}

func genStmt(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmt *lib.NodeList,
) {
	puts_fn("genStmt")

	stmtHead := head(stmt).Strval
	stmtRest := rest(stmt)

	if stmtHead == "set" {
		genSet(fnArgNames, lvarNames, stmt)
	} else if stmtHead == "call" {
		genCall(fnArgNames, lvarNames, stmt)
	} else if stmtHead == "call_set" {
		genCallSet(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "return" {
		genReturn(lvarNames, stmtRest)
	} else if stmtHead == "while" {
		genWhile(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "case" {
		genCase(fnArgNames, lvarNames, stmtRest)
	} else if stmtHead == "_cmt" {
		genVmComment(stmtRest.Get(0).Strval)
	} else if stmtHead == "_debug" {
		genDebug()
	} else {
		puts_kv_e("stmtHead", stmtHead)
		panic("Unsupported statement")
	}
}

func genStmts(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmts *lib.NodeList,
) {
	for i := 0; i < stmts.Len(); i++ {
		stmt := stmts.Get(i).List
		genStmt(fnArgNames, lvarNames, stmt)
	}
}

func genVar(
	fnArgNames *lib.Names,
	lvarNames *lib.Names,
	stmt *lib.NodeList,
) {
	fmt.Printf("  sub_sp 1\n")

	if stmt.Len() == 3 {
		varName := stmt.Get(1)
		expr := stmt.Get(2)
		_genSet(fnArgNames, lvarNames, varName.Strval, expr)
	}
}

func genFuncDef(topStmt *lib.NodeList) {
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
	asmPrologue()

	fmt.Println("")
	fmt.Println("  # 関数の処理本体")

	for i := 0; i < body.Len(); i++ {
		stmt := body.Get(i).List
		stmtHead := head(stmt).Strval

		if stmtHead == "var" {
			varName := stmt.Get(1).Strval
			lvarNames.Add(varName)
			genVar(fnArgNames, lvarNames, stmt)
		} else {
			genStmt(fnArgNames, lvarNames, stmt)
		}
	}

	fmt.Println("")

	asmEpilogue()
	fmt.Println(`  ret`)
}

func genTopStmts(topStmts *lib.NodeList) {
	topStmtsRest := rest(topStmts)
	for i := 0; i < topStmtsRest.Len(); i++ {
		topStmt := topStmtsRest.Get(i).List
		stmtRest := rest(topStmt)
		genFuncDef(stmtRest)
	}
}

func genBuiltinSetVram() {
	fmt.Println("")
	fmt.Println("label set_vram")
	asmPrologue()

	fmt.Println("  set_vram [bp:2] [bp:3]") // vram_addr value

	asmEpilogue()
	fmt.Println("  ret")
}

func genBuiltinGetVram() {
	fmt.Println("")
	fmt.Println("label get_vram")
	asmPrologue()

	fmt.Println("  get_vram [bp:2] reg_a") // vram_addr dest

	asmEpilogue()
	fmt.Println("  ret")
}

func Codegen() {
	json := lib.ReadStdinAll()
	tree := lib.ParseJson(json)

	fmt.Println(`  call main`)
	fmt.Println(`  exit`)

	genTopStmts(tree)

	fmt.Println("#>builtins")
	genBuiltinSetVram()
	genBuiltinGetVram()
	fmt.Println("#<builtins")
}
