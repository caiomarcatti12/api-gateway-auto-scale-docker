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
package docker

import (
	"fmt"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/config"
	"log"
	"net/http"
	"time"
)

// checkHealth performs a health check for a specific route using Retry and Liveness Probe.
func checkHealth(route config.RouteConfig) bool {
	client := &http.Client{
		Timeout: time.Duration(route.Retry.Period) * time.Second,
	}

	// Extract Liveness Probe configuration
	liveness := route.LivenessProbe
	url := fmt.Sprintf("%s://%s:%d/%s", route.Backend.Protocol, route.Backend.Host, route.Backend.Port, route.LivenessProbe.Path)

	log.Printf("Performing health check for service: %s", route.Backend.ContainerName)

	// Initial delay defined in the Liveness Probe
	if liveness.InitialDelaySeconds > 0 {
		log.Printf("Waiting %d seconds before initial health check...", liveness.InitialDelaySeconds)
		time.Sleep(time.Duration(liveness.InitialDelaySeconds) * time.Second)
	}

	// Attempts defined in RetryConfig
	for attempt := 1; attempt <= route.Retry.Attempts; attempt++ {
		resp, err := client.Get(url)

		// Success check
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Printf("Health check succeeded for %s on attempt %d",
				route.Backend.ContainerName, attempt)
			return true
		}

		log.Printf("Attempt %d failed for %s, error: %v",
			attempt, route.Backend.ContainerName, err)

		// If not the last attempt, wait for the retry period
		if attempt < route.Retry.Attempts {
			log.Printf("Waiting %d seconds before the next attempt...", route.Retry.Period)
			time.Sleep(time.Duration(route.Retry.Period) * time.Second)
		}
	}

	// If all attempts fail, wait for the grace period before termination
	log.Printf("Health check failed for %s after %d attempts. Waiting %d seconds before finalizing...",
		route.Backend.ContainerName, route.Retry.Attempts, route.TTL)

	time.Sleep(time.Duration(route.TTL) * time.Second)
	return false
}
