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

## How to Contribute

We are always open to contributions! If you'd like to help improve the project, whether through bug fixes, enhancements, or new features, follow our [Contribution Guide](docs/contributing.md) to understand the process and ensure your contribution is integrated smoothly.

## Code of Conduct

We are committed to providing a welcoming and inclusive community for everyone. We expect all project participants to follow our [Code of Conduct](docs/code_of_coduct.md). Please read and adhere to these guidelines to ensure a respectful and productive environment for all contributors.

## License

This project is licensed under the Apache 2.0 license. See the [LICENSE](docs/license) file for details.

