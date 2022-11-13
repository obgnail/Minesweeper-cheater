package main

import (
	"fmt"
	"math/rand"
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

type CellType int

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

	// 程序在策略3和策略4使用
	CellTypeMine CellType = 8
	CellTypeSafe CellType = 9
	// 优化使用
	CellTypeSafe2 CellType = 10
)

func Play() {
	for !Done() {
		RenewTable()
		// 正常来说只要将ValueEqualFlagAndUnknown和inductionOnUnknowns放在一个循环上就行了。
		// 分四次迭代后，场上能排出更多的flag,提高效率
		f1 := ignoreZeroCellDeco(ValueEqualFlagAndUnknown)
		f2 := ignoreZeroCellDeco(InductionOnUnknowns)
		Range(false, f1)
		Range(true, f1)
		Range(false, f2)
		if !progress {
			RandomPick() // 没有进展时随机选择
		}
		progress = false
	}
	Logger.Info("--- done ---")
}

// 效果等同于PlayGame,不过效率不佳
func PlayGame2() {
	for !Done() {
		RenewTable()
		Range(false, ValueEqualFlagAndUnknown)
		Range(false, InductionOnUnknowns)
		if !progress {
			RandomPick()
		}
		progress = false
	}
}

func Done() bool {
	// NOTE：当剩下的全部都是雷的时候，系统会直接判赢
	unknown := 0
	Range(false, func(row, col int) {
		if GetCellType(row, col) == CellTypeUnknown {
			unknown++
		}
	})
	if remainMine == unknown {
		return true
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
	table = Window2Table()
	if !showFlag {
		updateTable()
	}
}

func ValueEqualFlagAndUnknown(row, col int) {
	neighbors := GetNeighborMap(row, col)
	unknownNum := len(neighbors[CellTypeUnknown])
	flagNum := len(neighbors[CellTypeFlag])
	value := int(GetCellType(row, col))

	switch value {
	case int(CellTypeZero):
		SetFinish(row, col)
	case flagNum:
		if unknownNum != 0 {
			//Logger.Debug("Strategy1(ClearCell):")
			ClearCell(row, col, neighbors[CellTypeUnknown])
		}
		SetFinish(row, col)
	case unknownNum + flagNum:
		//Logger.Debug("Strategy2(MineFlag):")
		for _, cell := range neighbors[CellTypeUnknown] {
			FlagMine(cell.row, cell.col)
		}
		SetFinish(row, col)
	}
}

func RandomPick() {
	var pickTable []*Location
	Range(false, func(row, col int) {
		if GetCellType(row, col) == CellTypeUnknown {
			pickTable = append(pickTable, &Location{row: row, col: col})
		}
	})

	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(pickTable))
	c := pickTable[idx]
	//Logger.Debug("Strategy5(RandomPick):")
	randomPick(c.row, c.col)
}

func InductionOnUnknowns(row, col int) {
	unfinishedNumberNeighbors := GetUnfinishedNumberNeighbors(row, col)
	passSituations := getAllPassSituations(row, col, unfinishedNumberNeighbors)

	switch len(passSituations) {
	case 0:
		return
	case 1:
		handlePassSituationCell(passSituations[0].cells)
	default:
		for _, tryCell := range unfinishedNumberNeighbors {
			safeCountMap, mineCountMap := getCountMapByPassSituations(tryCell, passSituations)
			l := len(passSituations)
			safe := _getAlwaysPassCell(safeCountMap, l)
			mine := _getAlwaysPassCell(mineCountMap, l)
			handlePassSituationCell(safe)
			handlePassSituationCell(mine)
		}
	}
}

func _getAlwaysPassCell(MapCellToCount map[*Cell]int, passSituations int) []*Cell {
	var res []*Cell
	for cell, count := range MapCellToCount {
		if count == passSituations {
			res = append(res, cell)
		}
	}
	return res
}

