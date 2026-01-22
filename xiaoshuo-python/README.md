# 小说文件解析工具 (xiaoshuo-python)

这是一个专门用于解析小说文件的Python工具，用于处理EPUB和TXT格式的小说文件，提取章节信息和内容，以解决Go语言在文本解析方面的局限性。

## 功能特性

- 解析EPUB格式文件，提取章节标题和内容
- 解析TXT格式文件，智能识别章节划分
- 支持多种编码格式的TXT文件（自动检测）
- 识别多种章节标题格式（中文、英文）
- 输出结构化的JSON格式结果
- 与Go后端集成的API接口

## 支持的格式

- EPUB (.epub)
- TXT (.txt)

## 依赖包

- `ebooklib`: EPUB文件处理
- `beautifulsoup4`: HTML解析
- `chardet`: 文件编码检测
- `python-magic`: 文件类型检测
- `PyYAML`: 配置文件处理
- `requests`: HTTP请求处理

## 安装依赖

```bash
pip install -r requirements.txt
```

## 使用方法

### 命令行使用

```bash
# 基本解析
python novel_parser.py /path/to/novel.epub

# 指定输出文件
python novel_parser.py /path/to/novel.epub /path/to/output.json
```

### API模式使用

```bash
# API模式，输出JSON结果
python api_parser.py /path/to/novel.epub --api
```

## 与Go后端集成

Go后端可以通过执行Python脚本并解析其JSON输出来使用此解析器。

## 输出格式

```json
{
  "success": true,
  "error": null,
  "data": {
    "title": "小说标题",
    "author": "作者名",
    "chapters": [
      {
        "title": "章节标题",
        "content": "章节内容",
        "position": 1
      }
    ]
  }
}
```

## 配置文件

`config.json` 包含以下配置项：
- `python_path`: Python解释器路径
- `parser_script`: 解析脚本路径
- `timeout`: 解析超时时间（秒）
- `supported_formats`: 支持的格式列表
- `max_file_size`: 最大文件大小（字节）
- `temp_dir`: 临时文件目录