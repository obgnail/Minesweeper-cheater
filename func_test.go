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

func InitTestTable3() {
	table[0][0] = CellTypeUnknown
	table[0][1] = CellTypeOne
	table[0][2] = CellTypeZero
	table[0][3] = CellTypeZero

	table[1][0] = CellTypeUnknown
	table[1][1] = CellTypeTwo
	table[1][2] = CellTypeZero
	table[1][3] = CellTypeZero

	table[2][0] = CellTypeUnknown
	table[2][1] = CellTypeOne
	table[2][2] = CellTypeZero
	table[2][3] = CellTypeZero

	table[3][0] = CellTypeUnknown
	table[3][1] = CellTypeThree
	table[3][2] = CellTypeUnknown
	table[3][3] = CellTypeUnknown
}

func TestFunc(t *testing.T) {
	InitLogger()
	InitConfig()
	InitVar()
	InitTestTable()
	RangeTable(false, FindAlways)
	t.Log("--- done ---")
}

func TestFunc2(t *testing.T) {
	InitLogger()
	InitConfig()
	InitVar()
	InitTestTable2()
	RangeTable(false, FindAlways)
	t.Log("--- done ---")
}

func TestFunc3(t *testing.T) {
	InitLogger()
	InitConfig()
	InitVar()
	InitTestTable3()
	RangeTable(false, FindAlways)
	t.Log("--- done ---")
}

func TestFunc4(t *testing.T) {
	//CloseWindow()
	CloseMineSweeper()
}

func TestFunc5(t *testing.T) {
	failed := GameFailed()
	t.Log(failed)

	RestartGame()
}

func TestFunc6(t *testing.T) {
	var keyP uint16 = 80
	if _, err := UniKeyPress(keyP); err != nil {
		panic(err)
	}
}

func TestCombination(t *testing.T) {
	res := Combination(8, 8)
	t.Log(res)

	InitCache()
	res3 := GetCombinationFromCache(8, 8)
	t.Log(res3)

	for i := 1; i != 9; i++ {
		for j := 1; j != i+1; j++ {
			t.Log(i, j, GetCombinationFromCache(i, j))
		}
	}
}

func TestPopup(t *testing.T) {
	success := GameSuccess()
	t.Log(success)
}
