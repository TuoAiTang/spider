package leetcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/tuoaitang/spider/log"
	"github.com/tuoaitang/spider/model"
)

func GetToday() (*model.CodeProblem, error) {
	slug, err := getTodayTitleSlug()
	if err != nil {
		log.Error("getTodayTitleSlug err:%v", err)
		return nil, err
	}

	q, err := getProblemByTitleSlug(slug)
	if err != nil {
		log.Error("getProblemByTitleSlug err:%v", err)
		return nil, err
	}

	return q, nil
}

// getProblemByTitleSlug 根据titleSlug获取题目
func getProblemByTitleSlug(titleSlug string) (*model.CodeProblem, error) {
	url := "https://leetcode.cn/graphql/"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"operationName":"questionData","variables":{"titleSlug":"%s"},"query":"query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    questionFrontendId\n    categoryTitle\n    boundTopicId\n    title\n    titleSlug\n    content\n    translatedTitle\n    translatedContent\n    isPaidOnly\n    difficulty\n    likes\n    dislikes\n    isLiked\n    similarQuestions\n    contributors {\n      username\n      profileUrl\n      avatarUrl\n      __typename\n    }\n    langToValidPlayground\n    topicTags {\n      name\n      slug\n      translatedName\n      __typename\n    }\n    companyTagStats\n    codeSnippets {\n      lang\n      langSlug\n      code\n      __typename\n    }\n    stats\n    hints\n    solution {\n      id\n      canSeeDetail\n      __typename\n    }\n    status\n    sampleTestCase\n    metaData\n    judgerAvailable\n    judgeType\n    mysqlSchemas\n    enableRunCode\n    envInfo\n    book {\n      id\n      bookName\n      pressName\n      source\n      shortDescription\n      fullDescription\n      bookImgUrl\n      pressImgUrl\n      productUrl\n      __typename\n    }\n    isSubscribed\n    isDailyQuestion\n    dailyRecordStatus\n    editorType\n    ugcQuestionId\n    style\n    exampleTestcases\n    jsonExampleTestcases\n    __typename\n  }\n}\n"}`, titleSlug))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("authority", "leetcode.cn")
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "zh-CN")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cookie", "csrftoken=xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW; gr_user_id=0af8e097-9603-4b68-b866-8fbe46f88da0; _bl_uid=6Clhyc6e67F0F5hI4iLClIn2jsX7; aliyungf_tc=d7c758b5778e17dd35ecef145d11b630520f2134e337132184bca1f117b1a4be; Hm_lvt_f0faad39bcf8471e3ab3ef70125152c3=1672132431,1673001602; a2873925c34ecbd2_gr_session_id=52f336bf-55c7-40bb-b567-48e042c8d0c0; a2873925c34ecbd2_gr_session_id_52f336bf-55c7-40bb-b567-48e042c8d0c0=true; _gid=GA1.2.377141455.1673001603; __appToken__=; NEW_PROBLEMLIST_PAGE=1; Hm_lpvt_f0faad39bcf8471e3ab3ef70125152c3=1673001795; _gat_gtag_UA_131851415_1=1; _ga_PDVPZYN3CW=GS1.1.1673001606.2.1.1673001795.0.0.0; _ga=GA1.1.1772867861.1672132432; csrftoken=xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW")
	req.Header.Add("origin", "https://leetcode.cn")
	req.Header.Add("random-uuid", "e920c3bd-cd4e-bd4d-0697-30d817817145")
	req.Header.Add("referer", "https://leetcode.cn/problems/count-integers-with-even-digit-sum/")
	req.Header.Add("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Add("x-csrftoken", "xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW")
	req.Header.Add("x-definition-name", "question")
	req.Header.Add("x-operation-name", "questionData")
	req.Header.Add("x-timezone", "Asia/Shanghai")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sj, err := simplejson.NewJson(body)
	if err != nil {
		log.Error("new json err:%v", err)
		return nil, err
	}

	q := sj.Get("data").Get("question")
	cp := &model.CodeProblem{
		Abstract:  htmlToTxt(q.Get("translatedContent").MustString()),
		ID:        q.Get("questionId").MustString(),
		Index:     q.Get("questionFrontendId").MustString(),
		Title:     q.Get("translatedTitle").MustString(),
		TitleSlug: titleSlug,
		Level:     q.Get("difficulty").MustString(),
	}

	arr := q.Get("topicTags").MustArray()
	for i := 0; i < len(arr); i++ {
		cp.Tags = append(cp.Tags, q.Get("topicTags").GetIndex(i).Get("translatedName").MustString())
	}

	arr = q.Get("hints").MustArray()
	for i := 0; i < len(arr); i++ {
		cp.Hints = append(cp.Hints, q.Get("hints").GetIndex(i).MustString())
	}

	return cp, nil
}

// getTodayTitleSlug 获取今日题目的titleSlug
func getTodayTitleSlug() (string, error) {
	slug, err := do(func() (*http.Client, *http.Request, error) {
		url := "https://leetcode.cn/graphql/"
		method := "POST"

		payload := strings.NewReader(`{"query":"\n    query questionOfToday {\n  todayRecord {\n    date\n    userStatus\n    question {\n      questionId\n      frontendQuestionId: questionFrontendId\n      difficulty\n      title\n      titleCn: translatedTitle\n      titleSlug\n      paidOnly: isPaidOnly\n      freqBar\n      isFavor\n      acRate\n      status\n      solutionNum\n      hasVideoSolution\n      topicTags {\n        name\n        nameTranslated: translatedName\n        id\n      }\n      extra {\n        topCompanyTags {\n          imgUrl\n          slug\n          numSubscribed\n        }\n      }\n    }\n    lastSubmission {\n      id\n    }\n  }\n}\n    ","variables":{}}`)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			return nil, nil, err
		}

		req.Header.Add("authority", "leetcode.cn")
		req.Header.Add("accept", "*/*")
		req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
		req.Header.Add("authorization", "")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("cookie", "csrftoken=xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW; gr_user_id=0af8e097-9603-4b68-b866-8fbe46f88da0; _bl_uid=6Clhyc6e67F0F5hI4iLClIn2jsX7; aliyungf_tc=d7c758b5778e17dd35ecef145d11b630520f2134e337132184bca1f117b1a4be; Hm_lvt_f0faad39bcf8471e3ab3ef70125152c3=1672132431,1673001602; a2873925c34ecbd2_gr_session_id=52f336bf-55c7-40bb-b567-48e042c8d0c0; a2873925c34ecbd2_gr_session_id_52f336bf-55c7-40bb-b567-48e042c8d0c0=true; _gid=GA1.2.377141455.1673001603; __appToken__=; NEW_PROBLEMLIST_PAGE=1; _ga_PDVPZYN3CW=GS1.1.1673001606.2.0.1673001611.0.0.0; Hm_lpvt_f0faad39bcf8471e3ab3ef70125152c3=1673001612; _ga=GA1.2.1772867861.1672132432; csrftoken=xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW")
		req.Header.Add("origin", "https://leetcode.cn")
		req.Header.Add("random-uuid", "e920c3bd-cd4e-bd4d-0697-30d817817145")
		req.Header.Add("referer", "https://leetcode.cn/problemset/all/")
		req.Header.Add("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		req.Header.Add("x-csrftoken", "xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW")

		return client, req, nil
	}, func(sj *simplejson.Json) (interface{}, error) {
		return sj.Get("data").Get("todayRecord").GetIndex(0).Get("question").Get("titleSlug").MustString(), nil
	})

	if err != nil {
		log.Error("err:%v", err)
		return "", err
	}

	return slug.(string), err
}

func do(buildRequestAndClientFunc func() (*http.Client, *http.Request, error), parseResponseFunc func(sj *simplejson.Json) (interface{}, error)) (interface{}, error) {
	cli, req, err := buildRequestAndClientFunc()
	if err != nil {
		log.Error("build request and client err:%v", err)
		return "", err
	}
	res, err := cli.Do(req)
	if err != nil {
		log.Error("do err:%v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("read body err:%v", err)
		return nil, err
	}

	sj, err := simplejson.NewJson(body)
	if err != nil {
		log.Error("new json err:%v", err)
		return "", err
	}

	return parseResponseFunc(sj)
}

// htmlToTxt html转txt
func htmlToTxt(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error("new document err:%v", err)
		return ""
	}

	txt := doc.Text()

	return txt
}

func getSubmissionQuestion() ([]*model.SubmissionQuestion, error) {
	resp, err := do(func() (*http.Client, *http.Request, error) {
		url := "https://leetcode.cn/graphql/"
		method := "POST"

		payload := strings.NewReader(`{"operationName":"userProfileQuestions","variables":{"status":"ACCEPTED","skip":0,"first":300,"sortField":"LAST_SUBMITTED_AT","sortOrder":"DESCENDING","difficulty":[]},"query":"query userProfileQuestions($status: StatusFilterEnum!, $skip: Int!, $first: Int!, $sortField: SortFieldEnum!, $sortOrder: SortingOrderEnum!, $keyword: String, $difficulty: [DifficultyEnum!]) {\n  userProfileQuestions(status: $status, skip: $skip, first: $first, sortField: $sortField, sortOrder: $sortOrder, keyword: $keyword, difficulty: $difficulty) {\n    totalNum\n    questions {\n      translatedTitle\n      frontendId\n      titleSlug\n      title\n      difficulty\n      lastSubmittedAt\n      numSubmitted\n      lastSubmissionSrc {\n        sourceType\n        ... on SubmissionSrcLeetbookNode {\n          slug\n          title\n          pageId\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"}`)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			log.Error("new request err:%v", err)
			return nil, nil, err
		}
		req.Header.Add("authority", "leetcode.cn")
		req.Header.Add("accept", "*/*")
		req.Header.Add("accept-language", "zh-CN")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("cookie", "gr_user_id=0af8e097-9603-4b68-b866-8fbe46f88da0; _bl_uid=6Clhyc6e67F0F5hI4iLClIn2jsX7; aliyungf_tc=3195e356fb5bb20517395f10bf548d0ff8a701d58ab884368016f278d97780b5; Hm_lvt_f0faad39bcf8471e3ab3ef70125152c3=1673001602,1673402067,1675158408; _gid=GA1.2.1785451212.1675158409; a2873925c34ecbd2_gr_session_id=b2e1f89b-eed0-4f80-b366-f4840636641d; a2873925c34ecbd2_gr_session_id_b2e1f89b-eed0-4f80-b366-f4840636641d=true; __appToken__=; csrftoken=55ytv1b8CzeCXuab1qrKbzcwawvtDR9tQ8QgdTKQ8Vj7knjB23jmOPKCIWvFQgZ0; a2873925c34ecbd2_gr_last_sent_sid_with_cs1=b2e1f89b-eed0-4f80-b366-f4840636641d; a2873925c34ecbd2_gr_last_sent_cs1=0xcafebabe__; NEW_QUESTION_DETAIL_PAGE_V2=1; Hm_lpvt_f0faad39bcf8471e3ab3ef70125152c3=1675232834; _ga_PDVPZYN3CW=GS1.1.1675232684.6.1.1675232833.0.0.0; _ga=GA1.2.1772867861.1672132432; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMjIwNzc2IiwiX2F1dGhfdXNlcl9iYWNrZW5kIjoiZGphbmdvLmNvbnRyaWIuYXV0aC5iYWNrZW5kcy5Nb2RlbEJhY2tlbmQiLCJfYXV0aF91c2VyX2hhc2giOiI3YzM0NzhiNjUyN2U0ZmY0ZWVhYTVhNjA0ZTI0YzYyZTM0ZmQ2YjAxMjQ4MGJmYzFiYjY2NjU4OTU3YmE2NmVjIiwiaWQiOjIyMDc3NiwiZW1haWwiOiIiLCJ1c2VybmFtZSI6IjB4Q0FGRUJBQkVfXyIsInVzZXJfc2x1ZyI6IjB4Y2FmZWJhYmVfXyIsImF2YXRhciI6Imh0dHBzOi8vYXNzZXRzLmxlZXRjb2RlLmNuL2FsaXl1bi1sYy11cGxvYWQvdXNlcnMvdHVvYWl0YW5nLWpaN1ZRMVlFRDgvYXZhdGFyXzE1Mzc4MDMyMjcucG5nIiwicGhvbmVfdmVyaWZpZWQiOnRydWUsIl90aW1lc3RhbXAiOjE2NzUyMzI3MTkuNzU5OTcyNiwiZXhwaXJlZF90aW1lXyI6MTY3Nzc4MzYwMCwidmVyc2lvbl9rZXlfIjowfQ.MoW-ewuV_5Ib42YJlB5g1eXjHpy2t4wNsCbaIx0mHZ4; a2873925c34ecbd2_gr_cs1=0xcafebabe__; csrftoken=xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW")
		req.Header.Add("origin", "https://leetcode.cn")
		req.Header.Add("random-uuid", "e920c3bd-cd4e-bd4d-0697-30d817817145")
		req.Header.Add("referer", "https://leetcode.cn/progress/")
		req.Header.Add("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		req.Header.Add("x-csrftoken", "55ytv1b8CzeCXuab1qrKbzcwawvtDR9tQ8QgdTKQ8Vj7knjB23jmOPKCIWvFQgZ0")
		req.Header.Add("x-definition-name", "userProfileQuestions")
		req.Header.Add("x-operation-name", "userProfileQuestions")
		return client, req, nil
	}, func(sj *simplejson.Json) (interface{}, error) {
		var questions []*model.SubmissionQuestion
		bytes, err := sj.Get("data").Get("userProfileQuestions").Get("questions").MarshalJSON()
		if err != nil {
			log.Error("get questions err:%v", err)
			return nil, err
		}

		err = json.Unmarshal(bytes, &questions)
		if err != nil {
			log.Error("unmarshal questions err:%v", err)
			return nil, err
		}

		return questions, nil
	})

	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}

	return resp.([]*model.SubmissionQuestion), err
}

// getSubmissionBySlug get submission by slug
func getSubmissionBySlug(slug string) ([]*model.SubmissionItem, error) {
	resp, err := do(func() (*http.Client, *http.Request, error) {
		url := "https://leetcode.cn/graphql/"
		method := "POST"

		payload := strings.NewReader(fmt.Sprintf(`{"query":"\n    query submissionList($offset: Int!, $limit: Int!, $lastKey: String, $questionSlug: String!, $lang: String, $status: SubmissionStatusEnum) {\n  submissionList(\n    offset: $offset\n    limit: $limit\n    lastKey: $lastKey\n    questionSlug: $questionSlug\n    lang: $lang\n    status: $status\n  ) {\n    lastKey\n    hasNext\n    submissions {\n      id\n      title\n      status\n      statusDisplay\n      lang\n      langName: langVerboseName\n      runtime\n      timestamp\n      url\n      isPending\n      memory\n      submissionComment {\n        comment\n      }\n    }\n  }\n}\n    ","variables":{"questionSlug":"%s","offset":0,"limit":20,"lastKey":null,"status":null}}`, slug))

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
		req.Header.Add("authority", "leetcode.cn")
		req.Header.Add("accept", "*/*")
		req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
		req.Header.Add("authorization", "")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("cookie", "gr_user_id=0af8e097-9603-4b68-b866-8fbe46f88da0; _bl_uid=6Clhyc6e67F0F5hI4iLClIn2jsX7; aliyungf_tc=3195e356fb5bb20517395f10bf548d0ff8a701d58ab884368016f278d97780b5; Hm_lvt_f0faad39bcf8471e3ab3ef70125152c3=1673001602,1673402067,1675158408; _gid=GA1.2.1785451212.1675158409; a2873925c34ecbd2_gr_session_id=b2e1f89b-eed0-4f80-b366-f4840636641d; a2873925c34ecbd2_gr_session_id_b2e1f89b-eed0-4f80-b366-f4840636641d=true; __appToken__=; csrftoken=55ytv1b8CzeCXuab1qrKbzcwawvtDR9tQ8QgdTKQ8Vj7knjB23jmOPKCIWvFQgZ0; a2873925c34ecbd2_gr_last_sent_sid_with_cs1=b2e1f89b-eed0-4f80-b366-f4840636641d; a2873925c34ecbd2_gr_last_sent_cs1=0xcafebabe__; NEW_QUESTION_DETAIL_PAGE_V2=1; _ga=GA1.2.1772867861.1672132432; _ga_PDVPZYN3CW=GS1.1.1675232684.6.1.1675233476.0.0.0; Hm_lpvt_f0faad39bcf8471e3ab3ef70125152c3=1675233479; a2873925c34ecbd2_gr_cs1=0xcafebabe__; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMjIwNzc2IiwiX2F1dGhfdXNlcl9iYWNrZW5kIjoiZGphbmdvLmNvbnRyaWIuYXV0aC5iYWNrZW5kcy5Nb2RlbEJhY2tlbmQiLCJfYXV0aF91c2VyX2hhc2giOiI3YzM0NzhiNjUyN2U0ZmY0ZWVhYTVhNjA0ZTI0YzYyZTM0ZmQ2YjAxMjQ4MGJmYzFiYjY2NjU4OTU3YmE2NmVjIiwiaWQiOjIyMDc3NiwiZW1haWwiOiIiLCJ1c2VybmFtZSI6IjB4Q0FGRUJBQkVfXyIsInVzZXJfc2x1ZyI6IjB4Y2FmZWJhYmVfXyIsImF2YXRhciI6Imh0dHBzOi8vYXNzZXRzLmxlZXRjb2RlLmNuL2FsaXl1bi1sYy11cGxvYWQvdXNlcnMvdHVvYWl0YW5nLWpaN1ZRMVlFRDgvYXZhdGFyXzE1Mzc4MDMyMjcucG5nIiwicGhvbmVfdmVyaWZpZWQiOnRydWUsIl90aW1lc3RhbXAiOjE2NzUyMzI3MTkuNzU5OTcyNiwiZXhwaXJlZF90aW1lXyI6MTY3Nzc4MzYwMCwidmVyc2lvbl9rZXlfIjowLCJsYXRlc3RfdGltZXN0YW1wXyI6MTY3NTIzMzU0OH0.YiHOabG_ITb98AJ4AF7fmaEBGi-x91RKzYVbPf8XE10; csrftoken=xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW")
		req.Header.Add("origin", "https://leetcode.cn")
		req.Header.Add("random-uuid", "e920c3bd-cd4e-bd4d-0697-30d817817145")
		req.Header.Add("referer", "https://leetcode.cn/problems/regular-expression-matching/submissions/")
		req.Header.Add("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		req.Header.Add("x-csrftoken", "55ytv1b8CzeCXuab1qrKbzcwawvtDR9tQ8QgdTKQ8Vj7knjB23jmOPKCIWvFQgZ0")

		return client, req, nil
	}, func(sj *simplejson.Json) (interface{}, error) {
		var submissions []*model.SubmissionItem
		bytes, err := sj.Get("data").Get("submissionList").Get("submissions").MarshalJSON()
		if err != nil {
			log.Error("get submissions err:%v", err)
			return nil, err
		}

		err = json.Unmarshal(bytes, &submissions)
		if err != nil {
			log.Error("unmarshal submissions err:%v", err)
			return nil, err
		}

		return submissions, nil
	})

	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}

	return resp.([]*model.SubmissionItem), err
}

// getSubmissionDetail
func getSubmissionCode(subID string) (string, error) {
	resp, err := do(func() (*http.Client, *http.Request, error) {
		url := "https://leetcode.cn/graphql/"
		method := "POST"

		payload := strings.NewReader(fmt.Sprintf(`{"operationName":"mySubmissionDetail","variables":{"id":"%s"},"query":"query mySubmissionDetail($id: ID!) {\n  submissionDetail(submissionId: $id) {\n    id\n    code\n    runtime\n    memory\n    rawMemory\n    statusDisplay\n    timestamp\n    lang\n    isMine\n    passedTestCaseCnt\n    totalTestCaseCnt\n    sourceUrl\n    question {\n      titleSlug\n      title\n      translatedTitle\n      questionId\n      __typename\n    }\n    ... on GeneralSubmissionNode {\n      outputDetail {\n        codeOutput\n        expectedOutput\n        input\n        compileError\n        runtimeError\n        lastTestcase\n        __typename\n      }\n      __typename\n    }\n    submissionComment {\n      comment\n      flagType\n      __typename\n    }\n    __typename\n  }\n}\n"}`, subID))

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			log.Error("new request err:%v", err)
			return nil, nil, err
		}
		req.Header.Add("authority", "leetcode.cn")
		req.Header.Add("accept", "*/*")
		req.Header.Add("accept-language", "zh-CN")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("cookie", "gr_user_id=0af8e097-9603-4b68-b866-8fbe46f88da0; _bl_uid=6Clhyc6e67F0F5hI4iLClIn2jsX7; aliyungf_tc=3195e356fb5bb20517395f10bf548d0ff8a701d58ab884368016f278d97780b5; Hm_lvt_f0faad39bcf8471e3ab3ef70125152c3=1673001602,1673402067,1675158408; _gid=GA1.2.1785451212.1675158409; a2873925c34ecbd2_gr_session_id=b2e1f89b-eed0-4f80-b366-f4840636641d; a2873925c34ecbd2_gr_session_id_b2e1f89b-eed0-4f80-b366-f4840636641d=true; __appToken__=; csrftoken=55ytv1b8CzeCXuab1qrKbzcwawvtDR9tQ8QgdTKQ8Vj7knjB23jmOPKCIWvFQgZ0; a2873925c34ecbd2_gr_last_sent_sid_with_cs1=b2e1f89b-eed0-4f80-b366-f4840636641d; a2873925c34ecbd2_gr_last_sent_cs1=0xcafebabe__; NEW_QUESTION_DETAIL_PAGE_V2=1; _gat_gtag_UA_131851415_1=1; _ga=GA1.1.1772867861.1672132432; Hm_lpvt_f0faad39bcf8471e3ab3ef70125152c3=1675233579; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMjIwNzc2IiwiX2F1dGhfdXNlcl9iYWNrZW5kIjoiZGphbmdvLmNvbnRyaWIuYXV0aC5iYWNrZW5kcy5Nb2RlbEJhY2tlbmQiLCJfYXV0aF91c2VyX2hhc2giOiI3YzM0NzhiNjUyN2U0ZmY0ZWVhYTVhNjA0ZTI0YzYyZTM0ZmQ2YjAxMjQ4MGJmYzFiYjY2NjU4OTU3YmE2NmVjIiwiaWQiOjIyMDc3NiwiZW1haWwiOiIiLCJ1c2VybmFtZSI6IjB4Q0FGRUJBQkVfXyIsInVzZXJfc2x1ZyI6IjB4Y2FmZWJhYmVfXyIsImF2YXRhciI6Imh0dHBzOi8vYXNzZXRzLmxlZXRjb2RlLmNuL2FsaXl1bi1sYy11cGxvYWQvdXNlcnMvdHVvYWl0YW5nLWpaN1ZRMVlFRDgvYXZhdGFyXzE1Mzc4MDMyMjcucG5nIiwicGhvbmVfdmVyaWZpZWQiOnRydWUsIl90aW1lc3RhbXAiOjE2NzUyMzI3MTkuNzU5OTcyNiwiZXhwaXJlZF90aW1lXyI6MTY3Nzc4MzYwMCwidmVyc2lvbl9rZXlfIjowLCJsYXRlc3RfdGltZXN0YW1wXyI6MTY3NTIzMzU0OH0.YiHOabG_ITb98AJ4AF7fmaEBGi-x91RKzYVbPf8XE10; a2873925c34ecbd2_gr_cs1=0xcafebabe__; _ga_PDVPZYN3CW=GS1.1.1675232684.6.1.1675233605.0.0.0; csrftoken=xbbTltUILZnh84cdAbKIBdW2ClL4Us4VOahu01xXO3x6JA12NGMeTTNfgjUKcvJW")
		req.Header.Add("origin", "https://leetcode.cn")
		req.Header.Add("random-uuid", "e920c3bd-cd4e-bd4d-0697-30d817817145")
		req.Header.Add("referer", "https://leetcode.cn/submissions/detail/322647831/")
		req.Header.Add("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\"")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		req.Header.Add("x-csrftoken", "55ytv1b8CzeCXuab1qrKbzcwawvtDR9tQ8QgdTKQ8Vj7knjB23jmOPKCIWvFQgZ0")
		req.Header.Add("x-definition-name", "submissionDetail")
		req.Header.Add("x-operation-name", "mySubmissionDetail")
		req.Header.Add("x-timezone", "Asia/Shanghai")
		return client, req, nil
	}, func(sj *simplejson.Json) (interface{}, error) {
		s, err := sj.Get("data").Get("submissionDetail").Get("code").String()
		if err != nil {
			mj, _ := sj.MarshalJSON()
			log.Error("get code err:%v, resp:%v", err, string(mj))
			return "", err
		}

		return s, nil
	})

	if err != nil {
		log.Error("err:%v", err)
		return "", err
	}

	return resp.(string), err
}

func GetQuestionAndSubmissionBySlug(slug string) (*model.QuestionAndSubmission, error) {
	question, err := getProblemByTitleSlug(slug)
	if err != nil {
		log.Error("getProblemByTitleSlug err:%v", err)
		return nil, err
	}

	submissions, err := getSubmissionBySlug(slug)
	if err != nil {
		log.Error("getSubmissionBySlug err:%v", err)
		return nil, err
	}

	if len(submissions) == 0 {
		return nil, errors.New("no submission")
	}

	// 只选择成功AC的一个提交
	var neededSubmissions []*model.Submission
	for _, submission := range submissions {
		if submission.StatusDisplay != "Accepted" {
			continue
		}

		neededSubmissions = append(neededSubmissions, &model.Submission{
			Submission: submission,
			Code:       "",
		})
	}

	for _, submission := range neededSubmissions {
		code, err := getSubmissionCode(submission.Submission.Id)
		if err != nil {
			log.Error("getSubmissionCode err:%v", err)
			return nil, err
		}

		submission.Code = code
	}

	return &model.QuestionAndSubmission{
		Question:    question,
		Submissions: neededSubmissions,
	}, nil
}
