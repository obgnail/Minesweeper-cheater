package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	rowNum = 16
	colNum = 30
)

var (
	progress   = false // 如果一轮下来没有进展，随机选择一个
	remainMine = 99

	table  [rowNum][colNum]CellType
	finish [rowNum][colNum]bool // finish cells
	flag   [rowNum][colNum]bool // flag cells
)

type CellType = int

const (
	CellTypeFlag    CellType = -2
	CellTypeUnknown CellType = -1
	CellTypeZero    CellType = 0
	CellTypeOne     CellType = 1
	CellTypeTwo     CellType = 2
	CellTypeThree   CellType = 3
	CellTypeFour    CellType = 4
	CellTypeFive    CellType = 5
	CellTypeSix     CellType = 6
	CellTypeSeven   CellType = 7

	// 程序在策略3-6使用
	CellTypeMine CellType = 8  // 猜测是雷
	CellTypeSafe CellType = 9  // 猜测不是雷
	CellTypeDone CellType = 10 // 标记已处理,避免重复处理
)

type Option = string

const (
	optionRestart Option = "restart"
	optionAgain   Option = "again"
	optionStop    Option = "stop"
	optionExit    Option = "exit"
)

type strategy = string

const (
	strategyRandom strategy = "random"
	strategyEdge   strategy = "edge"
)

func MineSweeper() {
	InitLogger()
	InitCache()
	InitConfig()
	InitStore()
	for {
		InitWindow()
		InitVar()
		success := Play()
		if success {
			WhenSuccess()
		} else {
			WhenFailed()
		}
	}
}

func Play() (success bool) {
	HandlePopup(AgainGame)
	for !Done() {
		if failed := PickAndCheck(); failed {
			return false
		}

		RenewTable()
		OptimizedDetector() // or EasierDetector()
	}
	Logger.Info("--- DONE ---")
	return true
}

// OptimizedDetector 正常来说只要将FindEqual和FindAlways放在一个循环上就行了。
// 分三次迭代后，一次循环场上能排出更多的flag,提高效率
func OptimizedDetector() {
	f1 := ignoreZeroCellDeco(FindEqual)
	f2 := ignoreZeroCellDeco(FindAlways)
	RangeTable(false, f1)
	RangeTable(true, f1)
	RangeTable(false, f2)
}

// EasierDetector 简单地循环检测,没有重复检测,效率较低
func EasierDetector() {
	RangeTable(false, FindEqual)
	RangeTable(false, FindAlways)
}

// 没有进展时随机选择
func PickAndCheck() (failed bool) {
	if !progress {
		if err := RandomPick(); err != nil {
			return true
		}
		if failed = checkFail(); failed { // 检查是否踩雷
			Logger.Info("--- FAILED ---")
			return true
		}
	}
	progress = false
	return false
}

// ShowModalQuickly 双击边框,跳过动画,让模态框快点出现
func ShowModalQuickly() {
	doubleClick(LeftButton, rowNum, colNum)
}

func WhenSuccess() {
	ShowModalQuickly()
	time.Sleep(500 * time.Millisecond)

	switch whenSuccess {
	case optionStop:
		os.Exit(1)
	case optionAgain:
		HandlePopup(AgainGame)
	case optionExit:
		HandlePopup(ExitGame)
		os.Exit(1)
	default:
		os.Exit(1)
	}
}

func WhenFailed() {
	ShowModalQuickly()
	time.Sleep(500 * time.Millisecond)

	switch whenFailed {
	case optionStop:
		os.Exit(1)
	case optionRestart:
		HandlePopup(func() {
			RestartGame()
			ShowModalQuickly() // 去除那该死的【重玩】提示框
		})
	case optionAgain:
		HandlePopup(AgainGame)
	case optionExit:
		HandlePopup(ExitGame)
		os.Exit(1)
	default:
		os.Exit(1)
	}
}

func checkFail() bool {
	time.Sleep(100 * time.Millisecond)
	ShowModalQuickly()
	time.Sleep(600 * time.Millisecond)
	failed := GameFailed()
	return failed
}

