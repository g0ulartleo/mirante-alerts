package connections

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

type MySQLConnectionConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Tunnel   TunnelConfig
}

type MySQLConnection struct {
	DB        *sql.DB
	sshClient *ssh.Client
}

func getIntValue(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		log.Fatalf("warning: expected value is not an int or float: %v", value)
		return 0
	}
}

func NewMySQLConnectionConfig(config map[string]any) (*MySQLConnectionConfig, error) {
	for _, field := range []string{"host", "port", "user", "password", "database"} {
		if _, ok := config[field]; !ok {
			return nil, fmt.Errorf("missing required field: %s", field)
		}
	}
	c := MySQLConnectionConfig{
		Host:     config["host"].(string),
		Port:     getIntValue(config["port"]),
		User:     config["user"].(string),
		Password: config["password"].(string),
		Database: config["database"].(string),
	}
	if tunnelConfig, ok := config["tunnel"].(map[string]any); ok {
		tunnel, err := NewTunnelConfig(tunnelConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create tunnel config: %v", err)
		}
		c.Tunnel = *tunnel
	}
	if c.Host == "" {
		return nil, fmt.Errorf("missing required field: host")
	}
	if c.Port == 0 {
		c.Port = 3306
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

func (c *MySQLConnection) Close() error {
	var dbErr, sshErr error
	if c.DB != nil {
		dbErr = c.DB.Close()
	}
	if c.sshClient != nil {
		sshErr = c.sshClient.Close()
	}
	if dbErr != nil {
		return dbErr
	}
	return sshErr
}

func NewMySQLConnection(config MySQLConnectionConfig) (*MySQLConnection, error) {
	var db *sql.DB
	var sshClient *ssh.Client
	var err error

	if config.Tunnel.Host != "" {
		sshClient, err = NewSSHClient(config.Tunnel)
		if err != nil {
			return nil, fmt.Errorf("failed to create SSH client: %v", err)
		}

		dialName := "mysql+ssh"
		mysql.RegisterDialContext(dialName, func(ctx context.Context, addr string) (net.Conn, error) {
			return sshClient.Dial("tcp", addr)
		})

		dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
			config.User,
			config.Password,
			dialName,
			config.Host,
			config.Port,
			config.Database,
		)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			sshClient.Close()
			return nil, fmt.Errorf("failed to connect to database via SSH tunnel: %v", err)
		}
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %v", err)
		}
	}

	return &MySQLConnection{
		DB:        db,
		sshClient: sshClient,
	}, nil
}
