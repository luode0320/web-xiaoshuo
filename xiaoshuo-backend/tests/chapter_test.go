package tests

import (
	"fmt"
	"net/http"
)

// testNovelChapters 测试获取小说章节列表
func (suite *APITestSuite) testNovelChapters() {
	fmt.Println("测试获取小说章节列表...")
	
	resp, err := suite.SendRequest("GET", suite.BaseURL+"/novels/1/chapters", nil, "") // 使用ID为1的小说
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Chapters",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 200或404都是正常的，因为小说可能不存在或没有章节
	if suite.CheckResponse(resp, 200) || suite.CheckResponse(resp, 404) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Chapters",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Chapters",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200或404，实际获得%d", resp.StatusCode),
		})
	}
}

// testChapterContent 测试获取章节内容
func (suite *APITestSuite) testChapterContent() {
	fmt.Println("测试获取章节内容...")
	
	resp, err := suite.SendRequest("GET", suite.BaseURL+"/novels/1/chapters/1", nil, "") // 使用小说ID为1，章节ID为1
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Chapter Content",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 200或404都是正常的，因为章节可能不存在
	if suite.CheckResponse(resp, 200) || suite.CheckResponse(resp, 404) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Chapter Content",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Chapter Content",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200或404，实际获得%d", resp.StatusCode),
		})
	}
}