func Done() bool {
	// NOTE：当剩下的全部都是雷的时候，系统会直接判赢
	unknown := 0
	RangeTable(false, func(row, col int) {
		if GetCellType(row, col) == CellTypeUnknown {
			unknown++
		}
	})
	if remainMine != unknown {
		return false
	}

	for row := 0; row != rowNum; row++ {
		for col := 0; col != colNum; col++ {
			if finish[row][col] == false {
				return false
			}
		}
	}
	return true
}

func RenewTable() {
	MoveMouse(rowNum, colNum) // 将鼠标移动到范围外，防止错误解析图片

	newTable := Window2Table(func(row, col int) bool {
		ct := table[row][col]
		return ct != CellTypeUnknown && ct != CellTypeDone
	})

	RangeTable(false, func(row, col int) {
		ct := table[row][col]
		if ct != CellTypeUnknown && ct != CellTypeDone {
			newTable[row][col] = table[row][col]
		}
	})
	table = newTable
}

// 策略1-2: 数值=雷数 or 数值=未知单元格数+雷数
func FindEqual(row, col int) {
	neighbors := GetNeighborMap(row, col)
	_unknown := len(neighbors[CellTypeUnknown])
	_flag := len(neighbors[CellTypeFlag])

	switch GetCellType(row, col) {
	case CellTypeZero:
		SetFinish(row, col)
	case _flag:
		if _unknown != 0 {
			ClearCell(row, col, neighbors[CellTypeUnknown])
		}
		SetFinish(row, col)
	case _unknown + _flag:
		for _, cell := range neighbors[CellTypeUnknown] {
			FlagMine(cell.row, cell.col)
		}
		SetFinish(row, col)
	}
}

// 策略7
func RandomPick() (err error) {
	var pickTable []*Location
	RangeTable(false, func(row, col int) {
		if GetCellType(row, col) == CellTypeUnknown {
			pickTable = append(pickTable, &Location{row: row, col: col})
		}
	})

	if len(pickTable) == 0 {
		return fmt.Errorf("len(pickTable) == 0")
	}

	var row, col int
	switch guessStrategy {
	case strategyRandom:
		row, col = _random(pickTable)
	case strategyEdge:
		row, col = _edge(pickTable)
	}
	randomPick(row, col)
	return
}

func _random(pickTable []*Location) (row, col int) {
	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(pickTable))
	c := pickTable[idx]
	return c.row, c.col
}

func _edge(pickTable []*Location) (row, col int) {
	wrb := &WeightRandomBalance{}
	for _, c := range pickTable {
		neighbors := 0
		RangeNeighbors(c.row, c.col, func(idx int, cell *Cell) {
			if IsNumberCellType(cell.CellType) {
				neighbors++
			}
		})
		_ = wrb.Add(&Location{row: c.row, col: c.col}, neighbors+1)
	}
	next := wrb.Next().(*Location)
	return next.row, next.col
}

// 策略3-6：找出在所有情况中，总是雷 or 总是安全 的单元格
// NOTE: Always:在所有的情况中总是
func FindAlways(row, col int) {
	unfinishedNumberNeighbors := GetUnfinishedNumberNeighbors(row, col)
	passSituations := getAllPassSituations(row, col, unfinishedNumberNeighbors)

	switch len(passSituations) {
	case 0:
		return
	case 1:
		// 若只有一种情况,说明所有单元格一定是正确的
		handleAlwaysCell(passSituations[0].cells)
	default:
		// 若多于一种情况,那么就要找出[总是雷]或[总是安全]的单元格
		var alwaysCells []*Cell
		inMyView := alwaysCellsInMyView(passSituations)
		inNeighborsView := alwaysCellsInNeighborsView(passSituations, unfinishedNumberNeighbors)
		alwaysCells = append(alwaysCells, inMyView...)
		alwaysCells = append(alwaysCells, inNeighborsView...)
		alwaysCells = unique(alwaysCells)
		handleAlwaysCell(alwaysCells)
	}
}

