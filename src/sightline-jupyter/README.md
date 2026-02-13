# Sightline Jupyter with AI Tools

A Workbench JupyterLab environment with integrated AI assistant tools.

## What's Included

This devcontainer template provides:

- **JupyterLab**: Full scientific Python environment
- **Claude CLI**: Anthropic's Claude Code assistant
- **Gemini CLI**: Google's Gemini AI assistant
- **Sightline MCP Server**: Model Context Protocol server for wastewater viral surveillance data
- **Workbench Tools**: CLI tools for workspace and resource management

## Features

### AI Assistants

Both Claude and Gemini CLI tools are pre-installed and configured with the Sightline MCP server, enabling:

- Code generation and debugging assistance
- Data analysis help
- Access to wastewater viral surveillance data via MCP

### Sightline MCP Server

The integrated MCP server provides the `get_plant_virus_data` tool for querying:
- Viral activity levels at wastewater treatment plants
- SARS-CoV-2, Influenza, RSV, and Norovirus detection data
- Coverage across multiple US cities and states

### Scientific Computing

Built on the Workbench Jupyter image with:
- Python 3 with scientific libraries (NumPy, Pandas, Matplotlib, etc.)
- R kernel support
- JupyterLab extensions

## Usage

### Deploying in Workbench

1. Create a new Cloud Environment in Workbench
2. Select "Sightline Jupyter with AI Tools" template
3. Choose your cloud provider (GCP or AWS)
4. Launch the environment

### Using AI Assistants

Once deployed, open a terminal in JupyterLab:

**Claude CLI:**
```bash
claude
```

**Gemini CLI:**
```bash
gemini
```

Both assistants have access to the Sightline MCP server for viral surveillance queries.

### Example MCP Queries

Ask either assistant:
- "What viral activity is detected in San Diego wastewater?"
- "Show me SARS-CoV-2 levels across all plants"
- "What's the RSV status in Portland?"

## Configuration Options

- **cloud**: Choose `gcp` or `aws` (default: gcp)
- **login**: Auto-login to Workbench CLI - `true` or `false` (default: false)

## Technical Details

- **Base Image**: Workbench Jupyter (pre-built)
- **Container Name**: application-server
- **Port**: 8888 (JupyterLab)
- **User**: jupyter
- **Home Directory**: /home/jupyter
- **Workspace**: /workspace

## Files

- `Dockerfile`: Container image definition
- `docker-compose.yaml`: Service configuration
- `.devcontainer.json`: Devcontainer and feature configuration
- `devcontainer-template.json`: Template metadata for Workbench
