## Built-in Sentinels

### Endpoint Checker

The Endpoint Checker sentinel type performs HTTP operations on URLs and validates responses based on configuration.

#### Configuration

```yaml
id: 1
name: my-api-health-check
type: endpoint-checker
config:
  url: https://example.com
  expected_status: 200
  expected_body: "Hello, World!" # optional
```

