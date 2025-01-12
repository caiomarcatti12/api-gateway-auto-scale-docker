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
	"fmt"
	"strings"
	"sync"
)

// HostStore is the main storage for hosts and routes.
type HostStore struct {
	store map[string]HostData
}

// HostData stores the routes and CORS configuration for each host.
type HostData struct {
	CORS   CORSConfig             // CORS configuration specific to the host
	Routes map[string]RouteConfig // Mapping of routes by path
}

var (
	once     sync.Once
	instance *HostStore
)

// GetHostStore returns the Singleton instance of HostStore.
func GetHostStore() *HostStore {
	once.Do(func() {
		instance = &HostStore{
			store: make(map[string]HostData),
		}
	})
	return instance
}

// AddHost adds or updates a host in the HostStore.
func (hs *HostStore) AddHost(hostConfig HostConfig) {
	// Create the route map for the host.
	routeMap := make(map[string]RouteConfig)
	for _, route := range hostConfig.Routes {
		routeMap[route.Path] = route
	}

	// Add the host with its routes and CORS configuration.
	hs.store[hostConfig.Host] = HostData{
		CORS:   hostConfig.CORS,
		Routes: routeMap,
	}
}

// GetRoute retrieves a specific route of a host by its path.
func (hs *HostStore) GetRoute(host, path string) (RouteConfig, bool) {
	hostData, ok := hs.store[host]
	if !ok {
		return RouteConfig{}, false
	}

	prefix := hs.getPrefixPath(path)

	route, found := hostData.Routes[prefix]

	return route, found
}

// GetAllRoutes retrieves all routes of a specific host.
func (hs *HostStore) GetAllRoutes(host string) ([]RouteConfig, bool) {
	hostData, ok := hs.store[host]
	if !ok {
		return nil, false
	}

	// Convert the route map to a slice.
	routes := make([]RouteConfig, 0, len(hostData.Routes))
	for _, route := range hostData.Routes {
		routes = append(routes, route)
	}
	return routes, true
}

// GetCORS retrieves the CORS configuration of a host.
func (hs *HostStore) GetCORS(host string) (CORSConfig, bool) {
	hostData, ok := hs.store[host]
	if !ok {
		return CORSConfig{}, false
	}
	return hostData.CORS, true
}

// ListHosts returns all stored hosts.
func (hs *HostStore) ListHosts() []string {
	hosts := make([]string, 0, len(hs.store))
	for host := range hs.store {
		hosts = append(hosts, host)
	}
	return hosts
}

func (hs *HostStore) getPrefixPath(path string) string {
	prefixSplit := strings.Split(path, "/")

	if len(prefixSplit) > 1 {
		return fmt.Sprintf("/%s", prefixSplit[1])
	}
	return ""
}
