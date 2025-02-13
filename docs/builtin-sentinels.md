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

