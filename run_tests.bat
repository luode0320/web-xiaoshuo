@echo off
echo 开始运行小说阅读系统功能测试...

echo.
echo 检查Go环境...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到Go环境，请先安装Go
    pause
    exit /b 1
)

echo.
echo 检查后端服务是否正在运行...
curl -s -o nul -w "%%{http_code}" http://localhost:8888/api/v1/novels >nul 2>&1
if %errorlevel% equ 0 (
    echo 后端服务正在运行
) else (
    echo 后端服务未运行，正在启动...
    start /min cmd /c "cd /d "E:\web-xiaoshuo\xiaoshuo-backend" && go run main.go"
    
    echo 等待服务启动...
    timeout /t 10 >nul
    
    echo 再次检查服务状态...
    curl -s -o nul -w "%%{http_code}" http://localhost:8888/api/v1/novels >nul 2>&1
    if %errorlevel% neq 0 (
        echo 错误: 无法启动后端服务
        pause
        exit /b 1
    )
    echo 服务启动成功
)

echo.
echo 等待几秒确保服务完全启动...
timeout /t 5 >nul

echo.
echo 运行功能测试...
cd /d "E:\web-xiaoshuo"
go run test_runner.go

echo.
echo 测试完成！
pause