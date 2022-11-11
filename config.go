package main

import "time"

const (
	// 程序标题
	processTitle = "扫雷"
	// 鼠标点击时间间隔(过快可能会导致异常,比如两次单击变成双击)
	clickInterval = time.Millisecond * 300
	// 是否插棋
	ShowFlag = false
)
