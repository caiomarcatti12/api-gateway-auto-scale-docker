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
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/config"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/docker/container_store"
	"sync"
	"time"
)

var (
	containerMonitorMutex sync.Mutex
)

// CheckContainersToStop starts the continuous process of monitoring and stopping inactive containers.
func CheckContainersToStop() {
	for {
		monitorAndStopContainers()
		time.Sleep(5 * time.Second)
	}
}

// monitorAndStopContainers monitors and stops containers that are inactive beyond the timeout limit.
func monitorAndStopContainers() {
	containerMonitorMutex.Lock()
	defer containerMonitorMutex.Unlock()

	now := time.Now()

	hostStore := config.GetHostStore()

	hosts := hostStore.ListHosts()

	for _, host := range hosts {
		routes, _ := hostStore.GetAllRoutes(host)

		for _, route := range routes {
			container, _ := container_store.GetByContainerName(route.Backend.ContainerName)

			if container != nil {
				checkAndStopContainer(*container, route, now)
			}
		}
	}
}

// checkAndStopContainer checks if the container should be stopped based on TTL.
func checkAndStopContainer(container container_store.Container, route config.RouteConfig, now time.Time) {
	if isContainerExpired(container, route, now) {
		stopAndRemoveContainer(container)
	}
}

// isContainerExpired checks if the container has exceeded the allowed inactivity time.
func isContainerExpired(container container_store.Container, route config.RouteConfig, now time.Time) bool {
	return now.Sub(container.LastAccess) > time.Duration(route.TTL)*time.Second && container.IsActive
}

// stopAndRemoveContainer stops and removes the container from the store.
func stopAndRemoveContainer(container container_store.Container) {
	StopContainer(container.ID)

	container.IsActive = false
	container_store.Update(container)
}
