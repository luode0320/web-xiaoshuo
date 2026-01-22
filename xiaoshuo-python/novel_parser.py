"""
小说文件解析器
支持解析EPUB和TXT格式的小说文件，提取章节信息和内容
"""
import os
import json
import re
import chardet
from ebooklib import epub
from bs4 import BeautifulSoup
import sys
from pathlib import Path


class NovelParser:
    def __init__(self):
        self.supported_formats = ['.epub', '.txt']
    
    def detect_encoding(self, file_path):
        """检测文件编码"""
        with open(file_path, 'rb') as f:
            raw_data = f.read()
            result = chardet.detect(raw_data)
            return result['encoding']
    
    def parse_epub(self, file_path):
        """解析EPUB文件"""
        try:
            book = epub.read_epub(file_path)
            
            novel_info = {
                'title': '',
                'author': '',
                'chapters': []
            }
            
            # 获取元数据
            if book.metadata.get('http://purl.org/dc/elements/1.1/', {}).get('title'):
                novel_info['title'] = book.metadata['http://purl.org/dc/elements/1.1/']['title'][0][0]
            if book.metadata.get('http://purl.org/dc/elements/1.1/', {}).get('creator'):
                novel_info['author'] = book.metadata['http://purl.org/dc/elements/1.1/']['creator'][0][0]
            
            # 提取章节
            chapter_count = 0
            for item in book.get_items():
                if item.get_type() == epub.EpubHtml:
                    chapter_count += 1
                    soup = BeautifulSoup(item.get_content(), 'html.parser')
                    
                    # 尝试从标题中获取章节名
                    chapter_title = f"第{chapter_count}章"
                    title_elem = soup.find(['h1', 'h2', 'h3', 'h4', 'h5', 'h6'])
                    if title_elem:
                        chapter_title = title_elem.get_text().strip()
                    else:
                        # 尝试从内容中提取章节标题
                        content_text = soup.get_text()
                        # 匹配常见的章节标题模式
                        patterns = [
                            r'第[一二三四五六七八九十零\d]+[章节回部篇]',  # 第X章、第X节等
                            r'Chapter\s+\d+',  # Chapter X
                            r'chapter\s+\d+',  # chapter X
                            r'正文\s+',  # 正文
                            r'序章\s*',  # 序章
                            r'引子\s*',  # 引子
                            r'楔子\s*',  # 楔子
                            r'尾声\s*'   # 尾声
                        ]
                        for pattern in patterns:
                            match = re.search(pattern, content_text[:200])  # 只检查前200个字符
                            if match:
                                chapter_title = match.group().strip()
                                break
                    
                    chapter = {
                        'title': chapter_title,
                        'content': soup.get_text().strip(),
                        'position': chapter_count
                    }
                    
                    novel_info['chapters'].append(chapter)
            
            return novel_info
        except Exception as e:
            print(f"解析EPUB文件时出错: {str(e)}", file=sys.stderr)
            return None
    
    def parse_txt(self, file_path):
        """解析TXT文件"""
        try:
            # 检测编码
            encoding = self.detect_encoding(file_path)
            
            with open(file_path, 'r', encoding=encoding, errors='ignore') as f:
                content = f.read()
            
            novel_info = {
                'title': os.path.basename(file_path).replace('.txt', ''),
                'author': '',
                'chapters': []
            }
            
            # 按常见的章节模式分割内容
            # 支持多种章节标题格式
            chapter_patterns = [
                r'(第[一二三四五六七八九十零\d]+[章节回部篇]\s*[：:]?\s*[^
]*)',
                r'(Chapter\s+\d+\s*[：:]?\s*[^
]*)',
                r'(chapter\s+\d+\s*[：:]?\s*[^
]*)',
                r'(正文\s+[^
]*)',
                r'(序章\s*[^
]*)',
                r'(引子\s*[^
]*)',
                r'(楔子\s*[^
]*)',
                r'(尾声\s*[^
]*)',
                r'(第[一二三四五六七八九十零\d]+卷\s*[^
]*)',
                r'(新第[一二三四五六七八九十零\d]+[章节回部篇]\s*[^
]*)'  # 处理可能的OCR错误
            ]
            
            # 尝试使用各种模式分割
            chapters_found = False
            for pattern in chapter_patterns:
                # 使用正则表达式分割，保留分隔符
                parts = re.split(f'({pattern})', content)
                
                # 检查是否找到足够的章节（至少2章）
                chapter_titles = [part for part in parts if re.match(pattern, part.strip())]
                
                if len(chapter_titles) >= 2:  # 如果找到至少2个章节标题
                    chapters = []
                    current_title = "序章"  # 初始章节名
                    current_content = ""
                    
                    # 从第二部分开始（因为第一部分通常是章前内容）
                    i = 1
                    while i < len(parts):
                        part = parts[i].strip()
                        if re.match(pattern, part):
                            # 这是一个章节标题
                            if current_content.strip():  # 保存前一章节
                                chapters.append({
                                    'title': current_title.strip(),
                                    'content': current_content.strip(),
                                    'position': len(chapters) + 1
                                })
                            
                            current_title = part
                            current_content = ""
                        else:
                            # 这是章节内容
                            current_content += part
                        
                        i += 1
                    
                    # 添加最后一个章节
                    if current_content.strip():
                        chapters.append({
                            'title': current_title.strip(),
                            'content': current_content.strip(),
                            'position': len(chapters) + 1
                        })
                    
                    if chapters:
                        novel_info['chapters'] = chapters
                        chapters_found = True
                        break
            
            # 如果没有找到明确的章节划分，创建单个章节
            if not chapters_found:
                novel_info['chapters'] = [{
                    'title': '正文',
                    'content': content.strip(),
                    'position': 1
                }]
            
            return novel_info
        except Exception as e:
            print(f"解析TXT文件时出错: {str(e)}", file=sys.stderr)
            return None
    
    def parse_file(self, file_path):
        """根据文件扩展名选择解析方法"""
        file_path = Path(file_path)
        if not file_path.exists():
            print(f"文件不存在: {file_path}", file=sys.stderr)
            return None
        
        extension = file_path.suffix.lower()
        
        if extension not in self.supported_formats:
            print(f"不支持的文件格式: {extension}", file=sys.stderr)
            return None
        
        if extension == '.epub':
            return self.parse_epub(file_path)
        elif extension == '.txt':
            return self.parse_txt(file_path)
    
    def save_result(self, result, output_path):
        """保存解析结果到JSON文件"""
        try:
            with open(output_path, 'w', encoding='utf-8') as f:
                json.dump(result, f, ensure_ascii=False, indent=2)
            return True
        except Exception as e:
            print(f"保存结果时出错: {str(e)}", file=sys.stderr)
            return False


def main():
    if len(sys.argv) < 2:
        print("使用方法: python novel_parser.py <输入文件路径> [输出文件路径]")
        sys.exit(1)
    
    input_file = sys.argv[1]
    
    # 如果提供了输出路径，使用它，否则生成默认输出路径
    if len(sys.argv) > 2:
        output_file = sys.argv[2]
    else:
        input_path = Path(input_file)
        output_file = input_path.with_name(f"{input_path.stem}_parsed.json")
    
    # 创建解析器实例
    parser = NovelParser()
    
    # 解析文件
    result = parser.parse_file(input_file)
    
    if result:
        # 保存结果
        if parser.save_result(result, output_file):
            print(f"解析完成，结果已保存到: {output_file}")
            print(json.dumps(result, ensure_ascii=False, indent=2))
        else:
            print("保存结果失败", file=sys.stderr)
            sys.exit(1)
    else:
        print("解析失败", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()