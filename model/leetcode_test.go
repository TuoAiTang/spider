package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestionAndSubmission_ToOutPut(t *testing.T) {
	var qa *QuestionAndSubmission
	err := json.Unmarshal([]byte(`{"question":{"ID":"122","Index":"122","Title":"买卖股票的最佳时机 II","Tags":["贪心","数组","动态规划"],"Level":"Medium","Abstract":"给你一个整数数组 prices ，其中 prices[i] 表示某支股票第 i 天的价格。\n\n在每一天，你可以决定是否购买和/或出售股票。你在任何时候 最多 只能持有 一股 股票。你也可以先购买，然后在 同一天 出售。\n\n返回 你能获得的 最大 利润 。\n\n \n\n示例 1：\n\n输入：prices = [7,1,5,3,6,4]\n输出：7\n解释：在第 2 天（股票价格 = 1）的时候买入，在第 3 天（股票价格 = 5）的时候卖出, 这笔交易所能获得利润 = 5 - 1 = 4 。\n     随后，在第 4 天（股票价格 = 3）的时候买入，在第 5 天（股票价格 = 6）的时候卖出, 这笔交易所能获得利润 = 6 - 3 = 3 。\n     总利润为 4 + 3 = 7 。\n\n示例 2：\n\n输入：prices = [1,2,3,4,5]\n输出：4\n解释：在第 1 天（股票价格 = 1）的时候买入，在第 5 天 （股票价格 = 5）的时候卖出, 这笔交易所能获得利润 = 5 - 1 = 4 。\n     总利润为 4 。\n\n示例 3：\n\n输入：prices = [7,6,4,3,1]\n输出：0\n解释：在这种情况下, 交易无法获得正利润，所以不参与交易可以获得最大利润，最大利润为 0 。\n\n \n\n提示：\n\n\n 1 <= prices.length <= 3 * 104\n 0 <= prices[i] <= 104\n\n","Hints":null,"TitleSlug":"best-time-to-buy-and-sell-stock-ii"},"submissions":[{"submission":{"id":"314571326","title":"买卖股票的最佳时机 II","status":"AC","statusDisplay":"Accepted","lang":"golang","langName":"Go","runtime":"4 ms","timestamp":"1652757481","url":"/submissions/detail/314571326/","isPending":"Not Pending","memory":"2.9 MB"},"code":"func maxProfit(prices []int) int {\n if len(prices) < 2 {\n  return 0\n }\n    \n    var res int\n\n for i := 1; i < len(prices); i++ {\n        if prices[i] - prices[i-1] > 0 {\n            res += prices[i] - prices[i-1] \n        }\n }\n\n return res\n}"}]}`), &qa)
	assert.Nil(t, err)

	err = qa.ToOutPut("/Users/tuocheng/go/src/github.com/leetcode")
	if err != nil {
		t.Fatal(err)
	}
}
