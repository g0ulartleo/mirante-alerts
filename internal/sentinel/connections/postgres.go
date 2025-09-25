package connections

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresConnectionConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	SSLMode     string
	SSLRootCert string
	SSLVerify   bool
}

type PostgresConnection struct {
	DB *sql.DB
}

func NewPostgresConnectionConfig(config map[string]any) (*PostgresConnectionConfig, error) {
	for _, field := range []string{"host", "port", "user", "password", "database"} {
		if _, ok := config[field]; !ok {
			return nil, fmt.Errorf("missing required field: %s", field)
		}
	}

	c := PostgresConnectionConfig{
		Host:     config["host"].(string),
		Port:     getIntValue(config["port"]),
		User:     config["user"].(string),
		Password: config["password"].(string),
		Database: config["database"].(string),
	}

	if sslMode, ok := config["sslmode"].(string); ok {
		c.SSLMode = sslMode
	}

	if sslRootCert, ok := config["sslrootcert"].(string); ok {
		c.SSLRootCert = sslRootCert
	}

	if sslVerify, ok := config["sslverify"].(bool); ok {
		c.SSLVerify = sslVerify
	}

	if c.Host == "" {
		return nil, fmt.Errorf("missing required field: host")
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.User == "" {
		return nil, fmt.Errorf("missing required field: user")
	}
	if c.Password == "" {
		return nil, fmt.Errorf("missing required field: password")
	}
	if c.Database == "" {
		return nil, fmt.Errorf("missing required field: database")
	}

	return &c, nil
}

func (c *PostgresConnection) Close() error {
	return c.DB.Close()
}

func NewPostgresConnection(config PostgresConnectionConfig) (*PostgresConnection, error) {
	var db *sql.DB
	var err error

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
	)

	if config.SSLMode != "disable" && !config.SSLVerify {
		dsn += " sslrootcert= sslcert= sslkey="
	}

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PostgresConnection{
		DB: db,
	}, nil
}
