@echo off
echo.
echo ================================
echo   小说阅读系统启动脚本
echo ================================
echo.

REM 检查是否安装了必要的软件
where docker >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo Docker 已安装
) else (
    echo 警告: Docker 未安装，将使用传统方式启动
    goto traditional_start
)

where docker-compose >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo Docker Compose 已安装
    goto docker_start
) else (
    echo 警告: Docker Compose 未安装，将使用传统方式启动
    goto traditional_start
)

:docker_start
echo.
echo 使用 Docker Compose 启动系统...
echo.
cd /d "%~dp0"
docker-compose up -d
if %ERRORLEVEL% EQU 0 (
    echo.
    echo =================================
    echo   系统启动成功！
    echo   前端访问: http://localhost:3000
    echo   后端API: http://localhost:8888/api/v1
    echo =================================
) else (
    echo.
    echo =================================
    echo   Docker启动失败，请检查配置
    echo =================================
)
goto end

:traditional_start
echo.
echo 使用传统方式启动系统...
echo 请确保已安装 Go 和 Node.js
echo.
echo 启动后端服务...
start cmd /k "cd /d %~dp0xiaoshuo-backend && go run main.go"

timeout /t 5 /nobreak >nul

echo 启动前端服务...
start cmd /k "cd /d %~dp0xiaoshuo-frontend && npm run dev"

echo.
echo =================================
echo   系统启动中，请稍候...
echo   前端访问: http://localhost:3000
echo   后端API: http://localhost:8888/api/v1
echo =================================

:end
pause