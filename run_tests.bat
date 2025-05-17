@echo off
echo Running Cal.com Chatbot Tests

REM Check if Go is installed
go version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    exit /b 1
)

REM Download dependencies
echo Downloading dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo Failed to download dependencies
    exit /b 1
)

REM Run tests
echo Running tests...
go test ./test
if %ERRORLEVEL% NEQ 0 (
    echo Tests failed
    exit /b 1
)

echo All tests passed!
exit /b 0 