#!/bin/bash

# Readable test script for Sightline MCP Server

echo "=========================================="
echo "Sightline MCP Server - Test Suite"
echo "=========================================="
echo

# Function to extract and display tool results
parse_result() {
    local response="$1"
    local test_name="$2"

    echo "[$test_name]"
    echo "$response" | jq -r '.result.content[0].text // .result.tools // .result.serverInfo // "No content"' 2>/dev/null
    echo
}

# Run the server and capture output
OUTPUT=$(
    (
        echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}'
        echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
        echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_plant_virus_data","arguments":{"query":"North County WWTP"}}}'
        echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_plant_virus_data","arguments":{"query":"CA"}}}'
        echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_plant_virus_data","arguments":{"query":"Los Angeles"}}}'
        echo '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"get_plant_virus_data","arguments":{"query":"Eastside Treatment Facility"}}}'
        echo '{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"get_plant_virus_data","arguments":{"query":"Invalid Query"}}}'
    ) | ./sightline-mcp-server 2>/dev/null
)

# Parse each response
echo "Test 1: Initialize"
echo "------------------------------------------"
echo "$OUTPUT" | sed -n '1p' | jq -r '.result.serverInfo | "Server: \(.name) v\(.version)\nProtocol: \(.protocolVersion // "N/A")"' 2>/dev/null
echo "$OUTPUT" | sed -n '1p' | jq -r '.result | "Protocol Version: \(.protocolVersion)"' 2>/dev/null
echo

echo "Test 2: List Available Tools"
echo "------------------------------------------"
echo "$OUTPUT" | sed -n '2p' | jq -r '.result.tools[] | "Tool: \(.name)\nDescription: \(.description)\n"' 2>/dev/null
echo

echo "Test 3: Query by Plant Name 'North County WWTP'"
echo "------------------------------------------"
echo "$OUTPUT" | sed -n '3p' | jq -r '.result.content[0].text' 2>/dev/null
echo

echo "Test 4: Query by State 'CA' (Multiple Plants)"
echo "------------------------------------------"
echo "$OUTPUT" | sed -n '4p' | jq -r '.result.content[0].text' 2>/dev/null
echo

echo "Test 5: Query by City 'Los Angeles'"
echo "------------------------------------------"
echo "$OUTPUT" | sed -n '5p' | jq -r '.result.content[0].text' 2>/dev/null
echo

echo "Test 6: Query by Plant Name 'Eastside Treatment Facility'"
echo "------------------------------------------"
echo "$OUTPUT" | sed -n '6p' | jq -r '.result.content[0].text' 2>/dev/null
echo

echo "Test 7: Query Invalid (Error Test)"
echo "------------------------------------------"
echo "$OUTPUT" | sed -n '7p' | jq -r '.result.content[0].text' 2>/dev/null
IS_ERROR=$(echo "$OUTPUT" | sed -n '7p' | jq -r '.result.isError' 2>/dev/null)
echo "[Error returned: $IS_ERROR]"
echo

echo "=========================================="
echo "All tests complete!"
echo "=========================================="
