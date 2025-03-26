# mirante-alerts

mirante-alerts is an open-source monitoring system designed to watch over multiple projects and external services, providing simple red/green status indicators based on the health of its alarms. It features a web UI for real-time monitoring and a CLI for management.

## Current Status
This project was created initally for learning purposes and it is still under development, so production usage is not yet recommended.
If you have been using it, and feel that this is useful, consider contribuiting :) 

## Getting Started

### Prerequisites

- **Go:** The project is built with Go (see `go.mod` for version, currently Go 1.23.6).
- **Redis:** Required for task queue management and main alarm storage.
- **Docker (Optional):** For running via Docker Compose.

### Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/g0ulartleo/mirante-alerts.git
   cd mirante-alerts
   ```

2. **Configure Environment Variables**

   Create a `.env` file or set environment variables as necessary. The following variables are available:
   - `REDIS_ADDR` (default: `127.0.0.1:6379`)
   - `DB_DRIVER` (set to `redis` or `mysql` or `sqlite`)
   - `API_KEY`
   - For MySQL storage:
     - `DB_HOST`
     - `DB_PORT`
     - `DB_USER`
     - `DB_PASSWORD`
   - For email notifications:
     - `SMTP_HOST`
     - `SMTP_PORT`
     - `SMTP_USER`
     - `SMTP_PASSWORD`

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Setting Up Alarms**

   Alarms are configured via YAML files in the `config/alarms` directory. The directory structure reflects the URL path for an alarm's dashboard (if no `path` is defined for the alarm).

   Example alarm configuration:
   ```yaml
   id: my-alarm
   name: My Custom Alarm
   description: "Expects a 200 status code from some API"
   type: endpoint-checker
   interval: "30s"        # Alternatively, specify a cron expression in the `cron` field
   path: ['Project', 'APIs']
   config:
     url: "https://example.com"
     expected_status: 200
   notifications:
      email:
         to:
            - "test@example.com"
            - "test2@example.com"
      slack:
         webhook_url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXX"
   ```

   If you are hosting mirante in your servers, you can also manage alarms using the CLI.

   Start with setting up config, and then using `help` to see the available commands
   ```bash
   $ ./bin/cli config <your_endpoint> <api_key>
   $ ./bin/cli help
   ```

## Architecture

### Core Components

- **Alarms**: The main monitoring unit that defines what to check and how often
- **Sentinels**: The actual monitoring strategies that perform health checks
- **Signals**: The results of health checks, stored for historical tracking

### System Components

- **HTTP Server:** Serves the web UI that displays alarm status and history, and an admin API for CLI usage. Located in `cmd/http-server/`.
- **Worker Server:** Processes background tasks such as writing signals and executing sentinel checks. See `cmd/worker-server/`.
- **Scheduler:** Registers and executes periodic sentinel checks as well as cleanup tasks. Located in `cmd/scheduler/`.
- **CLI:** A command-line interface for managing alarms and signals. See `cmd/cli/`.

### Built-in Sentinels

- **EndpointChecker**: Performs HTTP operations on URLs and validates responses based on configuration
- **MySQLCountChecker**: Executes SQL queries that return counts and validates them against expected values
- See all built-in sentinels with configuration examples [here](docs/builtin-sentinels.md)


## License

Mirante Alerts is distributed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.
