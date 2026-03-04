### 项目简介
这是一个基于go-micro的一站式，包含开发、部署的微服务项目，可用于轻量型App、小游戏后端;

如要使用，具体的业务代码自行复用或者剥离;

### 技术栈
#### framework
- reverse proxy: traefik
- api: gin
- server: go-micro
#### driver
- pgx PostgreSQL
- go-redis Redis

### 中间件
- 缓存： Redis
- db: PostgreSQL
- ETCD

### deployment
- docker image
- k8s container manager

### api
- [x] api

### servers
- [x] account

### TEST
1. 本地debug，直接启动launch.json中的配置
2. 本地k8s集群测试，k3d/minikube

### go package
```bash
go install github.com/micro/micro/v5/cmd/protoc-gen-micro@latest
```