// 在当前单元格看来,某一位置的单元格[总是雷]或[总是安全]
func alwaysCellsInMyView(passSituations []*Situation) (alwaysCells []*Cell) {
	if len(passSituations) == 0 {
		return nil
	}

	for cellIdx := 0; cellIdx != len(passSituations[0].cells); cellIdx++ {
		alwaysEqual := true
		cellType0 := passSituations[0].cells[cellIdx].CellType

		for situationIdx := 1; situationIdx != len(passSituations); situationIdx++ {
			cellTypeN := passSituations[situationIdx].cells[cellIdx].CellType
			alwaysEqual = alwaysEqual && (cellType0 == cellTypeN)
			if !alwaysEqual {
				break
			}
		}

		if alwaysEqual {
			alwaysCells = append(alwaysCells, passSituations[0].cells[cellIdx])
		}
	}
	return
}

// 在所有邻居看来,某一位置的单元格[总是雷]或[总是安全]
func alwaysCellsInNeighborsView(passSituations []*Situation, unfinishedNumberNeighbors []*Cell) (alwaysCells []*Cell) {
	for _, tryCell := range unfinishedNumberNeighbors {
		passCount, safeCountMap, mineCountMap := getCountMapByPassSituations(tryCell, passSituations)
		safe := _getAlwaysCell(safeCountMap, passCount) // 总是安全
		mine := _getAlwaysCell(mineCountMap, passCount) // 总是雷
		alwaysCells = append(alwaysCells, safe...)
		alwaysCells = append(alwaysCells, mine...)
	}
	return
}

// 如果单元格是雷的次数不等于情况数，说明该单元格在所有情况中不总是雷
// 同理，如果单元格安全的次数不等于情况数，说明该单元格在所有情况中不总是安全的
func _getAlwaysCell(MapCellToCount map[*Cell]int, passSituationLength int) (alwaysCells []*Cell) {
	for cell, count := range MapCellToCount {
		if count == passSituationLength {
			alwaysCells = append(alwaysCells, cell)
		}
	}
	return
}

// 策略3-6最重要的函数
// 解出 tryCell 周围所有的未开发单元格 在所有可能成立的情况下，是 雷/不是雷 的次数
func getCountMapByPassSituations(tryCell *Cell, passSituations []*Situation) (int, map[*Cell]int, map[*Cell]int) {
	safeCountMap := make(map[string]int)
	mineCountMap := make(map[string]int)
	passCount := len(passSituations) // 合法的情况数
	for _, situation := range passSituations {
		_unique, safeCell, mineCell := resolveSituation(tryCell, situation)
		if !_unique {
			continue
		}
		if len(safeCell) == 0 || len(mineCell) == 0 { // 此策略是不合法的
			passCount--
			continue
		}
		for _, cell := range safeCell {
			key := fmt.Sprintf("%d-%d", cell.row, cell.col)
			safeCountMap[key]++
		}
		for _, cell := range mineCell {
			key := fmt.Sprintf("%d-%d", cell.row, cell.col)
			mineCountMap[key]++
		}
	}

	safe := _toCell(safeCountMap, CellTypeSafe)
	mine := _toCell(mineCountMap, CellTypeMine)

	return passCount, safe, mine
}

func _toCell(countMap map[string]int, cellType CellType) map[*Cell]int {
	res := make(map[*Cell]int)
	for cell, count := range countMap {
		s := strings.Split(cell, "-")
		r, _ := strconv.Atoi(s[0])
		c, _ := strconv.Atoi(s[1])
		_cell := &Cell{
			CellType: cellType,
			Location: &Location{row: r, col: c},
		}
		res[_cell] = count
	}
	return res
}

