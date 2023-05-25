package main

import (
	"github.com/kbinani/screenshot"
	"github.com/rivo/duplo"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path"
	"sort"
	"time"
)

const (
	startX = 86
	startY = 149
	endX   = 1416
	endY   = 858

	BoxHeight     = 44
	halfBoxHeight = BoxHeight / 2
)

var store *duplo.Store

var matchMap = map[string]CellType{
	"0-1.jpg": CellTypeZero,
	"1-1.jpg": CellTypeOne,
	"2-1.jpg": CellTypeTwo,
	"3-1.jpg": CellTypeThree,
	"4-1.jpg": CellTypeFour,
	"5-1.jpg": CellTypeFive,
	"6-1.jpg": CellTypeSix,
	"7-1.jpg": CellTypeSeven,

	"0-2.jpg": CellTypeZero,
	"1-2.jpg": CellTypeOne,
	"2-2.jpg": CellTypeTwo,
	"3-2.jpg": CellTypeThree,
	"4-2.jpg": CellTypeFour,
	"5-2.jpg": CellTypeFive,
	"6-2.jpg": CellTypeSix,
	"7-2.jpg": CellTypeSeven,

	"unknown-1.jpg": CellTypeUnknown,
	"unknown-2.jpg": CellTypeUnknown,
	"unknown-3.jpg": CellTypeUnknown,
	"unknown-4.jpg": CellTypeUnknown,
}

func InitConfig() {
	if !showFlag || clickInterval < 200*time.Microsecond {
		useDoubleClick = false
	}

	if clickInterval < 100*time.Millisecond {
		clickInterval = 100 * time.Millisecond
	}

	if showFlag {
		matchMap["flag-1.jpg"] = CellTypeFlag
		matchMap["flag-2.jpg"] = CellTypeFlag
		matchMap["flag-3.jpg"] = CellTypeFlag
		matchMap["flag-4.jpg"] = CellTypeFlag
	}
}

func InitStore() {
	store = duplo.New()
	for key := range matchMap {
		p := path.Join("image", key)
		img, err := OpenImage(p)
		if err != nil {
			Logger.Fatal(err)
		}
		hash, _ := duplo.CreateHash(img)
		store.Add(key, hash)
	}
}

func getNumFromImage(img image.Image) CellType {
	hash, _ := duplo.CreateHash(img)
	matches := store.Query(hash)
	sort.Sort(matches)
	key := matches[0].ID.(string)
	return matchMap[key]
}

func OpenImage(imagePath string) (image.Image, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	m, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Window2Table(skipFunc func(row, col int) bool) [rowNum][colNum]CellType {
	Logger.Debug("Parsing Window...")
	img, err := screenshot.Capture(startX, startY, endX-startX, endY-startY) // TODO 硬编码,在不同显示器不同分辨率下八成不行
	if err != nil {
		Logger.Fatal(err)
	}
	res := ParseImageToTable(img, skipFunc)
	Logger.Debug("Start Analyzing...")
	return res
}

func ParseImageToTable(img *image.RGBA, skipFunc func(row, col int) bool) [rowNum][colNum]CellType {
	res := [rowNum][colNum]CellType{}
	for row := 0; row != rowNum; row++ {
		boxY := row * (endY - startY) / rowNum
		nextBoxY := (row + 1) * (endY - startY) / rowNum
		for col := 0; col != colNum; col++ {

			// 跳过某些单元格,不去解析
			if skipFunc != nil && skipFunc(row, col) {
				continue
			}

			boxX := col * (endX - startX) / colNum
			nextBoxX := (col + 1) * (endX - startX) / colNum
			boxX, boxY, nextBoxX, nextBoxY = fixLen(boxX, boxY, nextBoxX, nextBoxY)
			// 多裁剪10个像素, 减少边框的影响
			newImg := GetSubImage(img, boxX+5, boxY+5, nextBoxX-5, nextBoxY-5)
			ret := getNumFromImage(newImg)
			res[row][col] = ret
			//name := fmt.Sprintf("image/%d_(%d,%d).jpg", ret, row, col)
			//SaveImage(name, newImg)
		}
	}
	return res
}

func fixLen(boxX, boxY, nextBoxX, nextBoxY int) (int, int, int, int) {
	if nextBoxX-boxX > BoxHeight {
		boxX = nextBoxX - BoxHeight
	}
	if nextBoxY-boxY > BoxHeight {
		boxY = nextBoxY - BoxHeight
	}
	return boxX, boxY, nextBoxX, nextBoxY
}

func GetSubImage(img *image.RGBA, x0, y0, x1, y1 int) *image.RGBA {
	subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
	newImg := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))
	for y := 0; y < y1-y0; y++ {
		for x := 0; x < x1-x0; x++ {
			r32, g32, b32, a32 := subImg.At(x0+x, y0+y).RGBA()
			r := uint8(r32 >> 8)
			g := uint8(g32 >> 8)
			b := uint8(b32 >> 8)
			a := uint8(a32 >> 8)
			newImg.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}
	return newImg
}

func SaveImage(name string, img *image.RGBA) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		return err
	}
	return nil
}
