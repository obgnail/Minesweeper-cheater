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

func MoveMouse(row, col int) {
	x, y := RowCol2XY(row, col)
	if err := SetCursorPos(x, y); err != nil {
		panic(err)
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
		panic(err)
	}
	time.Sleep(clickInterval)
}

func doubleClick(buttonType int, row, col int) {
	x, y := RowCol2XY(row, col)
	if _, err := MouseDoubleClick(buttonType, x, y); err != nil {
		panic(err)
	}
	time.Sleep(clickInterval)
}
