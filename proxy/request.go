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
	"api-gateway-knative-docker/config"
	"api-gateway-knative-docker/cors"
	"api-gateway-knative-docker/docker"
	"net/http"
	"net/url"
	"strings"
)

func HandleRequest(route config.Route, corsGlobal *cors.CORSConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAllowed := true
		if route.CORS != nil {
			isAllowed = cors.ResolveCors(w, r, route.CORS) // set CORS headers
		} else if corsGlobal != nil {
			isAllowed = cors.ResolveCors(w, r, corsGlobal) // set CORS headers
		}

		if !isAllowed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if route.Protocol == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Check if it's just a preflight (OPTIONS) request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		_, err := docker.StartContainer(route)
		if err != nil {
			http.Error(w, "Error starting container", http.StatusInternalServerError)
			return
		}

		serviceURL := &url.URL{
			Scheme: route.Protocol,
			Host:   route.Host + ":" + route.Port,
		}

		// Strip the route path from the request
		if route.StripPath {
			r.URL.Path = stripRoutePath(r.URL.Path, route.Path)
		}

		proxyToService(serviceURL)(w, r)
	}
}

func stripRoutePath(requestPath, routePath string) string {
	return strings.TrimPrefix(requestPath, routePath)
}
