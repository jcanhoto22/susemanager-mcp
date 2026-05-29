# MCP SUSE Manager Server

MCP (Model Context Protocol) server for SUSE Manager / Uyuni integration. Enables AI agents to manage SUSE Manager infrastructure through standard MCP tools.

Built with Go and the [mcp-golang](https://github.com/metoro-io/mcp-golang) library.

> **Disclaimer:** This project is provided as-is. Not all API functions have been fully tested against a live SUSE Manager instance. Use at your own risk.

## Features

### Authentication
- Connect, disconnect, test connection, get API version, check permissions

### System Management
- List all systems, get detailed information, search by hostname
- Hardware inventory (CPU, memory, DMI, devices)
- Software inventory (installed packages)
- Event history with pagination
- Schedule reboots, package updates, highstate, apply states

### Security / CVE
- List systems affected by CVEs with patch status
- Get relevant errata for systems
- Schedule errata application
- Get unscheduled errata

### Channel Management
- List all software channels, get channel details
- List channel packages
- View and change system channel subscriptions
- Create and delete custom channels
- Create and remove repositories
- Associate repositories to channels

### Content Lifecycle Management
- List, create, and remove lifecycle projects
- Build and promote projects through environments
- Create environments within a project
- Manage sources (attach, detach, list, replace)

### Resources
- `suma://systems` — all systems
- `suma://channels` — all channels
- `suma://summary` — dashboard summary

### Prompts
- `security-audit` — audit CVEs and affected systems
- `patch-status` — summary of patch status
- `system-overview` — comprehensive system overview

## Architecture

```
cmd/server/main.go          Entry point, wiring
internal/
  config/                   YAML + env config
  logger/                   slog-based logger
  api/                      API client (auth, system, audit, channel, errata)
  mcp/                      MCP server, tools, resources, prompts
  models/                   Go structs for API types
pkg/xmlrpc/                 XML-RPC transport with retry and TLS
```

### Flow

1. Server starts, loads config, authenticates to SUSE Manager
2. MCP server starts on stdio or HTTP transport
3. AI agent connects and discovers available tools
4. Each tool call marshals arguments, calls the SUSE Manager XML-RPC API, returns structured JSON

### Authentication

- Login on startup, automatic re-login on session expiry
- Session key is managed internally, never exposed to the client
- Passwords are never logged

## Quick Start

### Prerequisites

- Go 1.22+
- Access to a SUSE Manager 4.3+ / Uyuni instance

### Install

```bash
git clone <repo-url> mcp-susemanager
cd mcp-susemanager
make build
```

### Configure

Copy and edit the config file:

```bash
cp config.yaml.example config.yaml
```

```yaml
suse:
  url: "https://suse-manager.example.com/rpc/api"
  username: "admin"
  password: "your-password"
  insecure_skip_verify: false
  timeout: 60

server:
  port: 8080
  log_level: info
  transport: stdio

logging:
  format: text    # text or json
  level: info
```

Or use environment variables:

```bash
export SUSE_URL="https://suse-manager.example.com/rpc/api"
export SUSE_USERNAME="admin"
export SUSE_PASSWORD="your-password"
export SURE_INSECURE="false"
```

### Run

```bash
make run
```

Or directly:

```bash
./bin/mcp-susemanager -config config.yaml
```

## Docker

```bash
# Build
make docker-build

# Run with config file
docker run --rm -i -v $(pwd)/config.yaml:/etc/suse-mcp/config.yaml mcp-susemanager:latest

# Or with docker-compose
docker-compose up
```

### HTTP transport mode

```bash
SERVER_TRANSPORT=http docker-compose up
```

## MCP Tools (43)

### Auth (4)

| Tool | Description |
|---|---|
| `suse_connect` | Test connection to SUSE Manager |
| `suse_get_version` | Get API version |
| `suse_check_permissions` | Check user permissions |
| `suse_disconnect` | Logout |

### Systems (10)

| Tool | Description |
|---|---|
| `suse_list_systems` | List all systems |
| `suse_get_system_details` | Get system details by ID |
| `suse_search_system_by_hostname` | Search systems by regex hostname |
| `suse_get_system_hardware` | Get CPU, memory, DMI, devices |
| `suse_get_system_software` | List installed packages |
| `suse_get_system_events` | Get event history |
| `suse_schedule_reboot` | Schedule a reboot |
| `suse_schedule_package_update` | Schedule package update |
| `suse_schedule_highstate` | Schedule highstate |
| `suse_apply_states` | Apply Salt states |

### CVE / Security (4)

| Tool | Description |
|---|---|
| `suse_list_cve_systems` | List systems by CVE status |
| `suse_get_system_errata` | Get relevant errata |
| `suse_schedule_errata` | Apply errata to systems |
| `suse_get_unscheduled_errata` | Get unscheduled errata |

### Channels (13)

| Tool | Description |
|---|---|
| `suse_list_channels` | List software channels |
| `suse_get_channel_details` | Get channel details |
| `suse_list_channel_packages` | List channel packages |
| `suse_list_system_channels` | List system channel subscriptions |
| `suse_change_system_channels` | Change system channels |
| `suse_create_channel` | Create a new software channel |
| `suse_delete_channel` | Delete a custom software channel |
| `suse_list_arches` | List available channel architectures |
| `suse_create_repo` | Create a new repository |
| `suse_remove_repo` | Remove a repository |
| `suse_list_repos` | List all repositories |
| `suse_list_channel_repos` | List repos associated with a channel |
| `suse_associate_repo_to_channel` | Associate a repository with a channel |

### Content Lifecycle (12)

| Tool | Description |
|---|---|
| `suse_list_projects` | List all lifecycle projects |
| `suse_lookup_project` | Get project details by label |
| `suse_create_project` | Create a new lifecycle project |
| `suse_remove_project` | Remove a lifecycle project |
| `suse_build_project` | Build a lifecycle project |
| `suse_promote_project` | Promote a project environment |
| `suse_list_project_environments` | List project environments |
| `suse_create_environment` | Create an environment in a project |
| `suse_attach_source` | Attach a source (channel) to a project |
| `suse_detach_source` | Detach a source from a project |
| `suse_list_project_sources` | List sources attached to a project |
| `suse_set_sources` | Replace all sources of a given type |

## Client Configuration

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "suse-manager": {
      "command": "/path/to/mcp-susemanager/bin/mcp-susemanager",
      "args": ["-config", "/path/to/config.yaml"],
      "env": {}
    }
  }
}
```

### Docker with Claude Desktop

```json
{
  "mcpServers": {
    "suse-manager": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "-v", "/path/to/config.yaml:/etc/suse-mcp/config.yaml", "mcp-susemanager:latest"],
      "env": {}
    }
  }
}
```

### LM Studio

LM Studio supports MCP over HTTP transport. Start the server in HTTP mode:

```bash
make build
SERVER_TRANSPORT=http ./bin/mcp-susemanager -config config.yaml
```

Configure in LM Studio settings:

```json
{
  "mcpServers": {
    "suse-manager": {
      "url": "http://localhost:8080/mcp",
      "type": "http"
    }
  }
}
```

### Cursor

Add to Cursor MCP configuration:

```json
{
  "mcpServers": {
    "suse-manager": {
      "command": "/path/to/bin/mcp-susemanager",
      "args": ["-config", "/path/to/config.yaml"]
    }
  }
}
```

### VS Code (via Continue)

Add to `.continuerc.json`:

```json
{
  "experimental": {
    "mcpServers": {
      "suse-manager": {
        "command": "/path/to/bin/mcp-susemanager",
        "args": ["-config", "/path/to/config.yaml"]
      }
    }
  }
}
```

## Example Queries

Once connected, you can ask your AI agent:

> "List all systems in SUSE Manager"

> "Show me details for system ID 100"

> "Which systems are affected by CVE-2025-1234?"

> "Schedule a reboot for server web-01"

> "Apply all security errata to systems 101, 102, 103"

> "Run highstate on system 100"

> "List all software channels"

> "What packages are installed on system 100?"

> "Show me the hardware inventory for system 100"

> "Change the base channel of system 100 to sles15-sp5-updates"

## Available Prompts

### security-audit

```
Analyze the following data and provide a structured security assessment:
1. List all systems and their CVE patch status
2. Identify affected systems
3. Provide prioritized remediation recommendations
```

### patch-status

```
Analyze the patch status and provide a summary:
1. List all systems and their patch status
2. Identify out-of-date systems
3. Identify systems needing reboots
4. Provide a prioritized patching plan
```

### system-overview

```
Provide a comprehensive overview:
1. List all registered systems
2. For each system: hardware, packages, channels, events
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
# or
make test-short
# with race detector
make test-race
```

### Lint

```bash
make vet
# or with golangci-lint:
make lint
```

### Clean

```bash
make clean
```

## Configuration Reference

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `SUSE_URL` | SUSE Manager API endpoint | `https://localhost/rpc/api` |
| `SUSE_USERNAME` | API username | `admin` |
| `SUSE_PASSWORD` | API password | `""` |
| `SUSE_INSECURE` | Skip TLS verification | `false` |
| `SUSE_TIMEOUT` | API timeout in seconds | `60` |
| `SERVER_PORT` | MCP HTTP port | `8080` |
| `SERVER_TRANSPORT` | Transport type: stdio or http | `stdio` |
| `SERVER_LOG_LEVEL` | Log level | `info` |
| `LOG_LEVEL` | Log level (overrides server) | `info` |
| `LOG_FORMAT` | Log format: text or json | `text` |

