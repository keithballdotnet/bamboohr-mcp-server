package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// BambooHRClient represents a client for the BambooHR API
type BambooHRClient struct {
	BaseURL    string
	APIKey     string
	Company    string
	HTTPClient *http.Client
}

// TimeOffRequest represents a time-off request
type TimeOffRequest struct {
	ID         string `json:"id"`
	EmployeeID string `json:"employeeId"`
	Name       string `json:"name"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Created    string `json:"created"`
	Type       struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"type"`
	Amount struct {
		Unit   string        `json:"unit"`
		Amount FlexibleFloat `json:"amount"`
	} `json:"amount"`
	Notes  map[string]string `json:"notes"`
	Status struct {
		Status          string `json:"status"`
		LastChanged     string `json:"lastChanged"`
		LastChangedByID string `json:"lastChangedByUserId"`
	} `json:"status"`
	Actions struct {
		View    bool `json:"view"`
		Edit    bool `json:"edit"`
		Cancel  bool `json:"cancel"`
		Approve bool `json:"approve"`
		Deny    bool `json:"deny"`
		Bypass  bool `json:"bypass"`
	} `json:"actions"`
	Dates map[string]string `json:"dates"`
}

// TimeOffBalance represents time-off balance for an employee
type TimeOffBalance struct {
	TimeOffType    string        `json:"timeOffType"`
	Name           string        `json:"name"`
	Units          string        `json:"units"`
	Balance        FlexibleFloat `json:"balance"`
	End            string        `json:"end"`
	PolicyType     string        `json:"policyType"`
	UsedYearToDate FlexibleFloat `json:"usedYearToDate"`
}

// TimeOffRequestCreate represents the payload for creating a time-off request
type TimeOffRequestCreate struct {
	TimeOffTypeID string `json:"timeOffTypeId"`
	Start         string `json:"start"`
	End           string `json:"end"`
	EmployeeNote  string `json:"employeeNote,omitempty"`
}

// FlexibleFloat can unmarshal both string and float64 values
type FlexibleFloat float64

func (f *FlexibleFloat) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float64 first
	var floatVal float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		*f = FlexibleFloat(floatVal)
		return nil
	}

	// If that fails, try to unmarshal as string and convert
	var stringVal string
	if err := json.Unmarshal(data, &stringVal); err != nil {
		return err
	}

	// Handle empty string as 0
	if stringVal == "" {
		*f = FlexibleFloat(0)
		return nil
	}

	// Parse string to float
	floatVal, err := strconv.ParseFloat(stringVal, 64)
	if err != nil {
		return fmt.Errorf("cannot parse '%s' as float: %w", stringVal, err)
	}

	*f = FlexibleFloat(floatVal)
	return nil
}

