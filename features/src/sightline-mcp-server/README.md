# Sightline MCP Server

Mock MCP server that provides wastewater viral surveillance data for testing and development.

## What It Does

This is a test implementation that exposes viral activity data from wastewater treatment plants via the Model Context Protocol (MCP). It uses mock CSV data and requires no authentication or external dependencies.

## Installation

Add to your `devcontainer.json`:

```json
{
  "features": {
    "ghcr.io/verily-src/workbench-app-devcontainers/sightline-mcp-server:latest": {}
  }
}
```

Rebuild your devcontainer. The server installs at `/opt/wb-mcp-server/wb-mcp-server`.

## Setup

### With Claude CLI

```bash
claude mcp add --transport stdio sightline -- /opt/wb-mcp-server/wb-mcp-server
```

### With Gemini CLI

```bash
gemini mcp add --scope user sightline /opt/wb-mcp-server/wb-mcp-server
```

## Available Tool

### `get_wastewater_surveillance_data`

Retrieves viral activity levels from wastewater surveillance.

**Parameters:**
- `query` (string, required): Search term to filter results. Searches across plant name, city, or state (case-insensitive).

**Data Schema:**
- `city`: City location
- `state`: State location
- `plant_name`: Wastewater treatment plant name
- `virus`: Viral target (e.g., SARS-CoV-2, Influenza A, RSV, Norovirus)
- `level`: **BINARY** viral activity level
  - **0 = NOT HIGH** (normal/baseline activity or insufficient data)
  - **1 = HIGH** (elevated viral activity detected)
- `most_recent_date`: Latest sample collection date

## Example Queries

### Search by City

```
"What viral activity is detected in San Diego wastewater?"
```

Returns all virus data for San Diego treatment plants.

### Search by State

```
"Show me wastewater surveillance data for California"
```

Returns data for all California plants (searches by state code "CA" or full name).

### Search by Plant Name

```
"What viruses are detected at North County WWTP?"
```

Returns all viral activity for the specific treatment plant.

### Search by Virus Type

```
"Show me SARS-CoV-2 levels across all plants"
```

Returns SARS-CoV-2 activity data from all treatment plants.

## Mock Data

The server uses static test data from `/opt/wb-mcp-server/test.csv` containing:
- 6 cities across 5 states
- 6 wastewater treatment plants
- Detection of 5 viruses: SARS-CoV-2, Influenza A, Influenza B, RSV, Norovirus
- Recent data (February 2026)

## How It Works

- **No authentication required** - This is a mock server for testing
- **No external dependencies** - Reads from a local CSV file
- **Simple search** - Case-insensitive substring matching across plant name, city, and state fields

## Troubleshooting

### Server not responding

Test the server directly:
```bash
/opt/wb-mcp-server/wb-mcp-server
```

Then send a test request:
```json
{"jsonrpc":"2.0","id":1,"method":"tools/list"}
```

### "Failed to open CSV file" error

The server expects the CSV file at `/opt/wb-mcp-server/test.csv`. If it's missing, reinstall the devcontainer feature.

## Technical Details

- **Protocol**: Model Context Protocol (MCP) via stdio transport
- **Language**: Go 1.21+
- **Data format**: CSV with headers
- **Server name**: `sightline-mcp-server`
- **Version**: 1.0.0
