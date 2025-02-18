## Implementing a custom sentinel

 1. Create new file in `internal/sentinel/custom/`:

```go
// internal/sentinel/custom/website_uptime_checker.go
package custom

import (
    "github.com/g0ulartleo/mirante-alerts/internal/sentinel"
    "github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type WebsiteUptimeChecker struct {
    url string
    expectedStatus int
}

func (w WebsiteUptimeChecker) Check(ctx context.Context, alertID string) (signal.Signal, error) {
    // do something
    return signal.Signal{
        AlertID: alertID,
        Status: signal.StatusHealthy,
        Message: "Website is up",
        Timestamp: time.Now(),
    }, nil
}

func (w WebsiteUptimeChecker) Configure(config map[string]interface{}) error {
    if url, ok := config["url"]; !ok {
        return fmt.Errorf("url is required")
    } 
    if expectedStatus, ok := config["expected_status"]; !ok {
        return fmt.Errorf("expected_status is required")
    }
    w.url = url.(string)
    w.expectedStatus = expectedStatus.(int)
    return nil
}

func init() {
    sentinel.Factory.Register("website-uptime-checker", WebsiteUptimeChecker{})
}
```

3. Build the project with `custom` tag:

```bash
go build -tags custom ./cmd/worker-server
```

4. Create an alert configuration file using the sentinel:

```yaml
# alerts/my-website/uptime.yaml
id: my-website-uptime
type: website-uptime-checker
interval: 1m
config:
    url: https://example.com
    expected_status: 200
```
