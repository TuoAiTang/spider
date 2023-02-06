package model

import (
	"fmt"
)

type ZhiHuTopHub struct {
	List []*ZhiHuTopHubItem `json:"list"`
}

type ZhiHuTopHubItem struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Abstract    string `json:"abstract"`
	Type        string `json:"type"`
	AnswerCount int    `json:"answer_count"`
}

func (zth *ZhiHuTopHub) String() string {
	s := "【知乎热榜】\n"
	for i := 0; i < len(zth.List); i++ {
		v := zth.List[i]
		s += fmt.Sprintf("%d. 🔥%d %s[%s]\n", i+1, v.AnswerCount, v.Title, v.URL)
	}
	return s
}

func (zth *ZhiHuTopHub) Top(i int) *ZhiHuTopHub {
	if i <= 0 {
		i = 10
	}

	if i > len(zth.List) {
		i = len(zth.List)
	}
	zth.List = zth.List[:i]
	return zth
}
