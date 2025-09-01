# BambooHR MCP Server

An MCP (Model Context Protocol) server that provides access to BambooHR's time-off functionality.

## Features

This server provides the following tools:

### Time-Off Tools

1. **get_time_off_requests** - Retrieve time-off requests for an employee
   - `employeeId` (required): The ID of the employee
   - `start` (optional): Start date for filtering (YYYY-MM-DD format)
   - `end` (optional): End date for filtering (YYYY-MM-DD format)

2. **get_time_off_balance** - Get time-off balance for an employee
   - `employeeId` (required): The ID of the employee

3. **list_employees** - List all employees in the company directory

4. **create_time_off_request** - Create a new time-off request for an employee
   - `employeeId` (required): The ID of the employee to create the time-off request for
   - `start` (required): Start date for the time-off request (YYYY-MM-DD format)
   - `end` (required): End date for the time-off request (YYYY-MM-DD format)
   - `timeOffTypeId` (required): The ID of the time-off type (e.g., '1' for Vacation, '2' for Sick Days, '27' for Home Office days)
   - `amount` (required): The amount of time off in days (e.g., '1', '0.5', '2.5')
   - `notes` (optional): Optional notes for the time-off request
   - `skipManagerApproval` (optional): Whether to skip manager approval (true/false, default: false)

## Setup

### Prerequisites

- Go 1.21 or later
- BambooHR API key
- BambooHR company subdomain

### Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd bamboohr_mcp_server
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
export BAMBOOHR_API_KEY="your_api_key_here"
export BAMBOOHR_COMPANY="your_company_subdomain"
```

### Getting BambooHR Credentials

1. **API Key**: 
   - Log in to your BambooHR account
   - Go to Settings > API Keys
   - Generate a new API key

2. **Company Subdomain**:
   - This is the subdomain in your BambooHR URL
   - For example, if your BambooHR URL is `https://mycompany.bamboohr.com`, then your company subdomain is `mycompany`

### Running the Server

```bash
go run main.go
```

The server will start and listen for MCP requests via stdin/stdout.

## Usage with MCP Clients

This server implements the Model Context Protocol and can be used with any MCP-compatible client.

### VS Code Configuration

For VS Code with the MCP extension, use the configuration file `.vscode/mcp.json`:

```json
{
  "inputs": [],
  "servers": {
    "bamboohr": {
      "command": "./bamboohr-mcp-server",
      "env": {
        "BAMBOOHR_COMPANY": "your_company_subdomain",
        "BAMBOOHR_API_KEY": "YOUR_API_KEY"
      },
      "args": []
    }
  }
}
```

**Note**: Replace `YOUR_API_KEY` with your actual BambooHR API key. The repository includes a pre-commit hook that automatically replaces the API key with a placeholder when committing to prevent accidental exposure of sensitive credentials.

## Example Usage

Once configured, you can use the tools through your MCP client:

### Getting Time-Off Requests
```
Get all time-off requests for employee 157 in 2025
```

### Creating a Time-Off Request
```
Create a home office day request for September 5th, 2025 for employee 157
```
This will use:
- `employeeId`: "157"
- `start`: "2025-09-05"
- `end`: "2025-09-05"
- `timeOffTypeId`: "27" (Home Office days)
- `amount`: "1"

### Common Time-Off Type IDs
- `1` - Vacation
- `2` - Sick Days
- `5` - Sick Day Child
- `19` - Mobile work from abroad
- `27` - Home Office days

### Example Configuration for Claude Desktop

Add this to your Claude Desktop configuration file:

```json
{
  "mcpServers": {
    "bamboohr": {
      "command": "go",
      "args": ["run", "/path/to/bamboohr_mcp_server/main.go"],
      "env": {
        "BAMBOOHR_API_KEY": "your_api_key_here",
        "BAMBOOHR_COMPANY": "your_company_subdomain"
      }
    }
  }
}
```

## API Endpoints Used

This server uses the following BambooHR API endpoints:

- `GET /api/v1/time_off/requests` - Get time-off requests
- `GET /api/gateway.php/{company}/v1/employees/{id}/time_off/calculator` - Get time-off balances
- `GET /api/gateway.php/{company}/v1/employees/directory` - List employees
- `PUT /api/v1/employees/{id}/time_off/request` - Create new time-off request

## Authentication

The server uses HTTP Basic Authentication with the BambooHR API key as the username and an empty password.

## Error Handling

The server provides detailed error messages for:
- Missing or invalid API credentials
- Invalid employee IDs
- API rate limiting
- Network connectivity issues

## Development

### Building

```bash
go build -o bamboohr-mcp-server main.go
```

### Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions:
- Check the BambooHR API documentation: https://documentation.bamboohr.com/
- Create an issue in this repository