func (f FlexibleFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

// NewBambooHRClient creates a new BambooHR API client
func NewBambooHRClient(company, apiKey string) *BambooHRClient {
	return &BambooHRClient{
		BaseURL:    fmt.Sprintf("https://%s.bamboohr.com/api/gateway.php/%s/v1", company, company),
		APIKey:     apiKey,
		Company:    company,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// makeRequest performs an HTTP request to the BambooHR API
func (c *BambooHRClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Basic authentication with API key as username and empty password
	req.SetBasicAuth(c.APIKey, "")
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.HTTPClient.Do(req)
}

// makeRequestV1 performs an HTTP request to the newer BambooHR API v1 format
func (c *BambooHRClient) makeRequestV1(method, endpoint string, body io.Reader) (*http.Response, error) {
	baseURL := fmt.Sprintf("https://%s.bamboohr.com/api/v1", c.Company)
	url := baseURL + endpoint

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Basic authentication with API key as username and empty password
	req.SetBasicAuth(c.APIKey, "")
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.HTTPClient.Do(req)
}

// GetTimeOffRequests retrieves time-off requests for a given employee
func (c *BambooHRClient) GetTimeOffRequests(employeeID int, start, end string) ([]TimeOffRequest, error) {
	// Default to current year if no dates provided
	if start == "" || end == "" {
		now := time.Now()
		start = fmt.Sprintf("%d-01-01", now.Year())
		end = fmt.Sprintf("%d-12-31", now.Year())
	}

	// Try the exact format from documentation: /time_off/requests with required start/end params
	endpoint := "/time_off/requests"
	params := fmt.Sprintf("?start=%s&end=%s&employeeId=%d", start, end, employeeID)
	endpoint += params

	resp, err := c.makeRequestV1("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var requests []TimeOffRequest
	if err := json.Unmarshal(body, &requests); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return requests, nil
}

// GetTimeOffBalance retrieves time-off balance for an employee
func (c *BambooHRClient) GetTimeOffBalance(employeeID int) ([]TimeOffBalance, error) {
	endpoint := fmt.Sprintf("/employees/%d/time_off/calculator", employeeID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var balances []TimeOffBalance
	if err := json.Unmarshal(body, &balances); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return balances, nil
}

// CreateTimeOffRequest creates a new time-off request for an employee
func (c *BambooHRClient) CreateTimeOffRequest(employeeID int, request TimeOffRequestCreate) (*TimeOffRequest, error) {
	endpoint := fmt.Sprintf("/employees/%d/time_off/request", employeeID)

	// Marshal the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.makeRequestV1("PUT", endpoint, strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var createdRequest TimeOffRequest
	if err := json.Unmarshal(body, &createdRequest); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &createdRequest, nil
}

// Tool handlers

func handleGetTimeOffRequests(client *BambooHRClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		employeeIDStr, err := request.RequireString("employeeId")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("employeeId is required: %s", err.Error())), nil
		}

		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			return mcp.NewToolResultError("employeeId must be a valid integer"), nil
		}

		start := request.GetString("start", "")
		end := request.GetString("end", "")

		requests, err := client.GetTimeOffRequests(employeeID, start, end)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get time-off requests: %s", err.Error())), nil
		}

		data, err := json.MarshalIndent(requests, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %s", err.Error())), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func handleGetTimeOffBalance(client *BambooHRClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		employeeIDStr, err := request.RequireString("employeeId")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("employeeId is required: %s", err.Error())), nil
		}

		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			return mcp.NewToolResultError("employeeId must be a valid integer"), nil
		}

		balances, err := client.GetTimeOffBalance(employeeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get time-off balance: %s", err.Error())), nil
		}

		data, err := json.MarshalIndent(balances, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %s", err.Error())), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func handleListEmployees(client *BambooHRClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		endpoint := "/employees/directory"

		resp, err := client.makeRequest("GET", endpoint, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to make request: %s", err.Error())), nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return mcp.NewToolResultError(fmt.Sprintf("API error %d: %s", resp.StatusCode, string(body))), nil
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read response: %s", err.Error())), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

func handleCreateTimeOffRequest(client *BambooHRClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		employeeIDStr, err := request.RequireString("employeeId")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("employeeId is required: %s", err.Error())), nil
		}

		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			return mcp.NewToolResultError("employeeId must be a valid integer"), nil
		}

		timeOffTypeID, err := request.RequireString("timeOffTypeId")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("timeOffTypeId is required: %s", err.Error())), nil
		}

		startDate, err := request.RequireString("start")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("start date is required: %s", err.Error())), nil
		}

		endDate, err := request.RequireString("end")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("end date is required: %s", err.Error())), nil
		}

		// Optional employee note - for now we'll keep it empty if not provided
		employeeNote := request.GetString("employeeNote", "")

		// Create the request payload
		timeOffRequest := TimeOffRequestCreate{
			TimeOffTypeID: timeOffTypeID,
			Start:         startDate,
			End:           endDate,
			EmployeeNote:  employeeNote,
		}

		createdRequest, err := client.CreateTimeOffRequest(employeeID, timeOffRequest)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create time-off request: %s", err.Error())), nil
		}

		data, err := json.MarshalIndent(createdRequest, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %s", err.Error())), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func main() {
	// Get configuration from environment variables
	apiKey := os.Getenv("BAMBOOHR_API_KEY")
	company := os.Getenv("BAMBOOHR_COMPANY")

	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: BAMBOOHR_API_KEY environment variable is required\n")
		os.Exit(1)
	}

	if company == "" {
		fmt.Fprintf(os.Stderr, "Error: BAMBOOHR_COMPANY environment variable is required\n")
		os.Exit(1)
	}

	// Create BambooHR client
	client := NewBambooHRClient(company, apiKey)

	// Create MCP server
	s := server.NewMCPServer(
		"BambooHR Time-Off MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Define tools
	getTimeOffRequestsTool := mcp.NewTool(
		"get_time_off_requests",
		mcp.WithDescription("Get time-off requests for an employee"),
		mcp.WithString("employeeId",
			mcp.Required(),
			mcp.Description("The ID of the employee to get time-off requests for"),
		),
		mcp.WithString("start",
			mcp.Description("Start date for filtering requests (YYYY-MM-DD format). Optional."),
		),
		mcp.WithString("end",
			mcp.Description("End date for filtering requests (YYYY-MM-DD format). Optional."),
		),
	)

	getTimeOffBalanceTool := mcp.NewTool(
		"get_time_off_balance",
		mcp.WithDescription("Get time-off balance for an employee"),
		mcp.WithString("employeeId",
			mcp.Required(),
			mcp.Description("The ID of the employee to get time-off balance for"),
		),
	)

	listEmployeesTool := mcp.NewTool(
		"list_employees",
		mcp.WithDescription("List all employees in the company directory"),
	)

	createTimeOffRequestTool := mcp.NewTool(
		"create_time_off_request",
		mcp.WithDescription("Create a new time-off request for an employee"),
		mcp.WithString("employeeId",
			mcp.Required(),
			mcp.Description("The ID of the employee to create the time-off request for"),
		),
		mcp.WithString("timeOffTypeId",
			mcp.Required(),
			mcp.Description("The ID of the time-off type (e.g., '1' for Vacation, '27' for Home Office days)"),
		),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Description("Start date for the time-off request (YYYY-MM-DD format)"),
		),
		mcp.WithString("end",
			mcp.Required(),
			mcp.Description("End date for the time-off request (YYYY-MM-DD format)"),
		),
		mcp.WithString("employeeNote",
			mcp.Description("Optional note from the employee about the request"),
		),
	)

	// Add tools to server
	s.AddTool(getTimeOffRequestsTool, handleGetTimeOffRequests(client))
	s.AddTool(getTimeOffBalanceTool, handleGetTimeOffBalance(client))
	s.AddTool(listEmployeesTool, handleListEmployees(client))
	s.AddTool(createTimeOffRequestTool, handleCreateTimeOffRequest(client))

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