## Project Structure

```
cmd/server/main.go              Entrypoint
internal/
  api/
    client.go                   API client with session management
    system.go                   System API methods
    audit.go                    CVE/audit API methods
    channel.go                   Channel API methods
    content_lifecycle.go         Content Lifecycle API methods
    errata.go                    Errata API methods
  config/
    config.go                   YAML + env config loader
  logger/
    logger.go                   slog logger setup
  mcp/
    server.go                   MCP server (stdio + HTTP)
    tools.go                    43 tool registrations + handlers
    resources.go                3 resource registrations
    prompts.go                  3 prompt registrations
  models/
    common.go                   Shared types
    system.go                   System structs
    channel.go                  Channel structs
    cve.go                      CVE/audit structs
pkg/
  xmlrpc/
    client.go                   XML-RPC transport with retry
tests/
  integration/                  Integration tests
config.yaml.example             Example config
Dockerfile                      Multi-stage Docker build
docker-compose.yml              Docker Compose
Makefile                        Build automation
```

## API Coverage

### Current (43 tools)

- [x] Auth (4 tools)
- [x] Systems (10 tools)
- [x] CVE / Security (4 tools)
- [x] Channels (13 tools)
- [x] Content Lifecycle (12 tools)
- [x] Resources (3)
- [x] Prompts (3)

### Future

- [ ] Users (create, update, delete, roles)
- [ ] Groups (create, delete, manage members)
- [ ] Activation Keys (create, associate channels)
- [ ] Salt (pillars, grains, formulas)
- [ ] Scheduling (jobs, history)
- [ ] Image Build (profiles, builds)
- [ ] Configuration Management (config channels)
- [ ] Organizations (multi-org)
- [ ] Reporting (CSV, JSON export)
- [ ] Ansible integration

## Troubleshooting

### Connection refused

Check that the SUSE Manager API URL is correct and accessible:

```bash
curl -k https://your-suma-server/rpc/api
```

### Authentication failed

Verify credentials and check the user has API access in SUSE Manager.

### Session expires

The server automatically re-authenticates on session expiry. No action needed.

### Self-signed certificate

Set `insecure_skip_verify: true` in config or `SUSE_INSECURE=true` env var.

### XML-RPC errors

Enable debug logging:

```yaml
logging:
  level: debug
```

## License

MIT
