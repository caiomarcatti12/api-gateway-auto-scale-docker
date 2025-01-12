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
	"net/http"
	"strconv"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowedOrigins"`
	AllowedMethods   []string `yaml:"allowedMethods"`
	AllowedHeaders   []string `yaml:"allowedHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	ExposedHeaders   []string `yaml:"exposedHeaders"`
	MaxAge           int      `yaml:"maxAge"`
}

// ResolveCors configures CORS headers for an HTTP response based on the provided configuration.
func ResolveCors(w http.ResponseWriter, r *http.Request, corsConfig CORSConfig) bool {
	origin := r.Header.Get("Origin")

	if origin == "" {
		return true
	}

	isAllowed := false

	for _, allowedOrigin := range corsConfig.AllowedOrigins {
		if allowedOrigin == origin {
			isAllowed = true
			break
		}
	}

	if isAllowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	if len(corsConfig.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
	}
	if len(corsConfig.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
	}
	if len(corsConfig.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(corsConfig.ExposedHeaders, ", "))
	}
	if corsConfig.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if corsConfig.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAge))
	}

	return isAllowed
}
