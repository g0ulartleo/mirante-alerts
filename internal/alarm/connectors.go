package alarm

type Connector struct {
	ID     string
	Type   string
	Config map[string]any
}

type MySQLConnectorConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Tunnel   TunnelConfig
}

type TunnelConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	PemFile  string
}
