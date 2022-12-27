package main

import (
	"syscall"
	"time"
)

const (
	xWindow  = 0
	yWindow  = 0
	cxWindow = 1000
	cyWindow = 800

	cxScreen int32 = 1920
	cyScreen int32 = 1080
)

func InitWindow() {
	hWnd, err := FindWindow("", processTitle)
	if err != nil {
		Logger.Fatalf("no found process %s", processTitle)
	}
	Logger.Debugf("Found '%s' window: handle=0x%x", processTitle, hWnd)

	if err = SetWindowPos(hWnd, xWindow, yWindow, cxWindow, cyWindow); err != nil {
		Logger.Fatal(err)
	}
	Logger.Debugf("Set '%s' window Pos: x=%d, y=%d, cx=%d, cy=%d", processTitle, xWindow, yWindow, cxWindow, cyWindow)

	if err = SetForegroundWindow(hWnd); err != nil {
		Logger.Fatal(err)
	}
	Logger.Debugf("Foreground '%s' window", processTitle)
}

func CloseWindow() {
	hWnd, err := FindWindow("", processTitle)
	if err != nil {
		Logger.Fatal(err)
	}
	if err = PostMessage(hWnd, WM_CLOSE, 0, 0); err != nil {
		Logger.Fatal(err)
	}
}

func GameFailed() bool {
	hWnd, err := FindWindow("", failedTitle)
	if err != nil {
		return false
	}
	Logger.Debugf("Found '%s' window: handle=0x%x", failedTitle, hWnd)
	return true
}

func GameSuccess() bool {
	hWnd, err := FindWindow("", successTitle)
	if err != nil {
		return false
	}
	Logger.Debugf("Found '%s' window: handle=0x%x", successTitle, hWnd)
	return true
}

// 处理弹窗
func HandlePopup(f func()) {
	if GameFailed() || GameSuccess() {
		f()
		time.Sleep(2 * time.Second)
	}
}

func FindBox() (hWnd syscall.Handle, err error) {
	hWnd, err = FindWindow("", failedTitle)
	if err == nil {
		return
	}
	hWnd, err = FindWindow("", successTitle)
	if err == nil {
		return
	}
	return
}

func ForegroundBox() {
	hWnd, err := FindBox()
	if err != nil {
		return
	}
	if err = SetForegroundWindow(hWnd); err != nil {
		Logger.Fatal(err)
	}
}

// 新游戏
func AgainGame() {
	ForegroundBox()
	Logger.Debugf("Again Game...")
	var keyP uint16 = 80
	if _, err := UniKeyPress(keyP); err != nil {
		panic(err)
	}
}

// 重新开始
func RestartGame() {
	ForegroundBox()
	Logger.Debugf("Restart Game...")
	var keyR uint16 = 82
	if _, err := UniKeyPress(keyR); err != nil {
		panic(err)
	}
}

// 退出游戏
func ExitGame() {
	ForegroundBox()
	Logger.Debugf("Exit Game...")
	var keyX uint16 = 88
	if _, err := UniKeyPress(keyX); err != nil {
		panic(err)
	}
}

func CloseMineSweeper() {
	Logger.Debugf("Close MineSweeper...")
	hWnd, err := FindWindow("", processTitle)
	if err != nil {
		Logger.Fatal(err)
	}
	if err = PostMessage(hWnd, MB_OK, 0, 0); err != nil {
		Logger.Fatal(err)
	}
}

func MoveMouse(row, col int) {
	x, y := RowCol2XY(row, col)
	if err := SetCursorPos(x, y); err != nil {
		Logger.Fatal(err)
	}
	time.Sleep(clickInterval)
}

func RowCol2XY(row, col int) (x int32, y int32) {
	x = int32(startX + col*(endX-startX)/colNum + halfBoxHeight)
	y = int32(startY + row*(endY-startY)/rowNum + halfBoxHeight)
	return
}

func click(buttonType int, row, col int) {
	x, y := RowCol2XY(row, col)
	if _, err := MouseClick(buttonType, x, y); err != nil {
		Logger.Fatal(err)
	}
	time.Sleep(clickInterval)
}

func doubleClick(buttonType int, row, col int) {
	x, y := RowCol2XY(row, col)
	if _, err := MouseDoubleClick(buttonType, x, y); err != nil {
		Logger.Fatal(err)
	}
	time.Sleep(clickInterval)
}
