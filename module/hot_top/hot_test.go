package hottop

import (
	"fmt"
	"testing"
)

func Test_GetList(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetItHome()
	if err != nil {
		t.Errorf("GetList failed, err: %v", err)
	}
	printData(result)
	result, err = spider.GetJuejin()
	if err != nil {
		t.Errorf("GetList failed, err: %v", err)
	}
	printData(result)

}

func printData(result *HotData) {
	for _, v := range result.Data {
		fmt.Printf("id:%d,title:%s,hot:%d,\n", v.ID, v.Title, v.Hot)
	}
}
