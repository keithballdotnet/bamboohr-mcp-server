# BambooHR MCP Server Usage Examples

This document provides examples of how to use the BambooHR MCP Server tools.

## Setup

Before using the tools, make sure you have:
1. Set the `BAMBOOHR_API_KEY` environment variable
2. Set the `BAMBOOHR_COMPANY` environment variable
3. Started the MCP server

## Tool Examples

### 1. List Employees

Get a directory of all employees in your company:

```json
{
  "tool": "list_employees",
  "arguments": {}
}
```

**Expected Response:**
```json
{
  "employees": [
    {
      "id": "123",
      "firstName": "John",
      "lastName": "Doe",
      "jobTitle": "Software Engineer",
      "department": "Engineering"
    }
  ]
}
```

### 2. Get Time-Off Requests

Get all time-off requests for a specific employee:

```json
{
  "tool": "get_time_off_requests",
  "arguments": {
    "employeeId": "123"
  }
}
```

Get time-off requests for a specific date range:

```json
{
  "tool": "get_time_off_requests",
  "arguments": {
    "employeeId": "123",
    "start": "2024-01-01",
    "end": "2024-12-31"
  }
}
```

**Expected Response:**
```json
[
  {
    "id": 456,
    "employeeId": 123,
    "name": "John Doe",
    "start": "2024-06-15",
    "end": "2024-06-16",
    "created": "2024-05-01T10:00:00Z",
    "type": {
      "id": 1,
      "name": "Vacation",
      "icon": "🏖️"
    },
    "amount": {
      "unit": "days",
      "amount": 2.0
    },
    "notes": "Family vacation",
    "status": {
      "status": "approved",
      "lastChanged": "2024-05-02T14:30:00Z",
      "lastChangedByUserId": 789
    }
  }
]
```

### 3. Get Time-Off Balance

Get the current time-off balance for an employee:

```json
{
  "tool": "get_time_off_balance",
  "arguments": {
    "employeeId": "123"
  }
}
```

**Expected Response:**
```json
[
  {
    "policyType": "Vacation",
    "balance": 15.5,
    "used": 4.5,
    "available": 15.5,
    "unit": "days"
  },
  {
    "policyType": "Sick Leave",
    "balance": 8.0,
    "used": 2.0,
    "available": 8.0,
    "unit": "days"
  }
]
```

### 4. Create Time-Off Request

Create a new time-off request for an employee:

```json
{
  "tool": "create_time_off_request",
  "arguments": {
    "employeeId": "123",
    "timeOffTypeId": "1",
    "start": "2024-12-23",
    "end": "2024-12-23",
    "employeeNote": "Personal day off"
  }
}
```

**Expected Response:**
```json
{
  "id": "12345",
  "employeeId": "123",
  "name": "John Doe",
  "start": "2024-12-23",
  "end": "2024-12-23",
  "type": {
    "id": "1",
    "name": "Vacation",
    "icon": "palm-trees"
  },
  "amount": {
    "unit": "days",
    "amount": 1
  },
  "status": {
    "status": "approved",
    "lastChanged": "2024-12-01"
  }
}
```

**Time-Off Type IDs (common values):**
- `"1"` - Vacation
- `"2"` - Sick Days  
- `"5"` - Sick Day Child
- `"19"` - Mobile work from abroad
- `"27"` - Home Office days

**Required Arguments:**
- `employeeId`: The numeric ID of the employee
- `timeOffTypeId`: The ID of the time-off type (see list above)
- `start`: Start date in YYYY-MM-DD format
- `end`: End date in YYYY-MM-DD format

**Optional Arguments:**
- `employeeNote`: A note from the employee about the request

## Common Use Cases

### 1. Check Employee Time-Off Status

1. First, list employees to get their IDs
2. Use the employee ID to check their time-off balance
3. Get their recent time-off requests to see upcoming time off

### 2. Time-Off Planning

1. Check multiple employees' balances
2. Look at upcoming time-off requests
3. Plan coverage and workload distribution

### 3. HR Reporting

1. Get time-off data for all employees
2. Analyze patterns and usage
3. Generate reports for management

### 4. Create Time-Off Requests

1. Find your employee ID using `list_employees`
2. Choose the appropriate time-off type ID
3. Create the request with start/end dates
4. Add optional notes if needed

**Example: Booking Home Office days for all Fridays in October:**
```json
{
  "tool": "create_time_off_request",
  "arguments": {
    "employeeId": "157",
    "timeOffTypeId": "27",
    "start": "2025-10-03",
    "end": "2025-10-03",
    "employeeNote": "Friday home office"
  }
}
```

Repeat for each Friday: Oct 10, 17, 24, 31.

## Error Handling

The tools will return error messages for common issues:

- **Missing API Key**: "BAMBOOHR_API_KEY environment variable is required"
- **Missing Company**: "BAMBOOHR_COMPANY environment variable is required"
- **Invalid Employee ID**: "employeeId must be a valid integer"
- **API Errors**: "API error 404: Employee not found"

## Tips

1. **Employee IDs**: Always use numeric employee IDs, not names or email addresses
2. **Date Formats**: Use YYYY-MM-DD format for start and end dates
3. **Rate Limiting**: The BambooHR API has rate limits, so avoid making too many requests quickly
4. **Permissions**: Make sure your API key has permission to access time-off data

## Troubleshooting

### Common Issues

1. **"API error 401: Unauthorized"**
   - Check that your API key is correct
   - Verify your company subdomain is correct

2. **"API error 403: Forbidden"**
   - Your API key may not have permission to access time-off data
   - Contact your BambooHR administrator

3. **"employeeId must be a valid integer"**
   - Make sure you're using the numeric employee ID, not the name
   - Use the `list_employees` tool to find the correct ID

4. **Empty responses**
   - The employee may not have any time-off requests in the specified date range
   - Try expanding the date range or removing date filters