// 策略3-6最重要的函数
// 解出 tryCell周围所有的未开发单元格 在所有可能成立的情况下，是 雷/不是雷 的次数
func getCountMapByPassSituations(tryCell *Cell, passSituations []*Situation) (map[*Cell]int, map[*Cell]int) {
	safeCountMap := make(map[string]int)
	mineCountMap := make(map[string]int)
	for _, situation := range passSituations {
		safeCell, mineCell := resolveSituation(tryCell, situation)
		if len(safeCell) == 0 || len(mineCell) == 0 {
			continue // 此策略是不合理的
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

	toCell := func(m map[string]int, cellType CellType) map[*Cell]int {
		res := make(map[*Cell]int)
		for cell, count := range m {
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

	safe := toCell(safeCountMap, CellTypeSafe)
	mine := toCell(mineCountMap, CellTypeMine)

	return safe, mine
}

// 通过(row,col)的邻居、二级邻居筛掉不可能成立的情况，返回所有可能成立的情况
func getAllPassSituations(row, col int, unfinishedNumberNeighbors []*Cell) []*Situation {
	allSituation := getSituationList(row, col)

	var passSituations []*Situation
	for _, situation := range allSituation {
		thisSituationAlwaysPass := true // 此策略对所有人都能成功
		for _, tryCell := range unfinishedNumberNeighbors {
			ok := IsSituationPass(tryCell, situation)
			thisSituationAlwaysPass = thisSituationAlwaysPass && ok
			if !thisSituationAlwaysPass {
				break
			}
		}
		passSituations = append(passSituations, situation)
	}
	return passSituations
}

func handlePassSituationCell(cells []*Cell) {
	for _, cell := range cells {
		switch cell.CellType {
		case CellTypeSafe:
			//Logger.Debug("Strategy3(FlagSafe):")
			FlagSafe(cell.row, cell.col)
			setTableCell(cell.row, cell.col, CellTypeSafe2)
		case CellTypeMine:
			//Logger.Debug("Strategy4(FlagMine):")
			FlagMine(cell.row, cell.col)
		}
	}
}

func handlePassSituation333(passSituation *Situation) {
	passSituation.RangeCell(func(idx int, cell *Cell) (stop bool) {
		switch cell.CellType {
		case CellTypeSafe:
			//Logger.Debug("Strategy3(FlagSafe):")
			FlagSafe(cell.row, cell.col)
			setTableCell(cell.row, cell.col, CellTypeSafe2)
		case CellTypeMine:
			//Logger.Debug("Strategy4(FlagMine):")
			FlagMine(cell.row, cell.col)
		}
		return false
	})
}

func resolveSituation2(cell *Cell, neighbors map[CellType][]*Cell) (safe []*Cell, mine []*Cell) {
	if len(neighbors[CellTypeMine])+len(neighbors[CellTypeFlag]) != int(cell.CellType) {
		return nil, nil
	}
	safe = neighbors[CellTypeSafe]
	mine = neighbors[CellTypeMine]
	return safe, mine
}

// 根据situation,解出所有的Unknown
func resolveSituation(cell *Cell, situation *Situation) ([]*Cell, []*Cell) {
	neighbors := GetNeighborList(cell.row, cell.col)
	dugInSituation := 0

	situation.RangeCell(func(idx int, cell *Cell) (stop bool) {
		if cell.CellType != CellTypeMine && cell.CellType != CellTypeSafe {
			return false
		}
		for _, n := range neighbors {
			if n.row == cell.row && n.col == cell.col {
				n.CellType = cell.CellType
			}
		}
		if cell.CellType == CellTypeMine {
			dugInSituation++
		}
		return false
	})

	dugInNeighbors := 0
	for _, n := range neighbors {
		if n.CellType == CellTypeFlag {
			dugInNeighbors++
		}
	}

	for _, n := range neighbors {
		if n.CellType == CellTypeUnknown {
			if dugInSituation+dugInNeighbors < int(cell.CellType) {
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

	if len(mapByType[CellTypeMine])+len(mapByType[CellTypeFlag]) != int(cell.CellType) {
		return nil, nil
	}
	safe := mapByType[CellTypeSafe]
	mine := mapByType[CellTypeMine]

	return safe, mine
}

func Range(reverse bool, f func(row, col int)) {
	if !reverse {
		for row := 0; row != rowNum; row++ {
			for col := 0; col != colNum; col++ {
				f(row, col)
			}
		}
	} else {
		for col := colNum - 1; col != -1; col-- {
			for row := rowNum - 1; row != -1; row-- {
				f(row, col)
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
	isMine := 0
	situation.RangeCell(func(idx int, cell *Cell) (stop bool) {
		if cell.CellType == CellTypeMine && cell.IsNeighbor(cell.row, cell.col) {
			isMine++
		}
		return false
	})

	value := int(cell.CellType)
	m := GetNeighborMap(cell.row, cell.col)
	flag := len(m[CellTypeFlag])
	unknown := len(m[CellTypeUnknown])

	if isMine+flag > value || unknown+flag < value {
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
	var isMine []string
	var notMine []string
	for _, cell := range s.cells {
		str := fmt.Sprintf("(%d,%d)", cell.row, cell.col)
		if cell.CellType == CellTypeMine {
			isMine = append(isMine, str)
		} else if cell.CellType == CellTypeSafe {
			notMine = append(notMine, str)
		}
	}
	str := "Mine:"
	str += strings.Join(isMine, ", ")
	str += " | Safe:"
	str += strings.Join(notMine, ", ")
	return str
}

// 对于(row, col)来说,雷的所有可能分布情况
func getSituationList(row, col int) []*Situation {
	searchLocation := [][2]int{
		{row - 1, col - 1}, {row - 1, col}, {row - 1, col + 1}, {row, col + 1},
		{row + 1, col + 1}, {row + 1, col}, {row + 1, col - 1}, {row, col - 1},
	}
	dugCell := 0            // 已探明雷的数量
	var unknown []*Location // 可能是雷的位置
	for _, loc := range searchLocation {
		switch GetCellType(loc[0], loc[1]) {
		case CellTypeFlag:
			dugCell++
		case CellTypeUnknown:
			unknown = append(unknown, &Location{loc[0], loc[1]})
		}
	}

	unknownRemain := len(unknown)                        // 剩余未知的数量
	mineRemain := int(GetCellType(row, col)) - dugCell   // 剩余的雷数量
	situations := Combination(unknownRemain, mineRemain) // 总共有这么多种情况

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
				Location: unknown[idx],
			})
		}
		res = append(res, &Situation{cells: cells})
	}
	return res
}

func GetNeighborsByFunc(row, col int, pickFunc func(idx int, cell *Cell)) {
	searchLocation := [][2]int{
		{row - 1, col - 1}, {row - 1, col}, {row - 1, col + 1}, {row, col + 1},
		{row + 1, col + 1}, {row + 1, col}, {row + 1, col - 1}, {row, col - 1},
	}
	for idx, loc := range searchLocation {
		pickFunc(idx, &Cell{
			CellType: GetCellType(loc[0], loc[1]),
			Location: &Location{loc[0], loc[1]},
		})
	}
}

func GetNeighborMap(row, col int) map[CellType][]*Cell {
	res := make(map[CellType][]*Cell)
	GetNeighborsByFunc(row, col, func(idx int, cell *Cell) {
		res[cell.CellType] = append(res[cell.CellType], cell)
	})
	return res
}

func GetNeighborList(row, col int) []*Cell {
	var res []*Cell
	GetNeighborsByFunc(row, col, func(idx int, cell *Cell) {
		res = append(res, cell)
	})
	return res
}

// 包含一级邻居 和 以unknownCell为中心的二级邻居
func GetUnfinishedNumberNeighbors(row, col int) []*Cell {
	var cells []*Cell
	var unknown []*Cell
	GetNeighborsByFunc(row, col, func(idx int, cell *Cell) {
		if cell.IsLegal() && IsNumberCellType(cell.CellType) && !IsFinish(cell.row, cell.col) {
			cells = append(cells, cell)
		}
		if cell.CellType == CellTypeUnknown {
			unknown = append(unknown, cell)
		}
	})

	for _, c := range unknown {
		GetNeighborsByFunc(c.row, c.col, func(idx int, cell *Cell) {
			if cell.IsLegal() && IsNumberCellType(cell.CellType) && !IsFinish(cell.row, cell.col) {
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
	t := int(cellType)
	return 0 < t && t < 8
}

func ClearCell(row, col int, unknownCell []*Cell) {
	progress = true
	// 只有标记了flag的情况下才可以使用双击,否则需要一个个去点
	if showFlag {
		Logger.Debugf("DoubleClick: (%d,%d)", row, col)
		doubleClick(LeftButton, row, col)
	} else {
		for _, cell := range unknownCell {
			FlagSafe(cell.row, cell.col)
		}
	}
	for _, cell := range unknownCell {
		setTableCell(cell.row, cell.col, CellTypeSafe2)
	}
	SetFinish(row, col)
}

func randomPick(row, col int) {
	Logger.Debugf("RadomPick: (%d,%d)", row, col)
	progress = true
	click(LeftButton, row, col)
}

func FlagSafe(row, col int) {
	Logger.Debugf("FlagSafe: (%d,%d)", row, col)
	progress = true
	click(LeftButton, row, col)
}

func FlagMine(row, col int) {
	progress = true
	remainMine--
	table[row][col] = CellTypeFlag
	flag[row][col] = true
	if showFlag {
		Logger.Debugf("FlagMine: (%d,%d)", row, col)
		click(RightButton, row, col)
	}
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

func updateTable() {
	Range(false, func(row, col int) {
		if flag[row][col] == true {
			setTableCell(row, col, CellTypeFlag)
		}
	})
}

func InitTable() {
	for row := 0; row != rowNum; row++ {
		table[row] = [colNum]CellType{}
		for col := 0; col != colNum; col++ {
			table[row][col] = 0
		}
	}

	for row := 0; row != rowNum; row++ {
		finish[row] = [colNum]bool{}
		for col := 0; col != colNum; col++ {
			finish[row][col] = false
		}
	}
}
