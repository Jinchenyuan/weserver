### 技术栈
#### framework
- gateway: traefik (api, https, tsl)
- server: gin(用于api通信) + go-micro(实现为服务、服务发现；用于服务通信)
- admin: go-admin
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

### servers
- [ ] admin
- [ ] gateway
- [ ] map
- [ ] account
- [x] s3

### TEST
1. 本地debug，直接启动launch.json中的配置
2. 本地k8s集群测试，k3d/minikube