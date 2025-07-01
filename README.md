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

2. **Initial Setup**
   ```bash
   make setup
   ```
   This creates a sample `.env` file with all available configuration options and necessary directories.

3. **Configure Environment Variables**

   Edit the `.env` file created by the setup command. The following variables are available:
   - `REDIS_ADDR` (default: `127.0.0.1:6379`)
   - `DB_DRIVER` (set to `redis` or `mysql` or `sqlite`)
   - `API_KEY`
   - For MySQL storage:
     - `MYSQL_DB_HOST`
     - `MYSQL_DB_PORT`
     - `MYSQL_DB_USER`
     - `MYSQL_DB_PASSWORD`
   - For email notifications:
     - `SMTP_HOST`
     - `SMTP_PORT`
     - `SMTP_USER`
     - `SMTP_PASSWORD`
   - For HTTP server:
     - `HTTP_ADDR` (default: `127.0.0.1`)
     - `HTTP_PORT` (default: `40169`)
   - For OAuth authentication:
     - `OAUTH_CLIENT_ID`
     - `OAUTH_CLIENT_SECRET`
     - `OAUTH_JWT_SECRET`
   - For dashboard basic auth (optional):
     - `DASHBOARD_BASIC_AUTH_USERNAME`
     - `DASHBOARD_BASIC_AUTH_PASSWORD`

4. **Install Dependencies**
   ```bash
   go mod download
   ```

## Authentication

mirante-alerts supports two authentication methods:

### OAuth Authentication (Recommended for production/multi-user setups)

OAuth authentication allows you to control access using your existing Google or GitHub accounts, with fine-grained control over who can access your monitoring system.

#### Setting up OAuth

1. **Initialize OAuth Configuration**
   ```bash
   make init-oauth
   ```
   This creates a sample configuration file at `config/auth.yaml`.

2. **Configure OAuth Provider**

   **For Google OAuth:**
   - Go to [Google Cloud Console](https://console.developers.google.com/)
   - Create a new project or select existing
   - Enable Google+ API
   - Create OAuth 2.0 credentials
   - Set authorized redirect URI to: `http://your-domain:40169/auth/callback`

   **For GitHub OAuth:**
   - Go to [GitHub OAuth Apps](https://github.com/settings/applications/new)
   - Create a new OAuth App
   - Set Authorization callback URL to: `http://your-domain:40169/auth/callback`

3. **Update Configuration**

   **First, configure OAuth secrets in your `.env` file:**
   ```bash
   # OAuth Configuration
   OAUTH_CLIENT_ID=your-oauth-client-id
   OAUTH_CLIENT_SECRET=your-oauth-client-secret
   OAUTH_JWT_SECRET=your-secure-jwt-secret-key
   ```

   **Then, edit `config/auth.yaml` for non-sensitive settings:**
   ```yaml
   oauth:
     enabled: true
     provider: "google"  # or "github"
     redirect_url: "http://your-domain:40169/auth/callback"
     allowed_domains:
       - "@yourcompany.com"
       - "@contractor.yourcompany.com"
     allowed_emails:
       - "admin@yourcompany.com"
       - "developer@yourcompany.com"
     session_timeout: "24h"
   ```

4. **CLI Authentication**
   ```bash
   ./bin/cli auth http://your-domain:40169
   ```
   This will open your browser for authentication and save the token locally.

#### Access Control

You can control access using two methods:

- **Domain-based:** Allow all users with emails from specific domains
  ```yaml
  allowed_domains:
    - "@yourcompany.com"
    - "@contractor.yourcompany.com"
  ```

- **Email-based:** Allow specific individual email addresses
  ```yaml
  allowed_emails:
    - "john@company.com"
    - "jane@company.com"
  ```

### Basic API Key Authentication

For a simple API-Key setup

```bash
./bin/cli auth-key <your_endpoint> <api_key>
```

Set the `API_KEY` environment variable on your server.

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

   Start with setting up authentication, and then using `help` to see the available commands
   ```bash
   # For OAuth authentication
   $ ./bin/cli auth <your_endpoint>

   # Or for API key authentication
   $ ./bin/cli auth-key <your_endpoint> <api_key>

   $ ./bin/cli help
   ```

## Architecture


### System Components

- **HTTP Server:** Serves the web UI that displays alarm status and history, and an admin API for CLI usage. Located in `cmd/http-server/`.
- **Worker Server:** Processes background tasks such as writing signals and executing sentinel checks. See `cmd/worker-server/`.
- **Scheduler:** Registers and executes periodic sentinel checks as well as cleanup tasks. Located in `cmd/scheduler/`.
- **CLI:** A command-line interface for managing alarms and signals. See `cmd/cli/`.

### Built-in Sentinels

- **EndpointChecker**: Performs HTTP operations on URLs and validates responses based on configuration
- **MySQLCountChecker**: Executes SQL queries that return counts and validates them against expected values
- **SQSCountChecker**: Monitors the number of messages in an SQS queue and alerts if it exceeds a threshold
- See all built-in sentinels with configuration examples [here](docs/builtin-sentinels.md)


## License

Mirante Alerts is distributed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.
