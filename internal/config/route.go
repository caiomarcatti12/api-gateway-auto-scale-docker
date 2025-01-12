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

// HostConfig represents the configuration of a specific host.
type HostConfig struct {
	Host   string        `yaml:"host"`   // Host for which routes will be configured
	CORS   CORSConfig    `yaml:"cors"`   // CORS configuration specific to this host
	Routes []RouteConfig `yaml:"routes"` // List of routes for the host
}

// RouteConfig represents the configuration of a specific route.
type RouteConfig struct {
	Path          string              `yaml:"path"`          // Route path
	StripPath     bool                `yaml:"stripPath"`     // Indicates if the path should be removed
	TTL           int                 `yaml:"ttl"`           // Grace period for termination
	Backend       Backend             `yaml:"backend"`       // Backend configuration
	Retry         RetryConfig         `yaml:"retry"`         // Retry configuration
	LivenessProbe LivenessProbeConfig `yaml:"livenessProbe"` // Health check configuration
}

// Backend represents the backend configuration of a route.
type Backend struct {
	Protocol      string `yaml:"protocol"`      // Protocol (http or https)
	Host          string `yaml:"host"`          // Backend host
	Port          int    `yaml:"port"`          // Backend port
	ContainerName string `yaml:"containerName"` // Corresponding container name
}

// RetryConfig represents the retry configuration for a route.
type RetryConfig struct {
	Attempts int `yaml:"attempts"` // Number of retry attempts
	Period   int `yaml:"period"`   // Interval between retries in seconds
}

// LivenessProbeConfig represents the health check (liveness probe) configuration.
type LivenessProbeConfig struct {
	Path                string `yaml:"path"`                // Health check path
	SuccessThreshold    int    `yaml:"successThreshold"`    // Success threshold
	InitialDelaySeconds int    `yaml:"initialDelaySeconds"` // Initial delay before the health check
}
