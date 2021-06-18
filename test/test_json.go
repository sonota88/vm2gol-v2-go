package test

import (
	"github.com/sonota88/vm2gol-v2-go/lib"
)

func make_test_json_data_1() *lib.NodeList {
	return lib.NodeList_empty();
}

func make_test_json_data_2() *lib.NodeList {
	var tree lib.NodeList

	tree.Add(lib.Node_newInt(1))

	return &tree
}

func make_test_json_data_3() *lib.NodeList {
	var tree lib.NodeList

	tree.Add(lib.Node_newStr("fdsa"))

	return &tree
}

func make_test_json_data_4() *lib.NodeList {
	var tree lib.NodeList

	tree.Add(lib.Node_newInt(-123))
	tree.Add(lib.Node_newStr("fdsa"))

	return &tree
}

func make_test_json_data_5() lib.NodeList {
	var tree lib.NodeList
	var innerList lib.NodeList

	item0 := lib.Node_newList(&innerList)
	tree.Add(item0)

	return tree
}

func make_test_json_data_6() lib.NodeList {
	var tree lib.NodeList

	tree.Add(lib.Node_newInt(1))
	tree.Add(lib.Node_newStr("a"))

	{
		var innerList lib.NodeList
		innerList.Add(lib.Node_newInt(2))
		innerList.Add(lib.Node_newStr("b"))
		tree.Add(lib.Node_newList(&innerList))
	}

	tree.Add(lib.Node_newInt(3))
	tree.Add(lib.Node_newStr("c"))

	return tree
}

func TestJson() {
	// tree := make_test_json_data_1();
	// tree := make_test_json_data_2();
	// tree := make_test_json_data_3();
	// tree := make_test_json_data_4();
	// tree := make_test_json_data_5();
	// tree := make_test_json_data_6();

	json := lib.ReadStdinAll()
	// lib.Puts_e(json)

	tree := lib.ParseJson(json)
	lib.PrintAsJson(tree)
}
