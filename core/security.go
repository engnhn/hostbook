package core

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	ServiceName = "hostbook"
)

func SavePassword(hostName, password string) error {
	if password == "" {
		return nil
	}
	return keyring.Set(ServiceName, hostName, password)
}

func GetPassword(hostName string) (string, error) {
	password, err := keyring.Get(ServiceName, hostName)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to get password: %w", err)
	}
	return password, nil
}

func DeletePassword(hostName string) error {
	err := keyring.Delete(ServiceName, hostName)
	if err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete password: %w", err)
	}
	return nil
}
