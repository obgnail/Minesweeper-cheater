package main

import "testing"

func InitTestTable() {
	table[0][0] = CellTypeUnknown
	table[0][1] = CellTypeUnknown
	table[0][2] = CellTypeUnknown
	table[0][3] = CellTypeUnknown

	table[1][0] = CellTypeUnknown
	table[1][1] = CellTypeUnknown
	table[1][2] = CellTypeThree
	table[1][3] = CellTypeUnknown

	table[2][0] = CellTypeUnknown
	table[2][1] = CellTypeOne
	table[2][2] = CellTypeFour
	table[2][3] = CellTypeUnknown

	table[3][0] = CellTypeUnknown
	table[3][1] = CellTypeUnknown
	table[3][2] = CellTypeUnknown
	table[3][3] = CellTypeUnknown
}
func InitTestTable2() {
	table[0][0] = CellTypeUnknown
	table[0][1] = CellTypeUnknown
	table[0][2] = CellTypeUnknown
	table[0][3] = CellTypeUnknown
	table[0][4] = CellTypeUnknown

	table[1][0] = CellTypeOne
	table[1][1] = CellTypeOne
	table[1][2] = CellTypeTwo
	table[1][3] = CellTypeFlag
	table[1][4] = CellTypeUnknown

	table[2][0] = CellTypeOne
	table[2][1] = CellTypeFlag
	table[2][2] = CellTypeThree
	table[2][3] = CellTypeTwo
	table[2][4] = CellTypeUnknown

	table[3][0] = CellTypeOne
	table[3][1] = CellTypeThree
	table[3][2] = CellTypeUnknown
	table[3][3] = CellTypeUnknown
	table[3][4] = CellTypeUnknown
}

func TestFunc(t *testing.T) {
	InitLogger()
	InitFlag()
	InitTable()
	InitTestTable()
	Range(false, FindAlways)
	t.Log("--- done ---")
}

func TestFunc2(t *testing.T) {
	InitLogger()
	InitFlag()
	InitTable()
	InitTestTable2()
	Range(false, FindAlways)
	t.Log("--- done ---")
}