// 通过(row,col)的邻居、二级邻居筛掉不可能成立的情况，返回所有可能成立的情况
func getAllPassSituations(row, col int, unfinishedNumberNeighbors []*Cell) (passSituations []*Situation) {
	allSituation := getSituationList(row, col)
	for _, situation := range allSituation {
		thisSituationAlwaysPass := true // 此策略对所有人都能成功
		for _, tryCell := range unfinishedNumberNeighbors {
			ok := IsSituationPass(tryCell, situation)
			thisSituationAlwaysPass = thisSituationAlwaysPass && ok
			if !thisSituationAlwaysPass {
				break
			}
		}
		if thisSituationAlwaysPass {
			passSituations = append(passSituations, situation)
		}
	}
	return passSituations
}

func handleAlwaysCell(cells []*Cell) {
	for _, cell := range cells {
		switch cell.CellType {
		case CellTypeSafe:
			FlagSafe(cell.row, cell.col)
		case CellTypeMine:
			FlagMine(cell.row, cell.col)
		}
	}
}

// 更新neighbors
func renewNeighbor(neighbors []*Cell, situation *Situation) {
	situation.RangeCell(func(idx int, c *Cell) (stop bool) {
		if c.CellType != CellTypeMine && c.CellType != CellTypeSafe {
			return false
		}
		for _, n := range neighbors {
			if !(n.row == c.row && n.col == c.col) || n.CellType != CellTypeUnknown {
				continue
			}
			n.CellType = c.CellType
			break
		}
		return false
	})
}

// 根据situation,解出所有的Unknown
// unique: 是否只有一个解
//	- 若存在多个解,返回false,nil,nil;
//	- 若只有一个解且此解错误,返回true,nil,nil
//	- 若只有一个解且此解正确,返回此解(true,slice,slice)
func resolveSituation(cell *Cell, situation *Situation) (unique bool, safe []*Cell, mine []*Cell) {
	neighbors := GetNeighborListBySituation(cell.row, cell.col, situation)
	_flag, _unknown, _mine := 0, 0, 0
	for _, n := range neighbors {
		switch n.CellType {
		case CellTypeFlag:
			_flag++
		case CellTypeMine:
			_mine++
		case CellTypeUnknown:
			_unknown++
		}
	}

	value := cell.CellType
	// 存在多个解
	if _flag+_mine < value && value < _unknown+_flag+_mine {
		unique = false
		return
	}

	for _, n := range neighbors {
		if n.CellType == CellTypeUnknown {
			if _mine+_flag < value {
				n.CellType = CellTypeMine
			} else {
				n.CellType = CellTypeSafe
			}
		}
	}

	mapByType := make(map[CellType][]*Cell)
	for _, n := range neighbors {
		mapByType[n.CellType] = append(mapByType[n.CellType], n)
	}

	if len(mapByType[CellTypeMine])+len(mapByType[CellTypeFlag]) != value {
		return true, nil, nil
	}
	safe = mapByType[CellTypeSafe]
	mine = mapByType[CellTypeMine]

	return true, safe, mine
}

func RangeTable(reverse bool, rangeFunc func(row, col int)) {
	if rangeFunc == nil {
		panic("rangeFunc must not nil")
	}
	if !reverse {
		for row := 0; row != rowNum; row++ {
			for col := 0; col != colNum; col++ {
				rangeFunc(row, col)
			}
		}
	} else {
		for col := colNum - 1; col != -1; col-- {
			for row := rowNum - 1; row != -1; row-- {
				rangeFunc(row, col)
			}
		}
	}
}

// 优化装饰器:对于空白的Cell本身就是Finished状态了，不必执行后续的任务了
func ignoreZeroCellDeco(f func(row, col int)) func(row, col int) {
	return func(row, col int) {
		if GetCellType(row, col) == CellTypeZero || IsFinish(row, col) {
			SetFinish(row, col)
			return
		}
		f(row, col)
	}
}

func GetCellType(row, col int) CellType {
	if row < 0 || col < 0 || row > rowNum-1 || col > colNum-1 {
		return 0
	}
	return table[row][col]
}

