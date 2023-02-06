package leetcode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuoaitang/spider/log"
)

func TestGetToday(t *testing.T) {
	today, err := GetToday()
	assert.Nil(t, err)
	log.Info("today:%+v", today)
}

func Test_htmlToTxt(t *testing.T) {
	txt := htmlToTxt(h)
	fmt.Println(txt)
}

var h = `<p>给你一个正整数 <code>num</code> ，请你统计并返回 <strong>小于或等于</strong> <code>num</code> 且各位数字之和为 <strong>偶数</strong> 的正整数的数目。</p>\n\n<p>正整数的 <strong>各位数字之和</strong> 是其所有位上的对应数字相加的结果。</p>\n\n<p>&nbsp;</p>\n\n<p><strong>示例 1：</strong></p>\n\n<pre>\n<strong>输入：</strong>num = 4\n<strong>输出：</strong>2\n<strong>解释：</strong>\n只有 2 和 4 满足小于等于 4 且各位数字之和为偶数。    \n</pre>\n\n<p><strong>示例 2：</strong></p>\n\n<pre>\n<strong>输入：</strong>num = 30\n<strong>输出：</strong>14\n<strong>解释：</strong>\n只有 14 个整数满足小于等于 30 且各位数字之和为偶数，分别是： \n2、4、6、8、11、13、15、17、19、20、22、24、26 和 28 。\n</pre>\n\n<p>&nbsp;</p>\n\n<p><strong>提示：</strong></p>\n\n<ul>\n\t<li><code>1 &lt;= num &lt;= 1000</code></li>\n</ul>\n`

func Test_getAllSubmission(t *testing.T) {
	questions, err := getSubmissionQuestion()
	assert.Nil(t, err)

	for i, q := range questions {
		if i+1 < 282 {
			continue
		}

		//if (i+1)%5 == 0 {
		//	log.Info("sleep 3s")
		//	time.Sleep(3 * time.Second)
		//}

		log.Info("%d:%s, 进度%.2f%%", i+1, q.Title, float64(i+1)/289*100)
		qa, err := GetQuestionAndSubmissionBySlug(q.TitleSlug)
		assert.Nil(t, err)

		err = qa.ToOutPut("/Users/tuocheng/go/src/github.com/leetcode")
		assert.Nil(t, err)
	}
}
