package hot

import (
	"fmt"
	"testing"
)

func Test_GetItHome(t *testing.T) {
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

func Test_GetZhihu(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetZhihu()
	if err != nil {
		t.Errorf("GetZhihu failed, err: %v", err)
	}
	printData(result)
}

func Test_GetWeibo(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetWeibo()
	if err != nil {
		t.Errorf("GetWeibo failed, err: %v", err)
	}
	printData(result)
}

func Test_GetBilibili(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetBilibili()
	if err != nil {
		t.Errorf("GetBilibili failed, err: %v", err)
	}
	printData(result)
}

func Test_GetBaidu(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetBaidu()
	if err != nil {
		t.Errorf("GetBaidu failed, err: %v", err)
	}
	printData(result)
}

func Test_GetGitHub(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetGitHub()
	if err != nil {
		t.Errorf("GetGitHub failed, err: %v", err)
	}
	printData(result)
}

func Test_GetDouyin(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetDouyin()
	if err != nil {
		t.Errorf("GetDouyin failed, err: %v", err)
	}
	printData(result)
}

func Test_GetKuaishou(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetKuaishou()
	if err != nil {
		t.Errorf("GetKuaishou failed, err: %v", err)
	}
	printData(result)
}

func Test_GetToutiao(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetToutiao()
	if err != nil {
		t.Errorf("GetToutiao failed, err: %v", err)
	}
	printData(result)
}

func Test_GetCSDN(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetCSDN()
	if err != nil {
		t.Errorf("GetCSDN failed, err: %v", err)
	}
	printData(result)
}

func Test_GetTieba(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetTieba()
	if err != nil {
		t.Errorf("GetTieba failed, err: %v", err)
	}
	printData(result)
}

func Test_GetZhihuDaily(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetZhihuDaily()
	if err != nil {
		t.Errorf("GetZhihuDaily failed, err: %v", err)
	}
	printData(result)
}

func Test_GetCoolapk(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetCoolapk()
	if err != nil {
		t.Errorf("GetCoolapk failed, err: %v", err)
	}
	printData(result)
}

func Test_GetHupu(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetHupu()
	if err != nil {
		t.Errorf("GetHupu failed, err: %v", err)
	}
	printData(result)
}

func Test_GetHuxiu(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetHuxiu()
	if err != nil {
		t.Errorf("GetHuxiu failed, err: %v", err)
	}
	printData(result)
}

func Test_GetJianshu(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetJianshu()
	if err != nil {
		t.Errorf("GetJianshu failed, err: %v", err)
	}
	printData(result)
}

func Test_GetSmzdm(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetSmzdm()
	if err != nil {
		t.Errorf("GetSmzdm failed, err: %v", err)
	}
	printData(result)
}

func Test_GetSspai(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetSspai()
	if err != nil {
		t.Errorf("GetSspai failed, err: %v", err)
	}
	printData(result)
}

func Test_GetNetease(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetNetease()
	if err != nil {
		t.Errorf("GetNetease failed, err: %v", err)
	}
	printData(result)
}

func Test_GetQQ(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetQQ()
	if err != nil {
		t.Errorf("GetQQ failed, err: %v", err)
	}
	printData(result)
}

func Test_GetDoubanGroup(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetDoubanGroup()
	if err != nil {
		t.Errorf("GetDoubanGroup failed, err: %v", err)
	}
	printData(result)
}

func Test_GetAcfun(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetAcfun()
	if err != nil {
		t.Errorf("GetAcfun failed, err: %v", err)
	}
	printData(result)
}

func Test_GetDgtle(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetDgtle()
	if err != nil {
		t.Errorf("GetDgtle failed, err: %v", err)
	}
	printData(result)
}

func Test_GetDoubanMovie(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetDoubanMovie()
	if err != nil {
		t.Errorf("GetDoubanMovie failed, err: %v", err)
	}
	printData(result)
}

func Test_GetEarthquake(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetEarthquake()
	if err != nil {
		t.Errorf("GetEarthquake failed, err: %v", err)
	}
	printData(result)
}

func Test_GetGameres(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetGameres()
	if err != nil {
		t.Errorf("GetGameres failed, err: %v", err)
	}
	printData(result)
}

func Test_GetGeekpark(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetGeekpark()
	if err != nil {
		t.Errorf("GetGeekpark failed, err: %v", err)
	}
	printData(result)
}

func Test_GetGenshin(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetGenshin()
	if err != nil {
		t.Errorf("GetGenshin failed, err: %v", err)
	}
	printData(result)
}

func Test_GetGuokr(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetGuokr()
	if err != nil {
		t.Errorf("GetGuokr failed, err: %v", err)
	}
	printData(result)
}

func Test_Get36Kr(t *testing.T) {
	spider := NewSpider()
	result, err := spider.Get36Kr()
	if err != nil {
		t.Errorf("Get36kr failed, err: %v", err)
	}
	printData(result)
}

func Test_GetHackernews(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetHackernews()
	if err != nil {
		t.Errorf("GetHackernews failed, err: %v", err)
	}
	printData(result)
}

func Test_GetHelloGitHub(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetHelloGitHub()
	if err != nil {
		t.Errorf("GetHelloGitHub failed, err: %v", err)
	}
	printData(result)
}

func Test_GetHistory(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetHistory()
	if err != nil {
		t.Errorf("GetHistory failed, err: %v", err)
	}
	printData(result)
}

func Test_GetHonkai(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetHonkai()
	if err != nil {
		t.Errorf("GetHonkai failed, err: %v", err)
	}
	printData(result)
}

func Test_GetHostloc(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetHostloc()
	if err != nil {
		t.Errorf("GetHostloc failed, err: %v", err)
	}
	printData(result)
}

func Test_GetIfanr(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetIfanr()
	if err != nil {
		t.Errorf("GetIfanr failed, err: %v", err)
	}
	printData(result)
}

func Test_GetIthomeXijiayi(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetIthomeXijiayi()
	if err != nil {
		t.Errorf("GetIthomeXijiayi failed, err: %v", err)
	}
	printData(result)
}

func Test_GetMiyoushe(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetMiyoushe()
	if err != nil {
		t.Errorf("GetMiyoushe failed, err: %v", err)
	}
	printData(result)
}

func Test_GetNewsmth(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetNewsmth()
	if err != nil {
		t.Errorf("GetNewsmth failed, err: %v", err)
	}
	printData(result)
}

func Test_GetNgabbs(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetNgabbs()
	if err != nil {
		t.Errorf("GetNgabbs failed, err: %v", err)
	}
	printData(result)
}

func Test_GetNodeseek(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetNodeseek()
	if err != nil {
		t.Errorf("GetNodeseek failed, err: %v", err)
	}
	printData(result)
}

func Test_GetNytimes(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetNytimes()
	if err != nil {
		t.Errorf("GetNytimes failed, err: %v", err)
	}
	printData(result)
}

func Test_GetProducthunt(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetProducthunt()
	if err != nil {
		t.Errorf("GetProducthunt failed, err: %v", err)
	}
	printData(result)
}

func Test_GetSinaNews(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetSinaNews()
	if err != nil {
		t.Errorf("GetSinaNews failed, err: %v", err)
	}
	printData(result)
}

func Test_GetSina(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetSina()
	if err != nil {
		t.Errorf("GetSina failed, err: %v", err)
	}
	printData(result)
}

func Test_GetStarrail(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetStarrail()
	if err != nil {
		t.Errorf("GetStarrail failed, err: %v", err)
	}
	printData(result)
}

func Test_GetThepaper(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetThepaper()
	if err != nil {
		t.Errorf("GetThepaper failed, err: %v", err)
	}
	printData(result)
}

func Test_GetWeatheralarm(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetWeatheralarm()
	if err != nil {
		t.Errorf("GetWeatheralarm failed, err: %v", err)
	}
	printData(result)
}

func Test_GetWeread(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetWeread()
	if err != nil {
		t.Errorf("GetWeread failed, err: %v", err)
	}
	printData(result)
}

func Test_GetYystv(t *testing.T) {
	spider := NewSpider()
	result, err := spider.GetYystv()
	if err != nil {
		t.Errorf("GetYystv failed, err: %v", err)
	}
	printData(result)
}

func printData(result *HotData) {
	if result == nil || len(result.Data) == 0 {
		fmt.Println("result is empty")
		return
	}

	for _, v := range result.Data {
		fmt.Printf("id:%s,title:%s,hot:%d,url: %s\n", v.ID, v.Title, v.Hot, v.URL)
	}
}
