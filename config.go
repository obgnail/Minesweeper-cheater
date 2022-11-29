package main

import "time"

var (
	// 程序标题
	processTitle = "扫雷"
	// 当游戏失败时弹窗标题
	failedTitle = "游戏失败"
	// 当游戏成功时弹窗标题
	successTitle = "游戏胜利"
	// 鼠标点击时间间隔
	// 过快可能会导致异常: 比如先单击后双击,系统会判定为先双击后单击
	// 若<200ms,useDoubleClick 会自动调整为false
	// 不可<100ms,鼠标操作会错乱,此时会自动调整为100ms
	clickInterval = time.Millisecond * 120
	// 是否使用双击策略
	// 若 !showFlag || clickInterval<200ms 此设置自动调整为false;
	useDoubleClick = false
	// 是否插棋
	showFlag = true
	// 当失败时: restart/again/stop/exit
	whenFailed = "again"
	// 当成功时: again/stop/exit
	whenSuccess = "stop"
)
