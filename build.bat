@echo off
setlocal enabledelayedexpansion

:: Variables
set BINARY_NAME=apix
:: Get version: use tag if exists, otherwise use "dev"
for /f "delims=" %%i in ('git describe --tags --exact-match 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=dev
:: Get commit hash
for /f "delims=" %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT_HASH=%%i
if "%COMMIT_HASH%"=="" set COMMIT_HASH=unknown

echo.
echo Apix Build Script for Windows
echo Version: %VERSION%
echo.

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="build" goto build
if "%1"=="clean" goto clean
if "%1"=="test" goto test
if "%1"=="dev" goto dev
if "%1"=="version" goto version
goto unknown

:help
echo Usage: build.bat [target]
echo.
echo Targets:
echo   build     - Build binary for Windows
echo   clean     - Clean build artifacts
echo   test      - Run tests
echo   dev       - Build and show help
echo   version   - Show version information
echo   help      - Show this help message
echo.
goto end

:build
echo Building %BINARY_NAME% %VERSION%...
go build -ldflags="-s -w -X main.Version=%VERSION%" -o %BINARY_NAME%.exe main.go
if %errorlevel% neq 0 (
    echo Build failed!
    goto end
)
echo Build complete: %BINARY_NAME%.exe
goto end

:clean
echo Cleaning...
del /f /q %BINARY_NAME%.exe 2>nul
del /f /q apix-test*.exe 2>nul
rmdir /s /q dist 2>nul
rmdir /s /q release-assets 2>nul
rmdir /s /q artifacts 2>nul
echo Clean complete
goto end

:test
echo Running tests...
go test -v -race ./...
goto end

:dev
call :build
if %errorlevel% neq 0 goto end
echo.
echo Running %BINARY_NAME%...
%BINARY_NAME%.exe --help
goto end

:version
echo Version: %VERSION%
go version
goto end

:unknown
echo Unknown target: %1
echo Use 'build.bat help' for usage information
goto end

:end
echo.
endlocal
