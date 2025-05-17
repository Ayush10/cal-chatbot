# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "Using $goVersion"
} catch {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/"
    exit 1
}

# Check if .env file exists
if (-not (Test-Path ".env")) {
    Write-Host "Creating .env file from template..."
    
    # Create default .env file
    @"
# OpenAI API Configuration
OPENAI_API_KEY=your_openai_api_key
OPENAI_MODEL=gpt-4-turbo

# Cal.com API Configuration
CALCOM_API_KEY=your_calcom_api_key
CALCOM_API_URL=https://api.cal.com/v1
CALCOM_USERNAME=your_calcom_username

# Server Configuration
PORT=8080
DEBUG=true
"@ | Out-File -FilePath ".env" -Encoding utf8

    Write-Host "Please edit the .env file with your API keys before running the application." -ForegroundColor Yellow
    Write-Host "The .env file has been created. Press any key to continue..."
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
}

# Download dependencies
Write-Host "Downloading dependencies..."
go mod download
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to download dependencies" -ForegroundColor Red
    exit 1
}

# Run the application
Write-Host "Starting the Cal.com Chatbot server..."
go run cmd/server/main.go

# Return the exit code from the Go application
exit $LASTEXITCODE 