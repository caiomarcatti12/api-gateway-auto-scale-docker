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
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/docker/container_store"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var (
	updateContainerMutex sync.Mutex
	dockerClientInstance *client.Client
)

// CheckContainersActive starts the continuous process of verifying the containers.
func CheckContainersActive() {
	for {
		syncContainersState()
		time.Sleep(5 * time.Second)
	}
}

// syncContainersState is the main process that synchronizes the state of the containers.
func syncContainersState() {
	updateContainerMutex.Lock()
	defer updateContainerMutex.Unlock()

	cli, err := getDockerClient()
	if err != nil {
		log.Println("Error obtaining Docker client:", err)
		return
	}

	containers, err := listAllContainers(cli)
	if err != nil {
		log.Println("Error listing containers:", err)
		return
	}

	currentContainers := mapContainers(containers)
	activeContainers := container_store.GetAll()

	removeMissingContainers(activeContainers, currentContainers)
	updateOrAddContainers(activeContainers, currentContainers)
}

// listAllContainers lists all containers, including stopped ones.
func listAllContainers(cli *client.Client) ([]types.Container, error) {
	ctx := context.Background()
	return cli.ContainerList(ctx, container.ListOptions{All: true})
}

// mapContainers creates a map of the current containers with their relevant information.
func mapContainers(containers []types.Container) map[string]container_store.Container {
	currentContainers := make(map[string]container_store.Container)

	for _, container := range containers {
		for _, name := range container.Names {
			newContainer := createContainerObject(container, name)
			currentContainers[container.ID] = newContainer
		}
	}
	return currentContainers
}

// createContainerObject creates a Container instance based on the provided data.
func createContainerObject(container types.Container, name string) container_store.Container {
	return container_store.Container{
		ID:            container.ID,
		ContainerName: strings.ReplaceAll(name, "/", ""),
		LastAccess:    time.Now(),
		IsActive:      container.State == "running",
	}
}

// removeMissingContainers removes containers that are no longer present on the host.
func removeMissingContainers(activeContainers, currentContainers map[string]container_store.Container) {
	for containerID, storedContainer := range activeContainers {
		if _, exists := currentContainers[containerID]; !exists {
			container_store.Remove(containerID)
			log.Printf("Removed container: %s (%s)", storedContainer.ContainerName, storedContainer.ID)
		}
	}
}

// updateOrAddContainers adds or updates containers in the store.
func updateOrAddContainers(activeContainers, currentContainers map[string]container_store.Container) {
	for _, currentContainer := range currentContainers {
		if storedContainer, exists := activeContainers[currentContainer.ID]; exists {
			updateContainerIfChanged(storedContainer, currentContainer)
		} else {
			addNewContainer(currentContainer)
		}
	}
}

// updateContainerIfChanged updates a container if there is a change in its status.
func updateContainerIfChanged(storedContainer, currentContainer container_store.Container) {
	if storedContainer.IsActive != currentContainer.IsActive {
		storedContainer.IsActive = currentContainer.IsActive

		container_store.Update(storedContainer)

		log.Printf("Updated container: %s (%s) - IsActive: %v",
			storedContainer.ContainerName, storedContainer.ID, storedContainer.IsActive)
	}
}

// addNewContainer adds a new container to the store.
func addNewContainer(currentContainer container_store.Container) {
	container_store.Add(currentContainer)
	log.Printf("Added new container: %s (%s)", currentContainer.ContainerName, currentContainer.ID)
}