// Warn:优化使用的函数,没事别改动table元素的值
func setTableCell(row, col int, cellType CellType) {
	table[row][col] = cellType
}

type Cell struct {
	CellType
	*Location
}

func (c *Cell) String() string {
	return fmt.Sprintf("(%d,%d)[%d]", c.row, c.col, c.CellType)
}

type Location struct {
	row int
	col int
}

func (l *Location) Equal(row, col int) bool {
	return l.row == row && l.col == col
}

func (l *Location) IsNeighbor(row, col int) bool {
	return Abs(row-l.row) <= 1 && Abs(col-l.col) <= 1
}

func (l *Location) IsLegal() bool {
	if l.row < 0 || l.col < 0 || l.row > rowNum-1 || l.col > colNum-1 {
		return false
	}
	return true
}

// 对于cell来说,情况是否可能成立
func IsSituationPass(cell *Cell, situation *Situation) bool {
	neighbors := GetNeighborListBySituation(cell.row, cell.col, situation)
	_flag, _unknown, _mine := 0, 0, 0
	for _, n := range neighbors {
		switch n.CellType {
		case CellTypeFlag:
			_flag++
		case CellTypeMine:
			_mine++
		case CellTypeUnknown:
			_unknown++
		}
	}
	value := cell.CellType
	if _mine+_flag > value || _unknown+_mine+_flag < value {
		return false
	}
	return true
}

type Situation struct {
	cells []*Cell
}

func (s *Situation) RangeCell(rangeFunc func(idx int, cell *Cell) (stop bool)) {
	for idx, cell := range s.cells {
		if stop := rangeFunc(idx, cell); stop {
			break
		}
	}
}

func (s *Situation) String() string {
	var mine, safe []string
	for _, cell := range s.cells {
		str := fmt.Sprintf("%s", cell)
		if cell.CellType == CellTypeMine {
			mine = append(mine, str)
		} else if cell.CellType == CellTypeSafe {
			safe = append(safe, str)
		}
	}
	return fmt.Sprintf(
		"Mine: %s | Safe: %s",
		strings.Join(mine, ", "),
		strings.Join(safe, ", "),
	)
}

// 对于(row, col)来说,雷的所有可能分布情况
func getSituationList(row, col int) []*Situation {
	_flag := 0               // 已探明雷的数量
	var _unknown []*Location // 可能是雷的位置
	RangeNeighbors(row, col, func(idx int, cell *Cell) {
		switch cell.CellType {
		case CellTypeFlag:
			_flag++
		case CellTypeUnknown:
			_unknown = append(_unknown, cell.Location)
		}
	})

	unknownRemain := len(_unknown)                                   // 剩余未知的数量
	mineRemain := GetCellType(row, col) - _flag                      // 剩余的雷数量
	situations := GetCombinationFromCache(unknownRemain, mineRemain) // 总共有这么多种情况

	var res []*Situation
	for _, situation := range situations {
		var cells []*Cell
		for idx, value := range situation {
			var ct CellType
			if value == 1 {
				ct = CellTypeMine
			} else {
				ct = CellTypeSafe
			}
			cells = append(cells, &Cell{
				CellType: ct,
				Location: _unknown[idx],
			})
		}
		res = append(res, &Situation{cells: cells})
	}
	return res
}

func RangeNeighbors(row, col int, rangeFunc func(idx int, cell *Cell)) {
	if rangeFunc == nil {
		panic("rangeFunc must not nil")
	}
	searchLocation := [][2]int{
		{row - 1, col - 1}, {row - 1, col}, {row - 1, col + 1}, {row, col + 1},
		{row + 1, col + 1}, {row + 1, col}, {row + 1, col - 1}, {row, col - 1},
	}
	for idx, loc := range searchLocation {
		rangeFunc(idx, &Cell{
			CellType: GetCellType(loc[0], loc[1]),
			Location: &Location{loc[0], loc[1]},
		})
	}
}

