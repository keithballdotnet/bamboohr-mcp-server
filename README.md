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
   - `timeOffTypeId` (required): The ID of the time-off type (e.g., '1' for Vacation, '27' for Home Office days)
   - `start` (required): Start date for the time-off request (YYYY-MM-DD format)
   - `end` (required): End date for the time-off request (YYYY-MM-DD format)
   - `employeeNote` (optional): Optional note from the employee about the request

## Setup

### Prerequisites

- Go 1.21 or later
- BambooHR API key
- BambooHR company subdomain

### Installation

1. Clone this repository:
```bash
git clone https://github.com/keithballdotnet/bamboohr-mcp-server.git
cd bamboohr-mcp-server
```

2. Run the setup script to install Git hooks:
```bash
./scripts/setup.sh
```

3. Install dependencies:
```bash
go mod tidy
```

4. Set up environment variables:
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

**Note**: Replace `YOUR_API_KEY` with your actual BambooHR API key and `YOUR_COMPANY` with your company subdomain. The repository includes a pre-commit hook (installed via `./scripts/setup.sh`) that automatically replaces both sensitive values with placeholders when committing to prevent accidental exposure of credentials.

## Example Usage

Once configured, you can use the tools through your MCP client:

### Getting Time-Off Requests
```
Get all time-off requests for employee {employeeId} in 2025
```

### Creating a Time-Off Request
```
Create a home office day request for September 5th, 2025 for employee {employeeId}
```
This will use:
- `employeeId`: "{employeeId}"
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
go build -o .build/bamboohr-mcp-server
```

The binary will be created in the `.build` folder to keep the project root clean.

### Testing

```bash
go test ./...
```

### Building a Release

To build binaries for macOS, Linux, and Windows, follow these steps:

1. **Set Up Your Environment**:
   - Ensure you have Go installed (version 1.21 or later).

2. **Build the Binaries**:
   - Run the following commands:
     ```bash
     go build -o .build/bamboohr-mcp-server-macos main.go
     GOOS=linux GOARCH=amd64 go build -o .build/bamboohr-mcp-server-linux main.go
     GOOS=windows GOARCH=amd64 go build -o .build/bamboohr-mcp-server-windows.exe main.go
     ```

3. **Verify the Binaries**:
   - Ensure the binaries are created in the `.build` directory:
     ```bash
     ls .build
     ```

4. **Create a GitHub Release**:
   - Use the GitHub CLI to create a release:
     ```bash
     gh release create v1.0.0 \
       .build/bamboohr-mcp-server-macos#"macOS binary" \
       .build/bamboohr-mcp-server-linux#"Linux binary" \
       .build/bamboohr-mcp-server-windows.exe#"Windows binary" \
       --title "BambooHR MCP Server v1.0.0" \
       --notes "Initial release of the BambooHR MCP Server with binaries for macOS, Linux, and Windows."
     ```

5. **Publish the Release**:
   - The release will be available on the [Releases Page](https://github.com/keithballdotnet/bamboohr-mcp-server/releases).

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

## Releases

We provide pre-built binaries for macOS, Linux, and Windows. You can download the latest release from the [Releases Page](https://github.com/keithballdotnet/bamboohr-mcp-server/releases).

### Using the Binary

1. **Download the Binary**
   - Visit the [Releases Page](https://github.com/keithballdotnet/bamboohr-mcp-server/releases) and download the appropriate binary for your operating system.

2. **Make the Binary Executable** (if required):
   - On macOS/Linux:
     ```bash
     chmod +x bamboohr-mcp-server-<os>
     ```

3. **Run the Server**:
   - Set the required environment variables:
     ```bash
     export BAMBOOHR_API_KEY="your_api_key"
     export BAMBOOHR_COMPANY="your_company_subdomain"
     ```
   - Start the server:
     ```bash
     ./bamboohr-mcp-server-<os>
     ```

4. **Access the Server**:
   - The server communicates via stdin/stdout following the MCP protocol.
