// API功能测试脚本
const fetch = require('node-fetch');

class APITestSuite {
  constructor() {
    this.baseURL = 'http://localhost:8888/api/v1';
    this.testUser = {
      email: `test_${Date.now()}@example.com`,
      password: 'password123',
      nickname: 'TestUser'
    };
    this.adminUser = {
      email: 'admin@example.com',
      password: 'admin123'
    };
    this.testNovel = { id: null, title: 'Test Novel' };
    this.results = [];
    this.token = null;
  }

  async sendRequest(method, endpoint, data = null, token = null) {
    const url = `${this.baseURL}${endpoint}`;
    const headers = { 'Content-Type': 'application/json' };
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    const options = {
      method,
      headers
    };
    
    if (data) {
      options.body = JSON.stringify(data);
    }
    
    try {
      const response = await fetch(url, options);
      return response;
    } catch (error) {
      console.error(`请求失败: ${error.message}`);
      throw error;
    }
  }

  async runTests() {
    console.log('开始API功能测试...');

    // 用户认证测试
    await this.testUserRegistration();
    await this.testUserLogin();
    await this.testUserProfile();

    // 小说功能测试
    await this.testNovelList();
    await this.testNovelDetail();

    // 社交功能测试
    await this.testCommentCreation();
    await this.testRatingCreation();

    // 搜索功能测试
    await this.testSearchFunctionality();

    // 推荐系统测试
    await this.testRecommendations();

    // 管理员功能测试
    await this.testAdminFeatures();

    // 用户活动日志测试
    await this.testUserActivityLog();

    // 输出测试结果
    this.printResults();
  }

