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
	"context"
	"errors"
	"fmt"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/config"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/docker/container_store"
	"log"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var (
	mutexes      = make(map[string]*sync.Mutex)
	mutexesGuard = &sync.Mutex{} // Guard para proteger o acesso ao mapa de mutexes
	once         sync.Once
)

// getDockerClient garante que apenas uma instância do cliente Docker seja criada (singleton).
func getDockerClient() (*client.Client, error) {
	var err error
	once.Do(func() {
		dockerClientInstance, err = client.NewClientWithOpts(client.FromEnv)
	})
	return dockerClientInstance, err
}

// getMutexForService retorna o mutex associado a um serviço, criando um novo se necessário.
func getMutexForService(service string) *sync.Mutex {
	mutexesGuard.Lock()
	defer mutexesGuard.Unlock()

	if _, exists := mutexes[service]; !exists {
		log.Printf("creating new mutex for the service: %s", service)
		mutexes[service] = &sync.Mutex{}
	}
	return mutexes[service]
}

// StartContainer Funcionalidade de iniciar um container
func StartContainer(route config.RouteConfig) (bool, error) {
	if route.Backend.ContainerName == "" {
		log.Println("No services associated with the route, ignoring container start.")
		return true, nil
	}

	ctx := context.Background()
	cli, err := getDockerClient()

	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return false, err
	}

	log.Printf("Starting process for the service container: %s", route.Backend.ContainerName)

	containerService, exists := container_store.GetByContainerName(route.Backend.ContainerName)

	if !exists {
		log.Printf("Unable to find service for container %s", route.Backend.ContainerName)
	}

	serviceMutex := getMutexForService(route.Backend.ContainerName)
	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	log.Printf("Container for service %s is not running. Trying to start...", route.Backend.ContainerName)
	if err := cli.ContainerStart(ctx, containerService.ID, container.StartOptions{}); err != nil {
		log.Printf("Error starting container for service %s: %v", route.Backend.ContainerName, err)
		return false, err
	}

	log.Printf("Container started for service: %s", route.Backend.ContainerName)

	// Verificar o healthcheck do container
	if !checkHealth(route) {
		log.Printf("Healthcheck failed for container %s", route.Backend.ContainerName)
		return false, errors.New(fmt.Sprintf("Healthcheck failed for container %s", route.Backend.ContainerName))
	}

	log.Printf("Healthcheck successful for container: %s", route.Backend.ContainerName)

	log.Printf("Last access to updated service container %s.", route.Backend.ContainerName)
	container_store.UpdateAccessTime(containerService.ID)

	return true, nil
}

// StopContainer Funcionalidade de parar um container
func StopContainer(containerID string) {
	ctx := context.Background()
	cli, err := getDockerClient()
	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return
	}

	log.Printf("Starting stop process for container: %s", containerID)

	// Recupera o serviço associado ao containerID para obter o mutex correto
	service := getServiceForContainer(containerID)
	if service == "" {
		log.Printf("Error finding the service associated with the container: %s", containerID)
		return
	}

	serviceMutex := getMutexForService(service)
	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	log.Printf("Stopping container: %s of service: %s", containerID, service)
	err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v", containerID, err)
	} else {
		log.Printf("Container %s stopped successfully.", containerID)
	}
}

// getServiceForContainer é um placeholder para obter o serviço associado ao containerID
func getServiceForContainer(containerID string) string {
	containerInStore, exists := container_store.GetByID(containerID)
	if exists {
		return containerInStore.ContainerName
	}
	log.Printf("Unable to find service for container %s", containerID)
	return ""
}
