/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ConfigLoader struct {
	configDir string
}

// NewConfigLoader initializes a new ConfigLoader with the determined configuration directory.
func NewConfigLoader() (*ConfigLoader, error) {
	configDir, err := determineConfigDir()
	if err != nil {
		return nil, err
	}
	return &ConfigLoader{configDir: configDir}, nil
}

// LoadConfigs loads and parses configuration files from the config directory.
func (cl *ConfigLoader) LoadConfigs() error {
	files, err := cl.getConfigFiles()
	if err != nil {
		return err
	}

	configs, err := cl.parseConfigFiles(files)
	if err != nil {
		return err
	}

	if len(configs) == 0 {
		return errors.New("no config files found")
	}

	for _, config := range configs {
		GetHostStore().AddHost(config)
	}

	return nil
}

// determineConfigDir determines the configuration directory based on environment or defaults.
func determineConfigDir() (string, error) {
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		return envConfigPath, nil
	}

	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting binary path: %s", err.Error())
	}

	execDir := filepath.Dir(executable)
	if filepath.Base(os.Args[0]) == "main" {
		return "configs", nil
	}

	return execDir, nil
}

// getConfigFiles retrieves all YAML configuration files from the configuration directory.
func (cl *ConfigLoader) getConfigFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(cl.configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".yaml" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error listing configuration files: %s", err.Error())
	}

	return files, nil
}

// parseConfigFiles parses the content of configuration files into HostConfig objects.
func (cl *ConfigLoader) parseConfigFiles(files []string) ([]HostConfig, error) {
	var configs []HostConfig

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %s", file, err.Error())
		}

		var fileConfigs []HostConfig
		if err := yaml.Unmarshal(content, &fileConfigs); err != nil {
			return nil, fmt.Errorf("error when deserializing the file %s: %s", file, err.Error())
		}

		configs = append(configs, fileConfigs...)
	}

	return configs, nil
}
