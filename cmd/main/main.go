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
package main

import (
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/config"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/docker"
	"github.com/caiomarcatti12/api-gateway-auto-scale-docker/internal/proxy"
	"log"
	"net/http"
)

func main() {
	configLoader, err := config.NewConfigLoader()

	if err != nil {
		log.Fatal("Error loading config.")
		return
	}

	err = configLoader.LoadConfigs()

	if err != nil {
		log.Fatal("Error loading config.")
		return
	}

	// Defina um manipulador padrão para "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		routeConfig, exists := config.GetHostStore().GetRoute(r.Host, r.URL.Path)

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		corsConfig, exists := config.GetHostStore().GetCORS(r.Host)

		if exists {
			isAllowed := config.ResolveCors(w, r, corsConfig)

			if !isAllowed {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		proxy.HandleRequest(routeConfig)(w, r)
	})

	go docker.CheckContainersActive()
	go docker.CheckContainersToStop()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