func GetNeighborMap(row, col int) map[CellType][]*Cell {
	res := make(map[CellType][]*Cell)
	RangeNeighbors(row, col, func(idx int, cell *Cell) {
		res[cell.CellType] = append(res[cell.CellType], cell)
	})
	return res
}

func GetNeighborList(row, col int) []*Cell {
	var res []*Cell
	RangeNeighbors(row, col, func(idx int, cell *Cell) {
		res = append(res, cell)
	})
	return res
}

func GetNeighborListBySituation(row, col int, situation *Situation) []*Cell {
	neighbors := GetNeighborList(row, col)
	renewNeighbor(neighbors, situation)
	return neighbors
}

// 包含一级邻居 和 以unknownCell为中心的二级邻居
func GetUnfinishedNumberNeighbors(row, col int) []*Cell {
	var cells []*Cell
	var unknown []*Cell
	RangeNeighbors(row, col, func(idx int, cell *Cell) {
		if cell.IsLegal() && IsNumberCellType(cell.CellType) && !IsFinish(cell.row, cell.col) {
			cells = append(cells, cell)
		}
		if cell.CellType == CellTypeUnknown {
			unknown = append(unknown, cell)
		}
	})

	for _, c := range unknown {
		RangeNeighbors(c.row, c.col, func(idx int, cell *Cell) {
			if cell.IsLegal() && IsNumberCellType(cell.CellType) && !IsFinish(cell.row, cell.col) && !cell.Equal(row, col) {
				cells = append(cells, cell)
			}
		})
	}

	res := unique(cells)
	return res
}

// 以后都不必再看这个cell
func SetFinish(row, col int) {
	finish[row][col] = true
}

func IsFinish(row, col int) bool {
	return finish[row][col] == true
}

func IsNumberCellType(cellType CellType) bool {
	return CellTypeZero < cellType && cellType < CellTypeMine
}

func ClearCell(row, col int, unknownCell []*Cell) {
	progress = true
	// 只有标记了flag的情况下才可以使用双击,否则需要一个个去点
	if useDoubleClick {
		Logger.Debugf("DoubleClick: (%d,%d)", row, col)
		doubleClick(LeftButton, row, col)
		for _, cell := range unknownCell {
			setTableCell(cell.row, cell.col, CellTypeDone)
		}
	} else {
		for _, cell := range unknownCell {
			FlagSafe(cell.row, cell.col)
		}
	}
	SetFinish(row, col)
}

func randomPick(row, col int) {
	Logger.Debugf("RadomPick: (%d,%d)", row, col)
	progress = true
	click(LeftButton, row, col)
	setTableCell(row, col, CellTypeDone)
}

func FlagSafe(row, col int) {
	if GetCellType(row, col) == CellTypeDone {
		return
	}
	Logger.Debugf("FlagSafe: (%d,%d)", row, col)
	progress = true
	click(LeftButton, row, col)
	setTableCell(row, col, CellTypeDone)
}

func FlagMine(row, col int) {
	if ct := GetCellType(row, col); ct == CellTypeDone || ct == CellTypeFlag {
		return
	}
	progress = true
	remainMine--
	flag[row][col] = true
	if showFlag {
		Logger.Debugf("FlagMine: (%d,%d)", row, col)
		click(RightButton, row, col)
	}
	setTableCell(row, col, CellTypeFlag)
}

func unique(cells []*Cell) []*Cell {
	m := make(map[string]*Cell)
	for _, c := range cells {
		key := fmt.Sprintf("%d-%d", c.row, c.col)
		m[key] = c
	}
	var res []*Cell
	for _, cell := range m {
		res = append(res, cell)
	}
	return res
}

func InitVar() {
	progress = false
	remainMine = 99
	for row := 0; row != rowNum; row++ {
		table[row] = [colNum]CellType{}
		finish[row] = [colNum]bool{}
		flag[row] = [colNum]bool{}

		for col := 0; col != colNum; col++ {
			table[row][col] = CellTypeUnknown
			finish[row][col] = false
			flag[row][col] = false
		}
	}
}
