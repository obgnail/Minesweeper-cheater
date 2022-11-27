package main

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32 = syscall.MustLoadDLL("user32.dll")

	procFindWindow          = user32.MustFindProc("FindWindowW")
	procSetWindowPos        = user32.MustFindProc("SetWindowPos")
	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")

	procSetCursorPos = user32.MustFindProc("SetCursorPos")
	procSendInput    = user32.MustFindProc("SendInput")
	procPostMessage  = user32.MustFindProc("PostMessageW")
)

func handlerErrNo(r1, r2 uintptr, errNo syscall.Errno) (err error) {
	if 0 == int32(r1) {
		if errNo != 0 {
			err = error(errNo)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// stringToUTF16PtrElseNil String To UTF16Ptr if empty string trans to nil
func stringToUTF16PtrElseNil(str string) *uint16 {
	if str == "" {
		return nil
	}
	return syscall.StringToUTF16Ptr(str)
}

// FindWindow find window hWnd by name class="" if nil, nil mean ignore it
func FindWindow(class, title string) (syscall.Handle, error) {
	c := unsafe.Pointer(stringToUTF16PtrElseNil(class))
	t := unsafe.Pointer(stringToUTF16PtrElseNil(title))
	r1, r2, errNo := syscall.Syscall(procFindWindow.Addr(), 2, uintptr(c), uintptr(t), 0)
	err := handlerErrNo(r1, r2, errNo)
	return syscall.Handle(r1), err
}

func SetWindowPos(hWnd syscall.Handle, x, y, cx, cy int32) (err error) {
	HWND_TOP := 0
	SWP_SHOWWINDOW := 0x0040

	r1, r2, errNo := syscall.Syscall9(
		procSetWindowPos.Addr(),
		7,
		uintptr(hWnd),
		uintptr(HWND_TOP),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(SWP_SHOWWINDOW),
		0, 0)
	err = handlerErrNo(r1, r2, errNo)
	return
}

func SetForegroundWindow(hWnd syscall.Handle) (err error) {
	r1, r2, errNo := syscall.Syscall(procSetForegroundWindow.Addr(), 1, uintptr(hWnd), 0, 0)
	err = handlerErrNo(r1, r2, errNo)
	return
}

const (
	WM_CLOSE = 16
)

const (
	MB_OK = 0x00000000
)

func PostMessage(hWnd syscall.Handle, msg uint32, wParam, lParam uintptr) (err error) {
	r1, r2, errNo := syscall.Syscall6(procPostMessage.Addr(), 4,
		uintptr(hWnd),
		uintptr(msg),
		wParam,
		lParam,
		0,
		0)

	err = handlerErrNo(r1, r2, errNo)
	return
}

func SetCursorPos(x, y int32) (err error) {
	r1, r2, errNo := syscall.Syscall(procSetCursorPos.Addr(), 2, uintptr(x), uintptr(y), 0)
	err = handlerErrNo(r1, r2, errNo)
	return
}

// nInputs: The number of structures in the pInputs array.
// pInputs: expects a unsafe.Pointer to a slice of MOUSE_INPUT or KEYBD_INPUT or HARDWARE_INPUT structs.
// cbSize: The size, in bytes, of an INPUT structure. If cbSize is not the size of an INPUT structure, the function fails.
func SendInput(nInputs uint32, pInputs unsafe.Pointer, cbSize int32) (uint32, error) {
	r1, r2, errNo := syscall.Syscall(procSendInput.Addr(), 3,
		uintptr(nInputs),
		uintptr(pInputs),
		uintptr(cbSize))
	err := handlerErrNo(r1, r2, errNo)
	return uint32(r1), err
}

type MOUSE_INPUT struct {
	Type uint32
	Mi   MOUSEINPUT
}

type MOUSEINPUT struct {
	// 如果指定了 MOWSEEVENTF_ABSOLOTE 值，则 dX 和 dy 含有标准化的绝对坐标，其值在 0 到 65535 之间。事件程序将此坐标映射到显示表面。坐标（0，0）映射到显示表面的左上角，（65535，65535）映射到右下角。
	// 如果没指定 MOWSEEVENTF_ABSOLOTE，dX 和 dy 表示相对于上次鼠标事件产生的位置（即上次报告的位置）的移动。正值表示鼠标向右（或下）移动；负值表示鼠标向左（或上）移动。
	Dx          int32  // 指定鼠标沿 xWindow 轴的绝对位置或者从上次鼠标事件产生以来移动的数量
	Dy          int32  // 指定鼠标沿 yWindow 轴的绝对位置或者从上次鼠标事件产生以来移动的数量
	MouseData   uint32 // 如果 dwFlags 为 MOOSEEVENTF_WHEEL，MouseData 指定鼠标轮移动的数量。正值表示车轮向前旋转，远离用户；负值表示轮子向后旋转，朝向用户。一次车轮咔嗒声被定义为车轮_DELTA，即 120。如果 dwFlagsS 不是 MOOSEEVENTF_WHEEL，则 dWData 应为零。
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr // 鼠标的相对移动服从鼠标速度和加速度等级的设置，一个最终用户用鼠标控制面板应用程序来设置这些值，应用程序用函数 SystemParameterslnfo 来取得和设置这些值。
}

// MOUSEINPUT DwFlags
const (
	MOUSEEVENTF_ABSOLUTE        = 0x8000 // 是否采用绝对坐标, 如果不设置此位，参数含有相对数据：相对于上次位置的改动位置。
	MOUSEEVENTF_MOVE            = 0x0001 // 移动鼠标
	MOUSEEVENTF_MOVE_NOCOALESCE = 0x2000 // do not coalesce mouse moves
	MOUSEEVENTF_LEFTDOWN        = 0x0002 // 鼠标左键按下
	MOUSEEVENTF_LEFTUP          = 0x0004 // 鼠标左键抬起
	MOUSEEVENTF_RIGHTDOWN       = 0x0008 // 鼠标右键按下
	MOUSEEVENTF_RIGHTUP         = 0x0010 // 鼠标右键抬起
	MOUSEEVENTF_MIDDLEDOWN      = 0x0020 // 鼠标中键按下
	MOUSEEVENTF_MIDDLEUP        = 0x0040 // 鼠标中键抬起
	MOUSEEVENTF_VIRTUALDESK     = 0x4000 // 将坐标映射到整个桌面。必须与 MOUSEEVENTF_ABSOLUTE 一起使用。
	MOUSEEVENTF_WHEEL           = 0x0800 // 滚轮正向滚, mouseData 指定轮子的移动量。
	MOUSEEVENTF_HWHEEL          = 0x1000 // 滚轮反向滚
	MOUSEEVENTF_XDOWN           = 0x0080 // X键按下
	MOUSEEVENTF_XUP             = 0x0100 // X键抬起
)

// INPUT Type
const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2
)

const (
	LeftButton   int = 0
	RightButton  int = 1
	MiddleButton int = 2
)

// buttonType: LeftButton or RightButton or MiddleButton
func MouseClick(buttonType int, x, y int32) (uint32, error) {
	realX := 65535 * x / cxScreen // 转换后的x
	realY := 65535 * y / cyScreen // 转换后的y

	var input MOUSE_INPUT
	input.Type = INPUT_MOUSE
	input.Mi.Dx = realX
	input.Mi.Dy = realY
	if buttonType == 0 {
		input.Mi.DwFlags = MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE | MOUSEEVENTF_LEFTDOWN | MOUSEEVENTF_LEFTUP
	} else if buttonType == 1 {
		input.Mi.DwFlags = MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE | MOUSEEVENTF_RIGHTDOWN | MOUSEEVENTF_RIGHTUP
	} else if buttonType == 2 {
		input.Mi.DwFlags = MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE | MOUSEEVENTF_MIDDLEDOWN | MOUSEEVENTF_MIDDLEUP
	}
	input.Mi.MouseData = 0
	input.Mi.DwExtraInfo = 0
	input.Mi.Time = 0
	return SendInput(2, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}

func MouseDoubleClick(buttonType int, x, y int32) (ret uint32, err error) {
	if ret, err = MouseClick(buttonType, x, y); err != nil {
		return
	}
	if ret, err = MouseClick(buttonType, x, y); err != nil {
		return
	}
	return
}

type KEYBD_INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
}
type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
	Unused      [8]byte
}

// KEYBDINPUT DwFlags
const (
	KEYEVENTF_EXTENDEDKEY = 0x0001 // 若指定该值，则扫描码前一个值为 OXEO（224）的前缀字节。
	KEYEVENTF_KEYUP       = 0x0002 // 若指定该值，该键将被释放；若未指定该值，该键交被接下。
	KEYEVENTF_SCANCODE    = 0x0008
	KEYEVENTF_UNICODE     = 0x0004
)

func UniKeyPress(keyCode uint16) (uint32, error) {
	var input = [2]KEYBD_INPUT{}

	input[0].Type = INPUT_KEYBOARD
	input[0].Ki.WVk = 0
	input[0].Ki.WScan = keyCode
	input[0].Ki.DwFlags = KEYEVENTF_UNICODE

	input[1].Type = INPUT_KEYBOARD
	input[1].Ki.WVk = 0
	input[1].Ki.WScan = keyCode
	input[1].Ki.DwFlags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP

	ret, err := SendInput(2, unsafe.Pointer(&input[0]), int32(unsafe.Sizeof(input[0])))
	time.Sleep(5 * time.Millisecond)
	return ret, err
}
