package utils

import (
	"fmt"
	"xiaoshuo-backend/models"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/v2/search/query"
)

// SearchIndex 全文搜索索引管理器
type SearchIndex struct {
	index bleve.Index
}

var GlobalSearchIndex *SearchIndex

// InitSearchIndex 初始化搜索索引
func InitSearchIndex(indexPath string) error {
	var err error
	GlobalSearchIndex, err = NewSearchIndex(indexPath)
	return err
}

// NewSearchIndex 创建新的搜索索引
func NewSearchIndex(indexPath string) (*SearchIndex, error) {
	// 创建自定义映射
	mapping := bleve.NewIndexMapping()
	
	// 为小说创建映射
	novelMapping := bleve.NewDocumentMapping()
	
	// 为标题、作者、主角、描述等字段配置搜索
	titleFieldMapping := bleve.NewTextFieldMapping()
	titleFieldMapping.Analyzer = standard.Name
	novelMapping.AddFieldMappingsAt("Title", titleFieldMapping)
	
	authorFieldMapping := bleve.NewTextFieldMapping()
	authorFieldMapping.Analyzer = standard.Name
	novelMapping.AddFieldMappingsAt("Author", authorFieldMapping)
	
	protagonistFieldMapping := bleve.NewTextFieldMapping()
	protagonistFieldMapping.Analyzer = standard.Name
	novelMapping.AddFieldMappingsAt("Protagonist", protagonistFieldMapping)
	
	descriptionFieldMapping := bleve.NewTextFieldMapping()
	descriptionFieldMapping.Analyzer = standard.Name
	novelMapping.AddFieldMappingsAt("Description", descriptionFieldMapping)
	
	// 为关键词字段配置搜索
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = standard.Name
	novelMapping.AddFieldMappingsAt("Keywords", keywordFieldMapping)
	
	// 应用映射到小说类型
	mapping.AddDocumentMapping("novel", novelMapping)
	
	// 打开或创建索引
	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		// 如果已存在则打开
		if index, err = bleve.Open(indexPath); err != nil {
			return nil, fmt.Errorf("failed to create/open search index: %v", err)
		}
	}
	
	return &SearchIndex{
		index: index,
	}, nil
}

// IndexNovel 为小说建立索引
func (s *SearchIndex) IndexNovel(novel models.Novel) error {
	// 构建关键词字符串
	var keywordsStr string
	for _, keyword := range novel.Keywords {
		keywordsStr += keyword.Word + " "
	}
	
	// 创建用于索引的文档
	doc := struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Protagonist string `json:"protagonist"`
		Description string `json:"description"`
		Keywords    string `json:"keywords"`
	}{
		ID:          novel.ID,
		Title:       novel.Title,
		Author:      novel.Author,
		Protagonist: novel.Protagonist,
		Description: novel.Description,
		Keywords:    keywordsStr,
	}
	
	// 索引文档
	docID := fmt.Sprintf("novel_%d", novel.ID)
	return s.index.Index(docID, doc)
}

// IndexNovelContent 为小说内容建立索引
func (s *SearchIndex) IndexNovelContent(novelID uint, content string) error {
	doc := struct {
		ID      uint   `json:"id"`
		Content string `json:"content"`
	}{
		ID:      novelID,
		Content: content,
	}
	
	docID := fmt.Sprintf("content_%d", novelID)
	return s.index.Index(docID, doc)
}

// SearchNovels 搜索小说
func (s *SearchIndex) SearchNovels(queryStr string, page, size int) ([]uint, int, error) {
	// 创建布尔查询，组合多个字段的搜索
	boolQuery := bleve.NewBooleanQuery()
	
	// 为不同字段创建匹配查询
	titleQuery := query.NewMatchQuery(queryStr)
	titleQuery.SetField("title")
	boolQuery.AddShould(titleQuery)
	
	authorQuery := query.NewMatchQuery(queryStr)
	authorQuery.SetField("author")
	boolQuery.AddShould(authorQuery)
	
	descriptionQuery := query.NewMatchQuery(queryStr)
	descriptionQuery.SetField("description")
	boolQuery.AddShould(descriptionQuery)
	
	protagonistQuery := query.NewMatchQuery(queryStr)
	protagonistQuery.SetField("protagonist")
	boolQuery.AddShould(protagonistQuery)
	
	keywordsQuery := query.NewMatchQuery(queryStr)
	keywordsQuery.SetField("keywords")
	boolQuery.AddShould(keywordsQuery)

	// 创建搜索请求
	searchRequest := bleve.NewSearchRequest(boolQuery)
	
	// 设置分页
	searchRequest.From = (page - 1) * size
	searchRequest.Size = size
	
	// 执行搜索
	result, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, 0, fmt.Errorf("search failed: %v", err)
	}
	
	// 提取小说ID
	var novelIDs []uint
	for _, hit := range result.Hits {
		var novelID uint
		_, err := fmt.Sscanf(hit.ID, "novel_%d", &novelID)
		if err != nil {
			continue
		}
		novelIDs = append(novelIDs, novelID)
	}
	
	return novelIDs, int(result.Total), nil
}

// SearchNovelContent 搜索小说内容
func (s *SearchIndex) SearchNovelContent(queryStr string, page, size int) ([]uint, int, error) {
	// 创建内容搜索查询
	contentQuery := query.NewMatchQuery(queryStr)
	contentQuery.SetField("content")
	
	// 创建搜索请求
	searchRequest := bleve.NewSearchRequest(contentQuery)
	
	// 设置分页
	searchRequest.From = (page - 1) * size
	searchRequest.Size = size
	
	// 执行搜索
	result, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, 0, fmt.Errorf("content search failed: %v", err)
	}
	
	// 提取小说ID
	var novelIDs []uint
	for _, hit := range result.Hits {
		var novelID uint
		_, err := fmt.Sscanf(hit.ID, "content_%d", &novelID)
		if err != nil {
			continue
		}
		novelIDs = append(novelIDs, novelID)
	}
	
	return novelIDs, int(result.Total), nil
}

// DeleteNovelFromIndex 从索引中删除小说
func (s *SearchIndex) DeleteNovelFromIndex(novelID uint) error {
	docID := fmt.Sprintf("novel_%d", novelID)
	err1 := s.index.Delete(docID)
	
	contentDocID := fmt.Sprintf("content_%d", novelID)
	err2 := s.index.Delete(contentDocID)
	
	if err1 != nil {
		return err1
	}
	return err2
}

// Close 关闭索引
func (s *SearchIndex) Close() error {
	return s.index.Close()
}