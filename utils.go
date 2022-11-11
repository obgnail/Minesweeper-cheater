package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			DisableQuote:              true,
			DisableLevelTruncation:    true,
			FullTimestamp:             true,
			TimestampFormat:           "15:04:05",
		},
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//组合算法(从n中取出m个数)
func Combination(n int, m int) [][]int {
	if m < 1 || m > n {
		return [][]int{}
	}

	// 保存最终结果的数组，总数直接通过数学公式计算
	result := make([][]int, 0, combinationNum(n, m))
	// 保存每一个组合的索引的数组，1表示选中，0表示未选中
	indexes := make([]int, n)
	for i := 0; i < n; i++ {
		if i < m {
			indexes[i] = 1
		} else {
			indexes[i] = 0
		}
	}

	//第一个结果
	result = addTo(result, indexes)
	for {
		find := false
		//每次循环将第一次出现的 1 0 改为 0 1，同时将左侧的1移动到最左侧
		for i := 0; i < n-1; i++ {
			if indexes[i] == 1 && indexes[i+1] == 0 {
				find = true

				indexes[i], indexes[i+1] = 0, 1
				if i > 1 {
					moveOneToLeft(indexes[:i])
				}
				result = addTo(result, indexes)

				break
			}
		}

		//本次循环没有找到 1 0 ，说明已经取到了最后一种情况
		if !find {
			break
		}
	}

	return result
}

// 将ele复制后添加到arr中，返回新的数组
func addTo(arr [][]int, ele []int) [][]int {
	newEle := make([]int, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)
	return arr
}

// 计算组合数(从n中取m个数)
func combinationNum(n int, m int) int {
	return factorial(n) / (factorial(n-m) * factorial(m))
}

// 阶乘
func factorial(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func moveOneToLeft(leftNums []int) {
	//计算有几个1
	sum := 0
	for i := 0; i < len(leftNums); i++ {
		if leftNums[i] == 1 {
			sum++
		}
	}

	// 将前sum个改为1，之后的改为0
	for i := 0; i < len(leftNums); i++ {
		if i < sum {
			leftNums[i] = 1
		} else {
			leftNums[i] = 0
		}
	}
}
