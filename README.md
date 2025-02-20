# mirante-alerts
mirante-alerts is an open-source, lightweight monitoring system designed to watch over multiple projects and external services, providing simple red/green status indicators based on the health of its alarms.

## Features

- **Modular Alarms**: Configure alarms via YAML files that define the monitoring strategy, schedule (using intervals or cron expressions), and notification settings.
- **Tasks Management**: Uses [Asynq](https://github.com/hibiken/asynq) for reliable task processing (for signal writing, checking alarms, etc).
- **Dashboard**: A dashboard with hierarchical view of alarms.
- **Flexible Signal Storage**: Supports MySQL or Redis for storing signals.

## Alarms
Alarms use sentinels to check for specific aspects of your systems. Each different type of sentinel implements a specific monitoring strategy.

### Built-in Sentinels
- **EndpointValidator**: Performs HTTP operations on URLs and validates responses based on configuration
- See all built-in sentinels with configuration examples [here](docs/builtin-sentinels.md)

### Custom Sentinels

Create new sentinel types by implementing the Sentinel interface:

```go
type Sentinel interface {
    Check(ctx context.Context, alarmID string) (Signal, error)
    Configure(config map[string]interface{}) error
}
```

and then registering it via init function:

```go
func init() {
	sentinel.Factory.Register("my-sentinel", MySentinel{})
}
```

See the [custom-sentinels](docs/custom-sentinels.md) documentation for details.


### Adding a new alarm
Simply create a new yaml file in the `alarms` directory. The path of the file be reflected in the URL of the alarm.


## Components

- **HTTP Server:** Serves the web UI that displays alarm status and history. Located in `cmd/http-server/`.
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

   Alarms are configured via YAML files in the `alarms` directory. The directory structure reflects the URL path for an alarm’s dashboard.

   Example alarm configuration:
   ```yaml
   id: my-alarm
   name: My Custom Alarm
   description: "Expects a 200 status code and a body containing 'OK'"
   type: endpoint-checker
   interval: "30s"        # Alternatively, specify a cron expression in the `cron` field
   config:
     url: "https://example.com"
     expected_status: 200
     expected_body: "OK"  # Optional
   notifications:
      notify_missing_signals: false
      email:
         to: 
            - "test@example.com"
            - "test2@example.com"
      slack:
         webhook_url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXX"
   ```

## License

Mirante Alerts is distributed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.
