# mirante-alerts
mirante-alerts is an open-source, lightweight monitoring system designed to watch over multiple projects and external services, providing simple red/green status indicators based on the health of its sentinels.

## Alerts
Alerts use sentinels to check for specific aspects of your systems. Each different type of sentinel implements a specific monitoring strategy.

### Built-in Sentinels
- **EndpointValidator**: Performs HTTP operations on URLs and validates responses based on configuration
- See all built-in sentinels with configuration examples [here](docs/builtin-sentinels.md)

### Custom Sentinels
You can create your own sentinel types by implementing the Sentinel interface. Check out the [custom-sentinels](docs/custom-sentinels.md) documentation for details.

### Adding a new alert
Simply create a new yaml file in the `alerts` directory. The path of the file be reflected in the URL of the alert.
