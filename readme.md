### Project Overview
This is a one-stop microservices example project based on **wego**, covering both development and deployment. It is suitable for lightweight mobile apps or small game backends.

#### Drivers
- Database: **pgx** (PostgreSQL driver)
- Cache: **go-redis** (Redis client)

#### Middleware & Components
- Cache: **Redis**
- Database: **PostgreSQL**
- Service discovery & configuration: **etcd**

### Deployment
- Containerization: **Docker images**
- Orchestration: **Kubernetes** (container management)

### Implemented Modules
- [x] **API**
- [x] **account** service

### Testing
1. **Local debugging**: Directly run the configurations in `launch.json` on your local machine using VSCode
2. **Local Kubernetes testing**: Use **k3d** or **minikube**

