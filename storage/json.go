package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/engnhn/hostbook/core"
)

const (
	configDirName = ".hostbook"
	hostsFileName = "hosts.json"
	sshConfigName = "ssh_config"
)

type Storage struct {
	mu            sync.Mutex
	filePath      string
	sshConfigPath string
}

func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, configDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	socketsDir := filepath.Join(configDir, "sockets")
	if err := os.MkdirAll(socketsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create sockets directory: %w", err)
	}

	return &Storage{
		filePath:      filepath.Join(configDir, hostsFileName),
		sshConfigPath: filepath.Join(configDir, sshConfigName),
	}, nil
}

func (s *Storage) LoadHosts() ([]core.Host, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []core.Host{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %w", err)
	}

	var hosts []core.Host
	if err := json.Unmarshal(data, &hosts); err != nil {
		return nil, fmt.Errorf("failed to parse hosts file: %w", err)
	}

	return hosts, nil
}

func (s *Storage) SaveHosts(hosts []core.Host) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(hosts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal hosts: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write hosts file: %w", err)
	}

	return nil
}

func (s *Storage) SaveSSHConfig(content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.WriteFile(s.sshConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write ssh config file: %w", err)
	}

	return nil
}

func (s *Storage) GetSSHConfigPath() string {
	return s.sshConfigPath
}
