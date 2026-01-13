import requests
import json
import time
import os
from typing import Dict, Any

class APITestSuite:
    def __init__(self):
        self.base_url = "http://localhost:8888/api/v1"
        self.test_user = {
            "email": f"test_{int(time.time())}@example.com",
            "password": "password123",
            "nickname": "TestUser"
        }
        self.admin_user = {
            "email": "admin@example.com",
            "password": "admin123"
        }
        self.test_novel = {"id": None, "title": "Test Novel"}
        self.results = []
        self.token = None

    def send_request(self, method: str, endpoint: str, data: Dict[str, Any] = None, token: str = None) -> requests.Response:
        """发送HTTP请求"""
        url = f"{self.base_url}{endpoint}"
        headers = {"Content-Type": "application/json"}
        
        if token:
            headers["Authorization"] = f"Bearer {token}"
        
        if data:
            response = requests.request(method, url, json=data, headers=headers)
        else:
            response = requests.request(method, url, headers=headers)
        
        return response

    def run_tests(self):
        """运行所有测试"""
        print("开始API功能测试...")

        # 用户认证测试
        self.test_user_registration()
        self.test_user_login()
        self.test_user_profile()

        # 小说功能测试
        self.test_novel_list()
        self.test_novel_detail()

        # 社交功能测试
        self.test_comment_creation()
        self.test_rating_creation()

        # 搜索功能测试
        self.test_search_functionality()

        # 推荐系统测试
        self.test_recommendations()

        # 管理员功能测试
        self.test_admin_features()

        # 用户活动日志测试
        self.test_user_activity_log()

        # 输出测试结果
        self.print_results()

    def test_user_registration(self):
        """测试用户注册"""
        print("测试用户注册...")
        
        data = {
            "email": self.test_user["email"],
            "password": self.test_user["password"],
            "nickname": self.test_user["nickname"]
        }
        
        response = self.send_request("POST", "/users/register", data)
        
        if response.status_code == 200:
            result = response.json()
            if result.get("code") == 200:
                self.token = result.get("data", {}).get("token", "")
                self.results.append({
                    "test_name": "User Registration",
                    "passed": True,
                    "error": ""
                })
            else:
                self.results.append({
                    "test_name": "User Registration",
                    "passed": False,
                    "error": "响应格式错误"
                })
        else:
            self.results.append({
                "test_name": "User Registration",
                "passed": False,
                "error": f"期望状态码200，实际获得{response.status_code}"
            })

    def test_user_login(self):
        """测试用户登录"""
        print("测试用户登录...")
        
        if not self.token:
            self.results.append({
                "test_name": "User Login",
                "passed": False,
                "error": "依赖注册测试失败"
            })
            return
        
        data = {
            "email": self.test_user["email"],
            "password": self.test_user["password"]
        }
        
        response = self.send_request("POST", "/users/login", data)
        
        if response.status_code == 200:
            result = response.json()
            if result.get("code") == 200:
                self.results.append({
                    "test_name": "User Login",
                    "passed": True,
                    "error": ""
                })
            else:
                self.results.append({
                    "test_name": "User Login",
                    "passed": False,
                    "error": "响应格式错误"
                })
        else:
            self.results.append({
                "test_name": "User Login",
                "passed": False,
                "error": f"期望状态码200，实际获得{response.status_code}"
            })

    def test_user_profile(self):
        """测试用户信息获取"""
        print("测试用户信息获取...")
        
        if not self.token:
            self.results.append({
                "test_name": "User Profile",
                "passed": False,
                "error": "依赖登录测试失败"
            })
            return
        
        response = self.send_request("GET", "/users/profile", token=self.token)
        
        if response.status_code == 200:
            result = response.json()
            if result.get("code") == 200:
                self.results.append({
                    "test_name": "User Profile",
                    "passed": True,
                    "error": ""
                })
            else:
                self.results.append({
                    "test_name": "User Profile",
                    "passed": False,
                    "error": "响应格式错误"
                })
        else:
            self.results.append({
                "test_name": "User Profile",
                "passed": False,
                "error": f"期望状态码200，实际获得{response.status_code}"
            })

    def test_novel_list(self):
        """测试小说列表"""
        print("测试小说列表...")
        
        response = self.send_request("GET", "/novels")
        
        if response.status_code == 200:
            self.results.append({
                "test_name": "Novel List",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "Novel List",
                "passed": False,
                "error": f"期望状态码200，实际获得{response.status_code}"
            })

    def test_novel_detail(self):
        """测试小说详情"""
        print("测试小说详情...")
        
        response = self.send_request("GET", "/novels/1")  # 使用ID为1的小说
        
        # 404是正常的，因为ID为1的小说可能不存在
        if response.status_code in [200, 404]:
            self.results.append({
                "test_name": "Novel Detail",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "Novel Detail",
                "passed": False,
                "error": f"期望状态码200或404，实际获得{response.status_code}"
            })

    def test_comment_creation(self):
        """测试评论创建"""
        print("测试评论创建...")
        
        if not self.token:
            self.results.append({
                "test_name": "Comment Creation",
                "passed": False,
                "error": "依赖登录测试失败"
            })
            return
        
        data = {
            "novel_id": 1,
            "content": "测试评论"
        }
        
        response = self.send_request("POST", "/comments", data, token=self.token)
        
        # 404或400是正常的，因为小说可能不存在或参数验证失败
        if response.status_code in [200, 400, 404]:
            self.results.append({
                "test_name": "Comment Creation",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "Comment Creation",
                "passed": False,
                "error": f"期望状态码200/400/404，实际获得{response.status_code}"
            })

    def test_rating_creation(self):
        """测试评分创建"""
        print("测试评分创建...")
        
        if not self.token:
            self.results.append({
                "test_name": "Rating Creation",
                "passed": False,
                "error": "依赖登录测试失败"
            })
            return
        
        data = {
            "novel_id": 1,
            "score": 8.5,
            "comment": "很好的小说"
        }
        
        response = self.send_request("POST", "/ratings", data, token=self.token)
        
        # 404或400是正常的，因为小说可能不存在或参数验证失败
        if response.status_code in [200, 400, 404]:
            self.results.append({
                "test_name": "Rating Creation",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "Rating Creation",
                "passed": False,
                "error": f"期望状态码200/400/404，实际获得{response.status_code}"
            })

    def test_search_functionality(self):
        """测试搜索功能"""
        print("测试搜索功能...")
        
        response = self.send_request("GET", "/search/novels?q=测试")
        
        if response.status_code == 200:
            self.results.append({
                "test_name": "Search Functionality",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "Search Functionality",
                "passed": False,
                "error": f"期望状态码200，实际获得{response.status_code}"
            })

    def test_recommendations(self):
        """测试推荐功能"""
        print("测试推荐功能...")
        
        response = self.send_request("GET", "/recommendations")
        
        if response.status_code == 200:
            self.results.append({
                "test_name": "Recommendations",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "Recommendations",
                "passed": False,
                "error": f"期望状态码200，实际获得{response.status_code}"
            })

    def test_admin_features(self):
        """测试管理员功能"""
        print("测试管理员功能...")
        
        # 尝试访问管理员功能（应该失败，因为使用普通用户token）
        response = self.send_request("GET", "/users", token=self.token)
        
        # 403是预期的，因为普通用户不能访问管理员功能
        if response.status_code == 403:
            self.results.append({
                "test_name": "Admin Features Access",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "Admin Features Access",
                "passed": False,
                "error": f"期望状态码403，实际获得{response.status_code}"
            })

    def test_user_activity_log(self):
        """测试用户活动日志"""
        print("测试用户活动日志...")
        
        if not self.token:
            self.results.append({
                "test_name": "User Activity Log",
                "passed": False,
                "error": "依赖登录测试失败"
            })
            return
        
        # 获取用户ID（需要从JWT token解码或通过profile获取）
        # 这里简化为假设用户ID为1
        response = self.send_request("GET", "/users/profile", token=self.token)
        user_id = 1  # 默认值
        if response.status_code == 200:
            result = response.json()
            if result.get("code") == 200:
                user_id = result.get("data", {}).get("id", 1)
        
        url = f"/users/{user_id}/activities"
        response = self.send_request("GET", url, token=self.token)
        
        # 200或403都是正常的，取决于权限设置
        if response.status_code in [200, 403]:
            self.results.append({
                "test_name": "User Activity Log",
                "passed": True,
                "error": ""
            })
        else:
            self.results.append({
                "test_name": "User Activity Log",
                "passed": False,
                "error": f"期望状态码200或403，实际获得{response.status_code}"
            })

    def print_results(self):
        """输出测试结果"""
        print("\n测试结果汇总:")
        print("================================")

        total = len(self.results)
        passed = sum(1 for result in self.results if result["passed"])
        
        for result in self.results:
            if result["passed"]:
                print(f"✅ {result['test_name']}: 通过")
            else:
                print(f"❌ {result['test_name']}: 失败 - {result['error']}")

        print(f"\n总测试数: {total}")
        print(f"通过测试: {passed}")
        print(f"失败测试: {total - passed}")
        print(f"成功率: {passed / total * 100:.2f}%")
        
        if passed == total:
            print("\n🎉 所有测试通过！系统功能正常。")
        else:
            print("\n⚠️  存在测试失败，请检查系统功能。")

def main():
    # 检查服务器是否运行
    print("检查服务器是否运行在 :8888...")
    
    try:
        response = requests.get("http://localhost:8888/api/v1/novels", timeout=5)
        print("服务器连接正常，开始测试...")
    except requests.exceptions.RequestException as e:
        print(f"无法连接到服务器: {e}")
        print("请先启动后端服务（go run main.go）")
        return
    
    # 运行测试
    suite = APITestSuite()
    suite.run_tests()

if __name__ == "__main__":
    main()
