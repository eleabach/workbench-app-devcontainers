package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// MCP Protocol structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type CallToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Tool definition
var tools = []Tool{
	{
		Name:        "get_plant_virus_data",
		Description: "Retrieve viral activity levels from wastewater surveillance. Searches across plant name, city, and state to return matching viral activity data.\n\nData Schema:\n- city, state: Location (may be null if metadata missing)\n- plant_name: Wastewater treatment plant name\n- virus: Viral target (e.g., SARS-CoV-2, Influenza A, RSV, Norovirus)\n- level: Activity score (0-1 float) - computed via sliding window anomaly detection comparing recent concentrations to plant-specific historical baseline. 0 = normal activity OR insufficient historical data (<36 samples), 1 = high anomalous activity detected\n- most_recent_date: Latest sample collection date\n\nThe level score represents the proportion of recent measurements showing elevated activity (>1 std dev above baseline). Higher scores indicate statistically significant increases in viral concentration vs historical patterns for that specific plant-virus pair.",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search term to filter results. Searches across plant name (e.g., 'North County WWTP'), city (e.g., 'San Diego'), or state (e.g., 'California' or 'CA'). Case-insensitive matching.",
				},
			},
			Required: []string{"query"},
		},
	},
}

// getPlantVirusData reads the CSV and returns virus data matching the query
// Query searches across plant_name, city, and state fields
func getPlantVirusData(query string) (string, error) {
	file, err := os.Open("test.csv")
	if err != nil {
		return "", fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read CSV file: %v", err)
	}

	if len(records) < 2 {
		return "", fmt.Errorf("CSV file is empty or has no data rows")
	}

	// Find all matching entries (case-insensitive search across city, state, plant_name)
	// CSV columns: city, state, plant_name, virus, level, most_recent_date
	queryLower := strings.ToLower(strings.TrimSpace(query))
	plantGroups := make(map[string][]string) // plant_name -> list of result lines
	plantLocations := make(map[string]string) // plant_name -> location string

	for i, record := range records {
		if i == 0 {
			// Skip header row
			continue
		}
		if len(record) < 6 {
			continue
		}

		city := strings.TrimSpace(record[0])
		state := strings.TrimSpace(record[1])
		plantName := strings.TrimSpace(record[2])
		virus := strings.TrimSpace(record[3])
		level := strings.TrimSpace(record[4])
		date := strings.TrimSpace(record[5])

		// Check if query matches city, state, or plant_name
		if strings.Contains(strings.ToLower(city), queryLower) ||
			strings.Contains(strings.ToLower(state), queryLower) ||
			strings.Contains(strings.ToLower(plantName), queryLower) {

			// Build location string for this plant (only once per plant)
			if _, exists := plantLocations[plantName]; !exists {
				if city != "" && state != "" {
					plantLocations[plantName] = fmt.Sprintf("%s, %s", city, state)
				} else if city != "" {
					plantLocations[plantName] = city
				} else if state != "" {
					plantLocations[plantName] = state
				} else {
					plantLocations[plantName] = "Location unknown"
				}
			}

			// Add virus record
			plantGroups[plantName] = append(plantGroups[plantName],
				fmt.Sprintf("  - Virus: %s | Activity Level: %s | Date: %s", virus, level, date))
		}
	}

	if len(plantGroups) == 0 {
		return "", fmt.Errorf("no data found matching '%s'", query)
	}

	// Format output
	var output strings.Builder
	if len(plantGroups) > 1 {
		output.WriteString(fmt.Sprintf("Found %d plants matching '%s':\n\n", len(plantGroups), query))
	}

	for plantName, results := range plantGroups {
		output.WriteString(fmt.Sprintf("Plant: %s\n", plantName))
		output.WriteString(fmt.Sprintf("Location: %s\n\n", plantLocations[plantName]))
		output.WriteString(strings.Join(results, "\n"))
		output.WriteString("\n\n")
	}

	return strings.TrimSpace(output.String()), nil
}

// handleCallTool processes tool calls
func handleCallTool(params CallToolParams) CallToolResult {
	switch params.Name {
	case "get_plant_virus_data":
		query, ok := params.Arguments["query"].(string)
		if !ok {
			return CallToolResult{
				Content: []ContentItem{{Type: "text", Text: "Error: 'query' parameter required"}},
				IsError: true,
			}
		}

		output, err := getPlantVirusData(query)
		if err != nil {
			return CallToolResult{
				Content: []ContentItem{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}
		}

		return CallToolResult{
			Content: []ContentItem{{Type: "text", Text: output}},
			IsError: false,
		}

	default:
		return CallToolResult{
			Content: []ContentItem{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", params.Name)}},
			IsError: true,
		}
	}
}

// handleRequest processes JSON-RPC requests
func handleRequest(req JSONRPCRequest) JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: InitializeResult{
				ProtocolVersion: "2024-11-05",
				Capabilities:    map[string]interface{}{"tools": map[string]interface{}{}},
				ServerInfo:      ServerInfo{Name: "sightline-mcp-server", Version: "1.0.0"},
			},
		}
	case "notifications/initialized":
		// Client sends this notification after receiving initialize response
		// No response needed for notifications
		return JSONRPCResponse{}
	case "tools/list":
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  ListToolsResult{Tools: tools},
		}
	case "tools/call":
		var params CallToolParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &RPCError{Code: -32602, Message: "Invalid params"},
			}
		}
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  handleCallTool(params),
		}
	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32601, Message: "Method not found"},
		}
	}
}

func main() {
	fmt.Fprintln(os.Stderr, "Sightline MCP Server v1.0 starting...")
	fmt.Fprintf(os.Stderr, "Ready - %d tools available\n", len(tools))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			continue
		}

		response := handleRequest(req)
		// Only send response if there's a result or error (skip empty responses for notifications)
		if response.Result != nil || response.Error != nil {
			responseBytes, _ := json.Marshal(response)
			fmt.Println(string(responseBytes))
		}
	}
}
