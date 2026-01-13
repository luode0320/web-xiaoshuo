const puppeteer = require('puppeteer');
const { spawn } = require('child_process');
const path = require('path');

async function runFrontendTests() {
  let backendProcess, frontendProcess;
  
  try {
    console.log('启动后端服务器...');
    backendProcess = spawn('go', ['run', 'main.go'], {
      cwd: path.join(__dirname, '../xiaoshuo-backend'),
      stdio: 'pipe'
    });
    
    backendProcess.stdout.on('data', (data) => {
      const output = data.toString();
      console.log(`后端: ${output}`);
      
      if (output.includes('Server is running on')) {
        console.log('后端服务器已启动');
      }
    });
    
    backendProcess.stderr.on('data', (data) => {
      console.error(`后端错误: ${data}`);
    });
    
    console.log('等待后端启动...');
    await new Promise(resolve => setTimeout(resolve, 5000));
    
    console.log('启动前端开发服务器...');
    frontendProcess = spawn('npm', ['run', 'dev'], {
      cwd: path.join(__dirname, '../xiaoshuo-frontend'),
      stdio: 'pipe'
    });
    
    frontendProcess.stdout.on('data', (data) => {
      const output = data.toString();
      console.log(`前端: ${output}`);
      
      if (output.includes('Local:') || output.includes('Network:')) {
        console.log('前端开发服务器已启动');
      }
    });
    
    frontendProcess.stderr.on('data', (data) => {
      console.error(`前端错误: ${data}`);
    });
    
    console.log('等待前端启动...');
    await new Promise(resolve => setTimeout(resolve, 8000));
    
    // 启动浏览器测试
    console.log('启动浏览器测试...');
    const browser = await puppeteer.launch({ headless: false });
    const page = await browser.newPage();
    
    // 测试搜索功能
    await page.goto('http://localhost:3000/search');
    console.log('已导航到搜索页面');
    
    // 等待页面加载
    await page.waitForSelector('.search-container');
    
    // 测试搜索框
    await page.type('.el-input__inner', '测试小说');
    console.log('已输入搜索关键词');
    
    // 点击搜索按钮
    await page.click('.el-button--primary');
    console.log('已点击搜索按钮');
    
    // 等待搜索结果
    await page.waitForSelector('.search-results');
    console.log('搜索结果已加载');
    
    // 测试热门关键词
    const hotKeywords = await page.$$('.hot-keywords .el-tag');
    console.log(`找到 ${hotKeywords.length} 个热门关键词`);
    
    // 测试筛选功能
    await page.click('.el-collapse-item__header');
    console.log('已展开筛选条件');
    
    // 测试搜索统计（如果用户已登录）
    // 这里可以模拟用户登录来测试搜索历史功能
    
    await browser.close();
    console.log('浏览器测试完成');
    
  } catch (error) {
    console.error('测试过程中发生错误:', error);
  } finally {
    // 清理进程
    if (backendProcess) {
      backendProcess.kill();
    }
    if (frontendProcess) {
      frontendProcess.kill();
    }
  }
}

runFrontendTests();