package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

const (
	title     = "扫雷"
	sleepTime = time.Millisecond * 300
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
	logrus.SetLevel(logrus.DebugLevel)

	hWnd, err := FindWindow("", title)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Debugf("Found '%s' window: handle=0x%x", title, hWnd)

	if err = SetWindowPos(hWnd, xWindow, yWindow, cxWindow, cyWindow); err != nil {
		log.Fatal(err)
	}
	logrus.Debugf("Set '%s' window Pos: x=%d, y=%d, cx=%d, cy=%d", title, xWindow, yWindow, cxWindow, cyWindow)

	if err = SetForegroundWindow(hWnd); err != nil {
		log.Fatal(err)
	}
	logrus.Debugf("Foreground '%s' window", title)
}

func MoveMouse(row, col int) {
	x, y := RowCol2XY(row, col)
	if err := SetCursorPos(x, y); err != nil {
		panic(err)
	}
	time.Sleep(sleepTime)
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
	time.Sleep(sleepTime)
}

func doubleClick(buttonType int, row, col int) {
	x, y := RowCol2XY(row, col)
	if _, err := MouseDoubleClick(buttonType, x, y); err != nil {
		panic(err)
	}
	time.Sleep(sleepTime)
}
