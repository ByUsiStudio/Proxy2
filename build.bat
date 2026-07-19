@echo off
chcp 65001 >nul
echo ========================================
echo     HP-Lite 全自动编译脚本
echo ========================================
echo.

set BASE_DIR=%~dp0
set WEB_DIR=%BASE_DIR%hp-web
set SERVER_DIR=%BASE_DIR%hp-server-golang

echo [1/4] 检查 Node.js 是否安装...
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误：未检测到 Node.js，请先安装 Node.js
    pause
    exit /b 1
)
echo Node.js 版本：
node --version

echo.
echo [2/4] 构建前端项目...
cd /d "%WEB_DIR%"
if not exist "node_modules" (
    echo 正在安装依赖...
    npm install
    if %errorlevel% neq 0 (
        echo 错误：依赖安装失败
        pause
        exit /b 1
    )
)

echo 正在构建前端...
npm run build
if %errorlevel% neq 0 (
    echo 错误：前端构建失败
    pause
    exit /b 1
)
echo 前端构建成功！

echo.
echo [3/4] 复制前端静态文件到后端...
if not exist "%SERVER_DIR%\static" (
    mkdir "%SERVER_DIR%\static"
)
xcopy /e /y /q "%WEB_DIR%\dist\*" "%SERVER_DIR%\static\"
if %errorlevel% neq 0 (
    echo 错误：复制静态文件失败
    pause
    exit /b 1
)
echo 静态文件复制成功！

echo.
echo [4/4] 构建后端服务...
cd /d "%SERVER_DIR%"

echo 正在构建 Windows 版本...
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o ./target/hp-lite-server.exe main.go
if %errorlevel% neq 0 (
    echo 错误：Windows 版本构建失败
    pause
    exit /b 1
)

echo 正在构建 Linux AMD64 版本...
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ./target/hp-lite-server-amd64 main.go
if %errorlevel% neq 0 (
    echo 错误：Linux AMD64 版本构建失败
    pause
    exit /b 1
)

echo 正在构建 Linux ARM64 版本...
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
go build -o ./target/hp-lite-server-arm64 main.go
if %errorlevel% neq 0 (
    echo 错误：Linux ARM64 版本构建失败
    pause
    exit /b 1
)

echo 正在构建 Linux ARMv7 版本...
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
SET GOARM=7
go build -o ./target/hp-lite-server-armv7 main.go
if %errorlevel% neq 0 (
    echo 错误：Linux ARMv7 版本构建失败
    pause
    exit /b 1
)

echo.
echo ========================================
echo     编译完成！
echo ========================================
echo 输出目录：%SERVER_DIR%\target\
echo.
echo 文件列表：
dir "%SERVER_DIR%\target\" /b
echo.
pause