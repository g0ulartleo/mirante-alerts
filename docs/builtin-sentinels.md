## Built-in Sentinels

### Endpoint Checker

The Endpoint Checker sentinel type performs HTTP operations on URLs and validates responses based on configuration.

#### Configuration

```yaml
id: providers-apis-google-health-check
name: Google Health Check
type: endpoint-checker
config:
  url: https://example.com
  expected_status: 200
  expected_body: "Hello, World!" # optional
```

### MySQL Count Checker

The MySQL Count Checker sentinel type executes a SQL query that returns a count and validates it against an expected value.

#### Configuration

```yaml
id: users-count-check
name: Users Count Check
type: mysql-count-checker
config:
  connection:
    host: localhost
    port: 3306
    user: root
    password: secret
    database: myapp
  query: "SELECT COUNT(*) FROM users"
  expected: 100
```

### SQS Count Checker

The SQS Count Checker sentinel type monitors the number of messages in an Amazon SQS queue and alerts if it exceeds a specified threshold.

#### Configuration

```yaml
id: queue-backlog-monitor
name: SQS Queue Message Count Monitor
type: sqs-count-checker
config:
  queue_url: https://sqs.us-east-1.amazonaws.com/123456789012/my-queue
  max_message_count: 1000
  aws_region: us-east-1
```
