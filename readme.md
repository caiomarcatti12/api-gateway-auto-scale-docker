# API Gateway Auto Scale Docker

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](docs/license) ![Static Badge](https://img.shields.io/badge/Not%20Ready%20for%20Production-red)

API Gateway Auto Scale Docker is an innovative solution that simulates Knative's functionality on Docker. Inspired by Knative's ability to manage serverless applications on Kubernetes, this project aims to provide an alternative for those using Docker in various environments, whether production, staging, or development.

## Key Features

- **Dynamic Container Management**: The API Gateway can dynamically start containers based on incoming requests. If a container is offline, it will be automatically started.

- **Resource Efficiency**: Containers that remain idle for a configurable period are shut down, optimizing resource usage.

- **Intelligent Routing**: The API Gateway manages request routing to the appropriate container, ensuring fast and efficient responses.

## API Gateway Route Configuration

Our project's API Gateway uses a specific structure to configure and manage routes that redirect requests to specific Docker containers. This configuration covers aspects such as route paths, target services, retry attempts, health checks, and more.

To fully understand how to configure and the expected behavior of these routes, refer to the detailed guide available in [Route Configuration](docs/route_configuration.md).

## Development Environment Setup

To configure and start the project's development environment, see the [Development Guide](docs/development.md).

## How to use docker run

```bash
docker run -d \
  --name api-gateway-auto-scale \
  -p 8080:8080 \
  -v $(pwd)/configs/config.yaml:/app/config.yaml \
  caiomarcatti12/api-gateway-auto-scale-docker:v0.0.6
```

- Certifique-se de criar o arquivo `configs/config.yaml` conforme o exemplo abaixo.

### How to use docker compse

```yaml
version: "3.8"
services:
  api-gateway-auto-scale:
    image: caiomarcatti12/api-gateway-auto-scale-docker:v0.0.6
    container_name: api-gateway-auto-scale
    ports:
      - "8080:8080"
    volumes:
      - ./configs/config.yaml:/app/config.yaml
    restart: unless-stopped
```

- Salve o arquivo acima como `docker-compose.yaml` e execute:
  ```bash
  docker compose up -d
  ```

### Exemplo de arquivo de configuração (`configs/config.yaml`)

```yaml
- host: host.docker.internal
  cors:
    allowedOrigins:
      - "http://host.docker.internal"
    allowedMethods:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
    allowedHeaders:
      - "Authorization"
      - "Content-Type"
    allowCredentials: true
    exposedHeaders:
      - "X-Custom-Header"
    maxAge: 3600
  routes:
    - path: /my-app-route
      stripPath: true
      ttl: 3
      backend:
        protocol: "http"
        host: "host.docker.internal"
        port: 8002
        containerName: "my-app-container-name"
      retry:
        attempts: 3
        period: 5
      livenessProbe:
        path: healthcheck
        successThreshold: 1
        initialDelaySeconds: 3
```

## How to Contribute

We are always open to contributions! If you'd like to help improve the project, whether through bug fixes, enhancements, or new features, follow our [Contribution Guide](docs/contributing.md) to understand the process and ensure your contribution is integrated smoothly.

## Code of Conduct

We are committed to providing a welcoming and inclusive community for everyone. We expect all project participants to follow our [Code of Conduct](docs/code_of_coduct.md). Please read and adhere to these guidelines to ensure a respectful and productive environment for all contributors.

## License

This project is licensed under the Apache 2.0 license. See the [LICENSE](docs/license) file for details.

