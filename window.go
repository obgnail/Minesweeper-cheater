package main

import (
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
		Logger.Fatal(err)
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
	_, err := FindWindow("", failedTitle)
	if err != nil {
		return false
	}
	return true
}

// 新游戏
func AgainGame() {
	var keyP uint16 = 80
	if _, err := UniKeyPress(keyP); err != nil {
		panic(err)
	}
}

// 重新开始
func RestartGame() {
	var keyR uint16 = 82
	if _, err := UniKeyPress(keyR); err != nil {
		panic(err)
	}
}

// 退出游戏
func ExitGame() {
	var keyX uint16 = 88
	if _, err := UniKeyPress(keyX); err != nil {
		panic(err)
	}
}

func CloseMessageBox() {
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