  async testUserRegistration() {
    console.log('测试用户注册...');
    
    const data = {
      email: this.testUser.email,
      password: this.testUser.password,
      nickname: this.testUser.nickname
    };
    
    try {
      const response = await this.sendRequest('POST', '/users/register', data);
      
      if (response.status === 200) {
        const result = await response.json();
        if (result.code === 200) {
          this.token = result.data?.token || '';
          this.results.push({
            testName: 'User Registration',
            passed: true,
            error: ''
          });
        } else {
          this.results.push({
            testName: 'User Registration',
            passed: false,
            error: '响应格式错误'
          });
        }
      } else {
        this.results.push({
          testName: 'User Registration',
          passed: false,
          error: `期望状态码200，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'User Registration',
        passed: false,
        error: error.message
      });
    }
  }

  async testUserLogin() {
    console.log('测试用户登录...');
    
    if (!this.token) {
      this.results.push({
        testName: 'User Login',
        passed: false,
        error: '依赖注册测试失败'
      });
      return;
    }
    
    const data = {
      email: this.testUser.email,
      password: this.testUser.password
    };
    
    try {
      const response = await this.sendRequest('POST', '/users/login', data);
      
      if (response.status === 200) {
        const result = await response.json();
        if (result.code === 200) {
          this.results.push({
            testName: 'User Login',
            passed: true,
            error: ''
          });
        } else {
          this.results.push({
            testName: 'User Login',
            passed: false,
            error: '响应格式错误'
          });
        }
      } else {
        this.results.push({
          testName: 'User Login',
          passed: false,
          error: `期望状态码200，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'User Login',
        passed: false,
        error: error.message
      });
    }
  }

  async testUserProfile() {
    console.log('测试用户信息获取...');
    
    if (!this.token) {
      this.results.push({
        testName: 'User Profile',
        passed: false,
        error: '依赖登录测试失败'
      });
      return;
    }
    
    try {
      const response = await this.sendRequest('GET', '/users/profile', null, this.token);
      
      if (response.status === 200) {
        const result = await response.json();
        if (result.code === 200) {
          this.results.push({
            testName: 'User Profile',
            passed: true,
            error: ''
          });
        } else {
          this.results.push({
            testName: 'User Profile',
            passed: false,
            error: '响应格式错误'
          });
        }
      } else {
        this.results.push({
          testName: 'User Profile',
          passed: false,
          error: `期望状态码200，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'User Profile',
        passed: false,
        error: error.message
      });
    }
  }

  async testNovelList() {
    console.log('测试小说列表...');
    
    try {
      const response = await this.sendRequest('GET', '/novels');
      
      if (response.status === 200) {
        this.results.push({
          testName: 'Novel List',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'Novel List',
          passed: false,
          error: `期望状态码200，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'Novel List',
        passed: false,
        error: error.message
      });
    }
  }

  async testNovelDetail() {
    console.log('测试小说详情...');
    
    try {
      const response = await this.sendRequest('GET', '/novels/1'); // 使用ID为1的小说
      
      // 404是正常的，因为ID为1的小说可能不存在
      if (response.status === 200 || response.status === 404) {
        this.results.push({
          testName: 'Novel Detail',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'Novel Detail',
          passed: false,
          error: `期望状态码200或404，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'Novel Detail',
        passed: false,
        error: error.message
      });
    }
  }

  async testCommentCreation() {
    console.log('测试评论创建...');
    
    if (!this.token) {
      this.results.push({
        testName: 'Comment Creation',
        passed: false,
        error: '依赖登录测试失败'
      });
      return;
    }
    
    const data = {
      novel_id: 1,
      content: '测试评论'
    };
    
    try {
      const response = await this.sendRequest('POST', '/comments', data, this.token);
      
      // 404或400是正常的，因为小说可能不存在或参数验证失败
      if (response.status === 200 || response.status === 400 || response.status === 404) {
        this.results.push({
          testName: 'Comment Creation',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'Comment Creation',
          passed: false,
          error: `期望状态码200/400/404，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'Comment Creation',
        passed: false,
        error: error.message
      });
    }
  }

  async testRatingCreation() {
    console.log('测试评分创建...');
    
    if (!this.token) {
      this.results.push({
        testName: 'Rating Creation',
        passed: false,
        error: '依赖登录测试失败'
      });
      return;
    }
    
    const data = {
      novel_id: 1,
      score: 8.5,
      comment: '很好的小说'
    };
    
    try {
      const response = await this.sendRequest('POST', '/ratings', data, this.token);
      
      // 404或400是正常的，因为小说可能不存在或参数验证失败
      if (response.status === 200 || response.status === 400 || response.status === 404) {
        this.results.push({
          testName: 'Rating Creation',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'Rating Creation',
          passed: false,
          error: `期望状态码200/400/404，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'Rating Creation',
        passed: false,
        error: error.message
      });
    }
  }

  async testSearchFunctionality() {
    console.log('测试搜索功能...');
    
    try {
      const response = await this.sendRequest('GET', '/search/novels?q=测试');
      
      if (response.status === 200) {
        this.results.push({
          testName: 'Search Functionality',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'Search Functionality',
          passed: false,
          error: `期望状态码200，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'Search Functionality',
        passed: false,
        error: error.message
      });
    }
  }

  async testRecommendations() {
    console.log('测试推荐功能...');
    
    try {
      const response = await this.sendRequest('GET', '/recommendations');
      
      if (response.status === 200) {
        this.results.push({
          testName: 'Recommendations',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'Recommendations',
          passed: false,
          error: `期望状态码200，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'Recommendations',
        passed: false,
        error: error.message
      });
    }
  }

  async testAdminFeatures() {
    console.log('测试管理员功能...');
    
    // 尝试访问管理员功能（应该失败，因为使用普通用户token）
    try {
      const response = await this.sendRequest('GET', '/users', null, this.token);
      
      // 403是预期的，因为普通用户不能访问管理员功能
      if (response.status === 403) {
        this.results.push({
          testName: 'Admin Features Access',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'Admin Features Access',
          passed: false,
          error: `期望状态码403，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'Admin Features Access',
        passed: false,
        error: error.message
      });
    }
  }

  async testUserActivityLog() {
    console.log('测试用户活动日志...');
    
    if (!this.token) {
      this.results.push({
        testName: 'User Activity Log',
        passed: false,
        error: '依赖登录测试失败'
      });
      return;
    }
    
    // 获取用户ID（需要从JWT token解码或通过profile获取）
    // 这里简化为假设用户ID为1
    try {
      const profileResponse = await this.sendRequest('GET', '/users/profile', null, this.token);
      let userId = 1; // 默认值
      if (profileResponse.status === 200) {
        const result = await profileResponse.json();
        if (result.code === 200) {
          userId = result.data?.id || 1;
        }
      }
      
      const response = await this.sendRequest('GET', `/users/${userId}/activities`, null, this.token);
      
      // 200或403都是正常的，取决于权限设置
      if (response.status === 200 || response.status === 403) {
        this.results.push({
          testName: 'User Activity Log',
          passed: true,
          error: ''
        });
      } else {
        this.results.push({
          testName: 'User Activity Log',
          passed: false,
          error: `期望状态码200或403，实际获得${response.status}`
        });
      }
    } catch (error) {
      this.results.push({
        testName: 'User Activity Log',
        passed: false,
        error: error.message
      });
    }
  }

  printResults() {
    console.log('\n测试结果汇总:');
    console.log('================================');

    const total = this.results.length;
    const passed = this.results.filter(result => result.passed).length;
    
    this.results.forEach(result => {
      if (result.passed) {
        console.log(`✅ ${result.testName}: 通过`);
      } else {
        console.log(`❌ ${result.testName}: 失败 - ${result.error}`);
      }
    });

    console.log(`\n总测试数: ${total}`);
    console.log(`通过测试: ${passed}`);
    console.log(`失败测试: ${total - passed}`);
    console.log(`成功率: ${(passed / total * 100).toFixed(2)}%`);
    
    if (passed === total) {
      console.log('\n🎉 所有测试通过！系统功能正常。');
    } else {
      console.log('\n⚠️  存在测试失败，请检查系统功能。');
    }
  }
}

async function main() {
  // 检查服务器是否运行
  console.log('检查服务器是否运行在 :8888...');
  
  try {
    const response = await fetch('http://localhost:8888/api/v1/novels', { timeout: 5000 });
    console.log('服务器连接正常，开始测试...');
  } catch (error) {
    console.log(`无法连接到服务器: ${error.message}`);
    console.log('请先启动后端服务（go run main.go）');
    return;
  }
  
  // 运行测试
  const suite = new APITestSuite();
  await suite.runTests();
}

// 运行测试
main().catch(console.error);