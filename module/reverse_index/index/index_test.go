package index

import (
	"fmt"
	"testing"

	"github.com/fumiama/jieba"
)

func Test_SplitWords(t *testing.T) {
	seg, err := jieba.LoadDictionaryAt("dict.txt")
	if err != nil {
		panic(err)
	}

	fmt.Print("【全模式】：")
	fmt.Println(seg.CutAll("我来到北京清华大学"))

	fmt.Print("【精确模式】：")
	fmt.Println(seg.Cut("我来到北京清华大学", false))

	fmt.Print("【新词识别】：")
	fmt.Println(seg.Cut("他来到了网易杭研大厦", true))

	fmt.Print("【搜索引擎模式】：")
	fmt.Println(seg.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))

	fmt.Print("新词：")
	fmt.Println(seg.Cut("我喜欢人工智能AI", true))
}
