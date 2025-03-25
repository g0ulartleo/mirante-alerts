package connections

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type TunnelConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	PemFile  string
}

func NewTunnelConfig(config map[string]any) (*TunnelConfig, error) {
	for _, field := range []string{"host", "port", "user"} {
		if _, ok := config[field]; !ok {
			return nil, fmt.Errorf("missing required field: %s", field)
		}
	}
	c := TunnelConfig{
		Host: config["host"].(string),
		Port: getIntValue(config["port"]),
		User: config["user"].(string),
	}
	if _, ok := config["password"]; ok {
		c.Password = config["password"].(string)
	}
	if _, ok := config["pem_file"]; ok {
		c.PemFile = config["pem_file"].(string)
	}
	if c.Host != "" {
		if c.PemFile == "" && c.Password == "" {
			return nil, fmt.Errorf("missing required field: pem_file or password")
		}
		if c.User == "" {
			return nil, fmt.Errorf("missing required field: user")
		}
		if c.Port == 0 {
			c.Port = 22
		}
	}
	return &c, nil
}

func NewSSHClient(config TunnelConfig) (*ssh.Client, error) {
	tunnelAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	var sshConfig *ssh.ClientConfig

	if config.PemFile != "" {
		key, err := os.ReadFile(config.PemFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read PEM file: %v", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
		sshConfig = &ssh.ClientConfig{
			User:            config.User,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         40 * time.Second,
		}
	} else {
		sshConfig = &ssh.ClientConfig{
			User:            config.User,
			Auth:            []ssh.AuthMethod{ssh.Password(config.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         40 * time.Second,
		}
	}

	sshClient, err := ssh.Dial("tcp", tunnelAddr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH tunnel: %v", err)
	}
	return sshClient, nil
}
