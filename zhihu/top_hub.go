package zhihu

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/tuoaitang/spider/log"
	"github.com/tuoaitang/spider/model"
)

func GetTopHub() (*model.ZhiHuTopHub, error) {
	url := "https://www.zhihu.com/billboard"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("parse err:%v", err)
		return nil, err
	}

	js := doc.Find("#js-initialData").Text()
	//log.Info("js:%s", js)

	newJson, err := simplejson.NewJson([]byte(js))
	if err != nil {
		log.Error("new json err:%v", err)
		return nil, err
	}

	hotList := newJson.Get("initialState").Get("topstory").Get("hotList")
	hotListArray := hotList.MustArray()

	zth := &model.ZhiHuTopHub{}
	for i := 0; i < len(hotListArray); i++ {
		l := hotList.GetIndex(i)
		title := l.Get("target").Get("titleArea").Get("text").MustString()
		abstract := l.Get("target").Get("excerptArea").Get("text").MustString()
		link := l.Get("target").Get("link").Get("url").MustString()
		tp := l.Get("type").MustString()
		answerCount := l.Get("feedSpecific").Get("answerCount").MustInt()
		zth.List = append(zth.List, &model.ZhiHuTopHubItem{
			Title:       title,
			URL:         link,
			Abstract:    abstract,
			AnswerCount: answerCount,
			Type:        tp,
		})
	}

	return zth, nil
}
