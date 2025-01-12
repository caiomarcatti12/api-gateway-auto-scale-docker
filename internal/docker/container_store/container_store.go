package container_store

import (
	"time"
)

var containers = make(map[string]Container)
var containersBySvc = make(map[string]Container)

func Add(container Container) {
	containers[container.ID] = container
	containersBySvc[container.ContainerName] = container
}
func Update(container Container) {
	containers[container.ID] = container
	containersBySvc[container.ContainerName] = container
}

func Remove(containerID string) {
	if container, exists := containers[containerID]; exists {
		delete(containersBySvc, container.ContainerName)
		delete(containers, containerID)
	}
}

func GetByID(containerID string) (Container, bool) {
	container, exists := containers[containerID]
	return container, exists
}

func GetByContainerName(serviceName string) (*Container, bool) {
	container, exists := containersBySvc[serviceName]
	if !exists {
		return nil, false
	}
	return &container, true
}

func UpdateAccessTime(containerID string) {
	if container, exists := containers[containerID]; exists {
		container.LastAccess = time.Now()
		containers[containerID] = container
		containersBySvc[container.ContainerName] = container
	}
}

func GetAll() map[string]Container {
	return containers
}
