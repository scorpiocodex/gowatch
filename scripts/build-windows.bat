@echo off
REM Windows Build Script for GoWatch

setlocal EnableDelayedExpansion

echo ========================================
echo GoWatch - Windows Build Script
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    exit /b 1
)

echo [INFO] Go version:
go version
echo.

REM Parse command line arguments
set "COMMAND=%~1"
if "%COMMAND%"=="" set "COMMAND=help"

if /i "%COMMAND%"=="build" goto :build
if /i "%COMMAND%"=="test" goto :test
if /i "%COMMAND%"=="clean" goto :clean
if /i "%COMMAND%"=="install" goto :install
if /i "%COMMAND%"=="run" goto :run
if /i "%COMMAND%"=="help" goto :help
goto :help

:build
echo [BUILD] Building gowatch.exe...
if not exist "bin" mkdir bin
go build -o bin\gowatch.exe .\cmd\gowatch
if %errorlevel% equ 0 (
    echo [SUCCESS] Built: bin\gowatch.exe
) else (
    echo [ERROR] Build failed
    exit /b 1
)
goto :end

:test
echo [TEST] Running tests...
go test -v ./...
if %errorlevel% equ 0 (
    echo [SUCCESS] All tests passed
) else (
    echo [ERROR] Tests failed
    exit /b 1
)
goto :end

:clean
echo [CLEAN] Cleaning build artifacts...
if exist "bin" (
    rmdir /s /q bin
    echo [SUCCESS] Cleaned bin directory
)
if exist "coverage.out" del coverage.out
if exist "coverage.html" del coverage.html
echo [SUCCESS] Cleaned
goto :end

:install
echo [INSTALL] Installing gowatch...
go install .\cmd\gowatch
if %errorlevel% equ 0 (
    echo [SUCCESS] Installed to: %GOPATH%\bin\gowatch.exe
    echo Make sure %GOPATH%\bin is in your PATH
) else (
    echo [ERROR] Installation failed
    exit /b 1
)
goto :end

:run
call :build
if %errorlevel% equ 0 (
    echo.
    echo [RUN] Starting gowatch...
    .\bin\gowatch.exe %*
)
goto :end

:help
echo Usage: build-windows.bat [command]
echo.
echo Commands:
echo   build      Build the gowatch.exe binary
echo   test       Run all tests
echo   clean      Clean build artifacts
echo   install    Install to GOPATH\bin
echo   run        Build and run gowatch
echo   help       Show this help message
echo.
echo Examples:
echo   build-windows.bat build
echo   build-windows.bat test
echo   build-windows.bat run --help
goto :end

:end
echo.
endlocal