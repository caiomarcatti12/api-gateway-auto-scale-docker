# API Gateway Host and Route Configuration

This document outlines the new configuration structure for our API Gateway, which manages request redirection to specific Docker containers and monitors their health and lifecycle. The configuration has been updated to include **hosts** and **CORS**.

---

## YAML Configuration

The basic structure of the configuration file is as follows:

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

---

## Parameter Descriptions

### **HostConfig**
1. **host**: Defines the domain or IP of the host where routes will be configured.
2. **cors**: Configures **CORS** (Cross-Origin Resource Sharing) permissions:
    - **allowedOrigins**: List of allowed origins.
    - **allowedMethods**: Allowed HTTP methods (GET, POST, PUT, DELETE, etc.).
    - **allowedHeaders**: Allowed headers in requests.
    - **allowCredentials**: Specifies if credentials are allowed.
    - **exposedHeaders**: List of headers that can be exposed to the client.
    - **maxAge**: Maximum time, in seconds, that a CORS response can be cached.

### **RouteConfig**
1. **path**: Defines the route path for request redirection.
2. **stripPath**: Indicates whether the request path should be removed before redirection.
3. **ttl**: Specifies the maximum inactivity time, in seconds, before terminating the container.
4. **backend**: Contains the backend service configuration:
    - **protocol**: Protocol used (http or https).
    - **host**: Backend service's host or domain.
    - **port**: Port where the service is listening.
    - **containerName**: Name of the corresponding container.
5. **retry**: Configures retry attempts for unavailable services:
    - **attempts**: Maximum number of retry attempts.
    - **period**: Interval, in seconds, between retries.
6. **livenessProbe**: Configures the service's health check:
    - **path**: Path for the health check.
    - **successThreshold**: Minimum number of successful checks to consider the service healthy.
    - **initialDelaySeconds**: Initial waiting time before the first check.

---

## Behavior Based on Configuration

- **Hosts and Routes**:  
  The API Gateway redirects requests based on the **hosts** configuration. Each host can have multiple configured routes.

- **CORS**:  
  CORS permissions are configured for each host, ensuring that only allowed origins and methods can access resources.

- **Route Redirection**:  
  When a request arrives at a specific path, it is forwarded to the backend defined in the route configuration.

- **Health Checks**:  
  During startup, the API Gateway performs health checks on the specified path (`livenessProbe.path`).
    - If the check succeeds within the allowed attempts (`successThreshold`), the container is considered healthy.
    - Otherwise, the system will retry based on the `retry` configuration.

- **TTL (Time To Live)**:  
  If the container does not receive new requests within the configured time (`ttl`), it will be terminated.

- **Retry**:  
  If the container fails to start or becomes inaccessible, the API Gateway will retry according to the number and period defined in `retry`.

---

## Example

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

In this example:
1. The host `host.docker.internal` has a route at `/my-app-route` that redirects to the service `my-app-container-name`.
2. The backend uses the `http` protocol on port `8002`.
3. The container will be terminated after `3` seconds of inactivity (TTL).
4. If the service fails, the API Gateway will attempt to reconnect `3` times, with a `5`-second interval between attempts.
5. The health check will be performed at the path `/healthcheck` with an **initial delay** of 3 seconds and a **1-success tolerance** to consider it healthy.
6. CORS is configured to allow specific origins and methods, with response caching for up to 3600 seconds.

