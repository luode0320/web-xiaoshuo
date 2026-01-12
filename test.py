import requests
import time
import json
import subprocess
import os
from typing import Dict, Any

# 小说阅读系统测试脚本
class NovelSystemTester:
    def __init__(self, base_url: str = "http://localhost:8888"):
        self.base_url = base_url
        self.session = requests.Session()
        self.user_token = None
        self.user_id = None
        self.novel_id = None

    def test_api(self, method: str, endpoint: str, data: Dict[str, Any] = None, headers: Dict[str, str] = None, auth_required: bool = False) -> Dict[str, Any]:
        """通用API测试方法"""
        url = f"{self.base_url}/api/v1{endpoint}"
        
        if headers is None:
            headers = {}
        
        if auth_required and self.user_token:
            headers['Authorization'] = f'Bearer {self.user_token}'
        
        headers['Content-Type'] = 'application/json'
        
        try:
            if method.upper() == 'GET':
                response = self.session.get(url, headers=headers)
            elif method.upper() == 'POST':
                response = self.session.post(url, json=data, headers=headers)
            elif method.upper() == 'PUT':
                response = self.session.put(url, json=data, headers=headers)
            elif method.upper() == 'DELETE':
                response = self.session.delete(url, headers=headers)
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")
            
            print(f"{method} {endpoint} -> {response.status_code}")
            
            try:
                return response.json()
            except:
                return {"error": "Response is not JSON", "text": response.text}
                
        except requests.exceptions.ConnectionError:
            return {"error": f"Cannot connect to {url}"}
        except Exception as e:
            return {"error": str(e)}

    def test_homepage(self):
        """测试首页"""
        print("\n=== 测试首页 ===")
        response = self.test_api('GET', '/')
        print(f"首页访问结果: {response}")
        return True

    def test_user_registration(self):
        """测试用户注册"""
        print("\n=== 测试用户注册 ===")
        user_data = {
            "email": "test@example.com",
            "password": "password123",
            "nickname": "测试用户"
        }
        response = self.test_api('POST', '/users/register', user_data)
        print(f"注册结果: {response}")
        
        if response.get('code') == 200:
            self.user_token = response.get('data', {}).get('token')
            self.user_id = response.get('data', {}).get('user', {}).get('id')
            return True
        return False

    def test_user_login(self):
        """测试用户登录"""
        print("\n=== 测试用户登录 ===")
        login_data = {
            "email": "test@example.com",
            "password": "password123"
        }
        response = self.test_api('POST', '/users/login', login_data)
        print(f"登录结果: {response}")
        
        if response.get('code') == 200:
            self.user_token = response.get('data', {}).get('token')
            self.user_id = response.get('data', {}).get('user', {}).get('id')
            return True
        return False

    def test_get_profile(self):
        """测试获取用户信息"""
        print("\n=== 测试获取用户信息 ===")
        response = self.test_api('GET', '/users/profile', auth_required=True)
        print(f"获取用户信息结果: {response}")
        return response.get('code') == 200

    def test_update_profile(self):
        """测试更新用户信息"""
        print("\n=== 测试更新用户信息 ===")
        update_data = {
            "nickname": "更新后的测试用户"
        }
        response = self.test_api('PUT', '/users/profile', update_data, auth_required=True)
        print(f"更新用户信息结果: {response}")
        return response.get('code') == 200

    def test_get_novels(self):
        """测试获取小说列表"""
        print("\n=== 测试获取小说列表 ===")
        response = self.test_api('GET', '/novels')
        print(f"获取小说列表结果: {response}")
        return response.get('code') == 200

    def test_search_novels(self):
        """测试搜索小说"""
        print("\n=== 测试搜索小说 ===")
        response = self.test_api('GET', '/search/novels?q=测试')
        print(f"搜索小说结果: {response}")
        return response.get('code') == 200

    def test_get_categories(self):
        """测试获取分类"""
        print("\n=== 测试获取分类 ===")
        response = self.test_api('GET', '/categories')
        print(f"获取分类结果: {response}")
        return response.get('code') == 200

    def test_get_rankings(self):
        """测试获取排行榜"""
        print("\n=== 测试获取排行榜 ===")
        response = self.test_api('GET', '/rankings')
        print(f"获取排行榜结果: {response}")
        return response.get('code') == 200

    def test_get_reading_history(self):
        """测试获取阅读历史"""
        print("\n=== 测试获取阅读历史 ===")
        response = self.test_api('GET', '/users/reading-history', auth_required=True)
        print(f"获取阅读历史结果: {response}")
        return response.get('code') == 200

    def test_upload_novel(self):
        """测试上传小说功能(包括EPUB)"""
        print("\n=== 测试上传小说功能 ===")
        # 首先需要创建一个测试用的EPUB文件
        test_epub_path = "test.epub"
        try:
            # 创建一个简单的EPUB文件用于测试
            with open(test_epub_path, 'wb') as f:
                # 写入一个最小化的EPUB文件头（实际应用中需要使用合适的EPUB库）
                f.write(b"PK\x03\x04")  # ZIP文件头，EPUB是ZIP格式
                f.write(b"EPUB Test File")
            
            # 使用multipart/form-data上传
            url = f"{self.base_url}/api/v1/novels/upload"
            headers = {"Authorization": f"Bearer {self.user_token}"} if self.user_token else {}
            
            with open(test_epub_path, 'rb') as f:
                files = {'file': (test_epub_path, f, 'application/epub+zip')}
                data = {
                    'title': '测试EPUB小说',
                    'author': '测试作者',
                    'description': '这是一本用于测试EPUB功能的小说'
                }
                
                response = requests.post(url, files=files, data=data, headers=headers)
                print(f"上传EPUB结果: {response.status_code}")
                try:
                    result = response.json()
                    print(f"上传EPUB响应: {result}")
                    return result.get('code') == 200
                except:
                    print(f"上传EPUB响应非JSON: {response.text}")
                    return False
        except Exception as e:
            print(f"上传EPUB测试异常: {str(e)}")
            return False
        finally:
            # 清理测试文件
            if os.path.exists(test_epub_path):
                os.remove(test_epub_path)

    def run_all_tests(self):
        """运行所有测试"""
        print("开始运行小说阅读系统全流程测试...")
        
        # 检查服务是否运行
        try:
            response = requests.get(f"{self.base_url}/api/v1/novels", timeout=5)
            print("服务正在运行")
        except requests.exceptions.ConnectionError:
            print(f"服务未运行于 {self.base_url}，请先启动后端服务")
            print("启动命令: cd xiaoshuo-backend && go run main.go")
            return {"error": "服务未运行"}
        
        # 基础功能测试
        tests = [
            self.test_homepage,
            self.test_user_registration,
            self.test_user_login,
            self.test_get_profile,
            self.test_update_profile,
            self.test_get_novels,
            self.test_search_novels,
            self.test_get_categories,
            self.test_get_rankings,
        ]
        
        # 条件性测试 - 只有在用户登录成功的情况下才测试需要认证的功能
        if self.user_token:
            tests.extend([
                self.test_get_reading_history
            ])
        
        results = {}
        for test_func in tests:
            try:
                result = test_func()
                results[test_func.__name__] = result
                print(f"{test_func.__name__}: {'通过' if result else '失败'}")
            except Exception as e:
                print(f"{test_func.__name__}: 异常 - {str(e)}")
                results[test_func.__name__] = False
        
        # EPUB上传测试
        try:
            epub_result = self.test_upload_novel()
            results['test_upload_novel'] = epub_result
            print(f"test_upload_novel: {'通过' if epub_result else '失败'}")
        except Exception as e:
            print(f"test_upload_novel: 异常 - {str(e)}")
            results['test_upload_novel'] = False
        
        # 汇总结果
        passed = sum(1 for result in results.values() if result)
        total = len(results)
        
        print(f"\n=== 测试汇总 ===")
        print(f"总测试数: {total}")
        print(f"通过数: {passed}")
        print(f"失败数: {total - passed}")
        print(f"成功率: {passed/total*100:.2f}%" if total > 0 else "成功率: 0%")
        
        return results

def main():
    import sys
    if len(sys.argv) > 1 and sys.argv[1] == "--check-only":
        # 仅检查服务是否运行
        try:
            response = requests.get("http://localhost:8888/api/v1/novels", timeout=5)
            print("服务正在运行")
            sys.exit(0)
        except requests.exceptions.ConnectionError:
            print("服务未运行")
            sys.exit(1)
    
    # 创建测试器实例（默认测试本地运行的服务）
    tester = NovelSystemTester()
    
    # 运行所有测试
    results = tester.run_all_tests()
    
    # 保存测试结果
    with open('test_results.json', 'w', encoding='utf-8') as f:
        json.dump(results, f, ensure_ascii=False, indent=2)
    
    print("\n测试完成，结果已保存到 test_results.json")

if __name__ == "__main__":
    main()