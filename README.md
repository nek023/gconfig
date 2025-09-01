# gconfig

Enhanced Google Cloud SDK configuration management tool with interactive selection and bulk operations.

## Features

- **Interactive Switching**: Switch between gcloud configurations using fzf's fuzzy finder
- **Bulk Export**: Export configurations to YAML format for backup or migration
- **Bulk Import**: Import configurations from YAML files to quickly set up environments
- **List Configurations**: Display all available configurations in a formatted table

## Installation

### From Source

```bash
go install github.com/nek023/gconfig@latest
```

### Build Locally

```bash
git clone https://github.com/nek023/gconfig.git
cd gconfig
go build -o gconfig .
```

## Prerequisites

- Go 1.24.5 or higher
- [gcloud CLI](https://cloud.google.com/sdk/docs/install) installed and configured
- [fzf](https://github.com/junegunn/fzf) (required for interactive switching)

## Usage

### Switch Configuration

Interactively switch between gcloud configurations using fzf.

```bash
# Show all configurations and select interactively
gconfig switch

# Start with a search query
gconfig switch prod

# Filter configurations containing specific text
gconfig switch my-project
```

The currently active configuration is highlighted in red for easy identification.

### Export Configurations

Export gcloud configurations to YAML format.

```bash
# Export all configurations
gconfig export

# Export a specific configuration
gconfig export my-config

# Save to a file
gconfig export > configs.yaml

# Export specific configuration to file
gconfig export production > prod-config.yaml
```

### Import Configurations

Import configurations from YAML files.

```bash
# Import from file
gconfig import configs.yaml

# Import from stdin
cat configs.yaml | gconfig import

# Import from URL
curl -s https://example.com/configs.yaml | gconfig import
```

### List Configurations

Display all available gcloud configurations.

```bash
# List all configurations
gconfig list

# Filter output (using standard Unix tools)
gconfig list | grep prod
```

## Examples

### Backup and Restore

```bash
# Backup all configurations
gconfig export > backup-$(date +%Y%m%d).yaml

# Restore configurations
gconfig import backup-20240101.yaml
```

### Migration Between Machines

```bash
# On source machine
gconfig export > my-configs.yaml

# Transfer file to target machine, then:
gconfig import my-configs.yaml
```

### Team Configuration Sharing

```bash
# Export team configurations
gconfig export > team-configs.yaml

# Team members can import
gconfig import team-configs.yaml
```

### Quick Configuration Copy

```bash
# Copy configuration from one machine to another via SSH
gconfig export my-config | ssh user@remote-host 'gconfig import'
```

## Configuration File Format

The YAML format used for import/export follows the gcloud configuration structure:

```yaml
name: my-configuration
properties:
  core:
    account: user@example.com
    project: my-gcp-project
  compute:
    region: us-central1
    zone: us-central1-a
```

Multiple configurations can be stored in a single file using YAML document separators (`---`).

## Commands

| Command | Description |
|---------|-------------|
| `gconfig switch [query]` | Switch configuration interactively |
| `gconfig export [config]` | Export configuration(s) to YAML |
| `gconfig import [file]` | Import configurations from YAML |
| `gconfig list` | List all configurations |
| `gconfig help` | Show help information |

## Tips

- Use `gconfig switch` with partial configuration names for quick filtering
- Combine `gconfig export` with version control to track configuration changes
- Create configuration templates for different environments (dev, staging, prod)
- Use shell aliases for frequently used configurations:
  ```bash
  alias gcprod='gconfig switch production'
  alias gcdev='gconfig switch development'
  ```

## Development

### Building

```bash
go build -o gconfig .
```

### Testing

```bash
go test ./...
```

### Code Formatting

```bash
gofmt -w .
```

## Architecture

gconfig acts as a wrapper around the gcloud CLI, providing enhanced functionality while maintaining compatibility with existing gcloud configurations. It does not store any configuration data itself but operates directly on gcloud's configuration files.

### Key Components

- **cmd/**: Cobra-based CLI commands
- **internal/gcloud/**: Interface for gcloud CLI operations
- **internal/fzf/**: Interactive selection functionality

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

[nek023](https://github.com/nek023)