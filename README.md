# Cal.com Chatbot

An interactive chatbot using OpenAI's function calling capabilities to interact with the Cal.com API.

## Features

- Book new meetings through natural language
- List scheduled events
- Cancel existing events
- Optional web interface for interaction

## Setup

1. Install Go (version 1.21 or later) from [golang.org/dl](https://golang.org/dl/)
2. Clone this repository
3. Run the setup script:
   - Windows: `.\run_app.ps1` (PowerShell) or `run_tests.bat` (Command Prompt)
   - The script will create a `.env` file automatically if one doesn't exist
4. Edit the `.env` file with your API keys:
   ```
   OPENAI_API_KEY=your_openai_api_key
   CALCOM_API_KEY=your_calcom_api_key
   CALCOM_USERNAME=your_calcom_username
   ```
5. The script handles downloading dependencies and starting the server

## API Endpoints

- `POST /api/chat` - Send a message to the chatbot
- `GET /api/events` - Get all scheduled events
- (More endpoints to be added)

## Testing

The project includes a comprehensive test suite:

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test component interactions with mocked dependencies
3. **Mock Implementations**: Test without actual API calls to OpenAI and Cal.com

To run the tests:
- Windows: `.\run_tests.bat` (Command Prompt)
- PowerShell: `go test ./test`

Note: Most tests will automatically skip if API keys are not provided in the environment.

## Project Structure

```
cal-chatbot/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/             # REST API handlers
│   ├── calcom/          # Cal.com API client
│   ├── chatbot/         # Chatbot logic and OpenAI integration
│   └── models/          # Data structures
├── web/                 # Web interface (optional)
├── config/              # Configuration
└── README.md
``` 