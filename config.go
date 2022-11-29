package main

import "time"

const (
	// 程序标题
	processTitle = "扫雷"
	// 当游戏失败时弹窗标题
	failedTitle = "游戏失败"
	// 当游戏成功时弹窗标题
	successTitle = "游戏胜利"
	// 鼠标点击时间间隔,过快可能会导致异常:
	// 若<200ms,务必将showFlag设置为false.因为此时单元格可能会先单击后双击,系统对此会判定为先双击后单击
	// 不可<100ms.因为系统顶不住,鼠标操作会错乱
	clickInterval = time.Millisecond * 100
	// 是否插棋
	showFlag = false
)
