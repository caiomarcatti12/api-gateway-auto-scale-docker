## Development Environment Setup

To configure and start the development environment for the "API Gateway Auto Scale Docker" project, follow the steps below:

### 1. Prerequisites:
- Ensure you have [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine.

### 2. Clone the Repository:
Clone the repository to your local machine using the following command:
```bash
git clone https://github.com/caiomarcatti12/api-gateway-auto-scale-docker.git
cd api-gateway-auto-scale-docker
```

### 3. Configure Routes:
- Before starting the project, you need to configure the routes. Create a file named `config.yaml` in the project root.
- Use the `config.example.yaml` file as a reference for the structure of the `config.yaml` file.
- Configure the routes as per your requirements. If you need to add more routes or understand the configuration of the existing ones, refer to the [Route Configuration Guide](route_configuration.md).

### 4. Start the Project:
- With the `config.yaml` file configured, you can start the project using Docker Compose. Run the following command in the project root directory:
```bash
docker-compose up -d
```
This command will start all the services defined in the `docker-compose.yaml` file in "detached" mode, i.e., in the background.

### 5. Verification:
- After starting the services, you can check if all containers are running by using the command:
```bash
docker ps
```
- Now, the project should be running and ready to handle requests based on the configured routes.

### 6. Stop the Project:
- When you want to stop the project, use the following command:
```bash
docker-compose down
```

