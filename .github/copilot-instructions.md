# BambooHR MCP Server - Copilot Instructions

This is a Go-based Model Context Protocol (MCP) server that provides BambooHR time-off functionality through three tools: `get_time_off_requests`, `get_time_off_balance`, and `list_employees`.

## Architecture Overview

**Single-file MCP Server** (`main.go`): All functionality is contained in one file with clear separation:
- `BambooHRClient` struct handles HTTP communication with BambooHR REST API
- Three tool handlers wrap client methods for MCP protocol
- `main()` function sets up server with environment-based configuration

**Key Pattern - Tool Handler Factory**: Each tool uses the factory pattern `handleXXX(client) server.ToolHandlerFunc` returning closures that capture the client instance. This allows clean separation between business logic and MCP protocol handling.

**Authentication Strategy**: Uses HTTP Basic Auth with API key as username, empty password. The `makeRequest()` method centralizes this pattern across all API calls.

## Critical Development Workflows

**Environment Setup**: Always set both required env vars before running:
```bash
export BAMBOOHR_API_KEY="your_key" BAMBOOHR_COMPANY="your_subdomain"
go run main.go
```

**Testing Strategy**: Run `go test -v` for unit tests. Integration testing requires valid BambooHR credentials. Tests focus on client creation and validation, not API calls.

**Building**: Use `go build -o .build/bamboohr-mcp-server` to create standalone binary in the `.build` folder. The server communicates via stdin/stdout following MCP protocol. Always build binaries to the `.build` directory to keep the project root clean.

## Project-Specific Conventions

**Error Handling Pattern**: Always use `mcp.NewToolResultError()` for tool failures, never return Go errors directly. See tool handlers for examples.

**JSON Response Handling**: All API responses are passed through as formatted JSON strings using `json.MarshalIndent()` for readability.

**Struct Design**: BambooHR response structs match API exactly with json tags. The `TimeOffRequest` and `TimeOffBalance` structs are faithful representations of the API schema.

## Integration Points

**MCP Protocol**: Uses `github.com/mark3labs/mcp-go/server` for MCP server implementation. Tools are defined with `mcp.NewTool()` and registered with `s.AddTool()`.

**BambooHR API**: Base URL pattern is `https://{company}.bamboohr.com/api/gateway.php/{company}/v1`. All endpoints require authentication and return JSON.

**Client Integration**: Designed for Claude Desktop - see `claude-desktop-config.json` for configuration pattern. Server runs as subprocess communicating via stdio.

## Key Files & Patterns

- `main.go`: Single source of truth - client, handlers, server setup, and main function
- `main_test.go`: Focuses on client creation validation, not API integration
- `README.md`: Comprehensive setup guide with BambooHR credential instructions
- `USAGE.md`: Tool usage examples with JSON request/response patterns
- `.env.example`: Template showing required environment variables
