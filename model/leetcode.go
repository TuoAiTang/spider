package model

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tuoaitang/spider/log"
)

type CodeProblem struct {
	ID        string
	Index     string
	Title     string
	Tags      []string
	Level     string
	Abstract  string
	Hints     []string
	TitleSlug string
}

func (cp *CodeProblem) String() string {
	s := fmt.Sprintf("「每日一题 %s - %s」", cp.Index, cp.Title)

	s += "\n「难度」: " + cp.Level

	s += "\n「标签」: " + strings.Join(cp.Tags, "、")
	s += "\n「题目」: " + cp.Abstract

	var hints []string
	for i, hint := range cp.Hints {
		hints = append(hints, fmt.Sprintf("「提示%d」%s", i+1, hint))
	}

	s += "\n" + strings.Join(hints, "\n")
	s += "\n「题解」: https://leetcode-cn.com/problems/" + cp.TitleSlug + "/solution/"
	return s
}

type SubmissionQuestion struct {
	TranslatedTitle string `json:"translatedTitle"`
	FrontendId      string `json:"frontendId"`
	TitleSlug       string `json:"titleSlug"`
	Title           string `json:"title"`
	Difficulty      string `json:"difficulty"`
	LastSubmittedAt int    `json:"lastSubmittedAt"`
	NumSubmitted    int    `json:"numSubmitted"`
}

type SubmissionItem struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Status        string `json:"status"`
	StatusDisplay string `json:"statusDisplay"`
	Lang          string `json:"lang"`
	LangName      string `json:"langName"`
	Runtime       string `json:"runtime"`
	Timestamp     string `json:"timestamp"`
	Url           string `json:"url"`
	IsPending     string `json:"isPending"`
	Memory        string `json:"memory"`
}

type QuestionAndSubmission struct {
	Question    *CodeProblem  `json:"question"`
	Submissions []*Submission `json:"submissions"`
}

type Submission struct {
	Submission *SubmissionItem `json:"submission"`
	Code       string          `json:"code"`
}

func (q *QuestionAndSubmission) ToOutPut(rootFolder string) error {
	// create folder
	key := strings.ReplaceAll(q.Question.TitleSlug, "-", "_")
	path := fmt.Sprintf("%s/lc_%s_%s", rootFolder, q.Question.Index, key)
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, fs.ModePerm)
		if err != nil {
			return err
		}
	} else {
		log.Info("path %s already exists", path)
	}

	// save question text
	questionPath := fmt.Sprintf("%s/question.md", path)
	if _, err := os.Stat(questionPath); err != nil {
		f, err := os.Create(questionPath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(q.Question.String())
		if err != nil {
			return err
		}
	} else {
		log.Info("path %s already exists", questionPath)
	}

	for _, sub := range q.Submissions {
		var ext, header string
		switch sub.Submission.Lang {
		case "golang":
			ext = ".go"
			header = fmt.Sprintf("package lc_%s_%s", q.Question.Index, key)
		case "java":
			ext = ".java"
		case "python3":
			ext = ".py"
		}

		ts, err := strconv.ParseInt(sub.Submission.Timestamp, 10, 64)
		if err != nil {
			return err
		}

		subFile := fmt.Sprintf("%s/submission_%s%s", path, time.Unix(ts, 0).Format("2006_01_02_15_04_05"), ext)
		if _, err := os.Stat(subFile); err != nil {
			f, err := os.Create(subFile)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = f.WriteString(header + "\n\n" + sub.Code)
			if err != nil {
				return err
			}
		}

	}

	//cmd := exec.Command("go", "fmt", fmt.Sprintf("lc_%s_%s", q.Question.Index, key))
	//output, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	//if err != nil {
	//	log.Error("go fmt failed: %s", err)
	//	return err
	//}

	return nil
}
