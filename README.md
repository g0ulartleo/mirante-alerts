# mirante-alerts
mirante-alerts is an open-source, lightweight monitoring system designed to watch over multiple projects and external services, providing simple red/green status indicators based on the health of its alerts.

## Features

- **Modular Alerts**: Configure alerts via YAML files that define the monitoring strategy, schedule (using intervals or cron expressions), and notification settings.
- **Tasks Management**: Uses [Asynq](https://github.com/hibiken/asynq) for reliable task processing (for signal writing, checking alerts, etc).
- **Dashboard**: A dashboard with hierarchical view of alerts.
- **Flexible Signal Storage**: Supports MySQL or Redis for storing signals.


## Alerts
Alerts use sentinels to check for specific aspects of your systems. Each different type of sentinel implements a specific monitoring strategy.

### Built-in Sentinels
- **EndpointValidator**: Performs HTTP operations on URLs and validates responses based on configuration
- See all built-in sentinels with configuration examples [here](docs/builtin-sentinels.md)

### Custom Sentinels

Create new sentinel types by implementing the Sentinel interface:

```go
type Sentinel interface {
    Check(ctx context.Context, alertID string) (Signal, error)
    Configure(config map[string]interface{}) error
}
```

and then registering it with the sentinel factory (TBD):

```go
sentinel.Register("my-sentinel", MySentinel{})
```

See the [custom-sentinels](docs/custom-sentinels.md) documentation for details. (TBD)


### Adding a new alert
Simply create a new yaml file in the `alerts` directory. The path of the file be reflected in the URL of the alert.



## Components

- **HTTP Server:** Serves the web UI that displays alert status and history. Located in `cmd/http-server/`.
- **Worker Server:** Processes background tasks such as writing signals and executing sentinel checks. See `cmd/worker-server/`.
- **Scheduler:** Registers and executes periodic sentinel checks as well as cleanup tasks. Located in `cmd/scheduler/`.

## Getting Started

### Prerequisites

- **Go:** The project is built with Go (see `go.mod` for version, currently Go 1.23.6).
- **Redis:** Required for task queue management.
- **MySQL:** Required if you choose to use MySQL for signal storage (use the `mysql` driver build tag when needed).
- **Docker (Optional):** For running Redis and MySQL via Docker Compose.

### Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/g0ulartleo/mirante-alerts.git
   cd mirante-alerts
   ```

2. **Configure Environment Variables**

   Create a `.env` file or set environment variables as necessary. Some important variables include:
   - `REDIS_ADDR` (default: `127.0.0.1:6379`)
   - `DB_DRIVER` (set to `mysql` or `memory`)
   - For MySQL storage:
     - `DB_HOST`
     - `DB_PORT`
     - `DB_USER`
     - `DB_PASSWORD`

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Setting Up Alerts**

   Alerts are configured via YAML files in the `alerts` directory. The directory structure reflects the URL path for an alert’s dashboard.

   Example alert configuration:
   ```yaml
   id: my-alert
   name: My Custom Alert
   type: endpoint-checker
   config:
     url: "https://example.com"
     expected_status: 200
     expected_body: "OK"  # Optional
   interval: "30s"        # Alternatively, specify a cron expression in the `cron` field
   ```

## License

Mirante Alerts is distributed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.
