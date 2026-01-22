"""
API接口脚本 - 用于Go后端调用Python解析器
"""
import sys
import json
import argparse
from novel_parser import NovelParser


def api_parse_file(input_path, output_path=None):
    """
    API函数：解析小说文件
    :param input_path: 输入文件路径
    :param output_path: 输出JSON文件路径（可选）
    :return: 解析结果的JSON字符串
    """
    parser = NovelParser()
    result = parser.parse_file(input_path)
    
    if result is None:
        # 返回错误信息
        error_result = {
            "success": False,
            "error": "解析失败",
            "data": None
        }
        return json.dumps(error_result, ensure_ascii=False)
    
    # 返回成功结果
    success_result = {
        "success": True,
        "error": None,
        "data": result
    }
    
    # 如果提供了输出路径，保存到文件
    if output_path:
        try:
            with open(output_path, 'w', encoding='utf-8') as f:
                json.dump(success_result, f, ensure_ascii=False, indent=2)
        except Exception as e:
            error_result = {
                "success": False,
                "error": f"保存文件失败: {str(e)}",
                "data": None
            }
            return json.dumps(error_result, ensure_ascii=False)
    
    return json.dumps(success_result, ensure_ascii=False)


def main():
    parser = argparse.ArgumentParser(description='小说文件解析API')
    parser.add_argument('input_path', help='输入文件路径')
    parser.add_argument('-o', '--output', help='输出JSON文件路径')
    parser.add_argument('--api', action='store_true', help='API模式，输出JSON结果')
    
    args = parser.parse_args()
    
    result_json = api_parse_file(args.input_path, args.output)
    
    if args.api:
        print(result_json)
    else:
        # 解析并打印结果
        result = json.loads(result_json)
        if result["success"]:
            data = result["data"]
            print(f"书名: {data.get('title', 'Unknown')}")
            print(f"作者: {data.get('author', 'Unknown')}")
            print(f"章节数: {len(data.get('chapters', []))}")
            for i, chapter in enumerate(data.get('chapters', [])[:3], 1):  # 只显示前3章
                print(f"  第{i}章: {chapter['title'][:50]}...")
            if len(data.get('chapters', [])) > 3:
                print(f"  ... 还有{len(data.get('chapters', [])) - 3}章")
        else:
            print(f"解析失败: {result['error']}")


if __name__ == "__main__":
    main()