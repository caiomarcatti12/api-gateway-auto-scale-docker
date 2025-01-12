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
package proxy

import (
	"fmt"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/config"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/docker"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/docker/container_store"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// HandleRequest processes an incoming request and routes it to the appropriate backend service.
func HandleRequest(route config.RouteConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if route.Backend.Protocol == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Check if it's just a preflight (OPTIONS) request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if route.Backend.ContainerName != "" {
			containerService, exists := container_store.GetByContainerName(route.Backend.ContainerName)

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if !containerService.IsActive {
				_, err := docker.StartContainer(route)
				if err != nil {
					http.Error(w, "Error starting container", http.StatusInternalServerError)
					return
				}
			}

			log.Printf("Last access to the service container %s updated.", route.Backend.ContainerName)
			container_store.UpdateAccessTime(containerService.ID)
		}

		serviceURL := &url.URL{
			Scheme: route.Backend.Protocol,
			Host:   fmt.Sprintf("%s:%d", route.Backend.Host, route.Backend.Port),
		}

		// Strip the route path from the request
		if route.StripPath {
			r.URL.Path = stripRoutePath(r.URL.Path, route.Path)
		}

		proxyToService(serviceURL)(w, r)
	}
}

// stripRoutePath removes the route's base path from the request path.
func stripRoutePath(requestPath, routePath string) string {
	return strings.TrimPrefix(requestPath, routePath)
}
