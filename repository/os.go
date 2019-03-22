package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	// LocalConfig is path to gamesave local Git repository
	LocalConfig string = ".gamesave.json"
)

// IOSRepository is interface for interaction with local files
// include configuration files
type IOSRepository interface {
	Copy(src, dst string) error
	GetConfig(key string) string
	SetConfig(key, value string) error
}

// OSRepository is the implementation of IOSRepository
type OSRepository struct{}

// Copy force copies file or directory from src to dst
func (rep *OSRepository) Copy(src, dst string) error {
	cmd := exec.Command("cp", "-rf", src, dst)
	output, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Print(string(output))
	} else {
		err = errors.New(string(output))
	}
	return err
}

// GetConfig get config by the key from LocalConfig
// returns empty string if key not exist
func (rep *OSRepository) GetConfig(key string) string {
	var config map[string]string
	data, err := ioutil.ReadFile(LocalConfig)
	if err != nil {
		fmt.Printf("Error on get config: %v", err)
		return ""
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error on get config: %v", err)
		return ""
	}
	if val, ok := config[key]; ok {
		return val
	}
	return ""
}

// SetConfig set config by the key from LocalConfig
// overwrite value of existing key
func (rep *OSRepository) SetConfig(key, value string) error {
	var config map[string]string
	err := createIfNotExist()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(LocalConfig)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	config[key] = value
	byt, err := json.MarshalIndent(config, "", "  ")
	return ioutil.WriteFile(LocalConfig, byt, 0644)
}

func createIfNotExist() error {
	if _, err := os.Stat(LocalConfig); os.IsNotExist(err) {
		return ioutil.WriteFile(LocalConfig, []byte("{}"), 0644)
	}
	return nil
}
