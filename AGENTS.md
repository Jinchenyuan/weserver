# AGENTS.md

This repository is the backend service project for `we`. It currently centers on two implemented business areas, `account` and `storyline`, and is built around the `wego` runtime plus a small set of clear layers.

## Project Overview

- Entry points:
  - `api/storyline/main.go`: unified HTTP API process for `account` + `storyline`
  - `service/account/main.go`: microservice process
  - `service/storyline/main.go`: microservice process
- Main request path:
  - Gin HTTP route -> API handler -> go-micro client -> service handler -> model -> PostgreSQL
- Shared runtime:
  - `wego` creates and owns the HTTP server, microservice server, PostgreSQL connection, Redis client, etcd registry, and global logger
- Current business scope:
  - Implemented and usable: `account`, `storyline`
  - Repository still contains some legacy or placeholder `s3` references in Swagger, scripts, and deploy files

## Main Areas

- `api`
  - API gateway layer
  - Receives HTTP requests, binds request bodies, calls microservice clients, and shapes JSON responses
- `service`
  - Business service layer
  - Registers go-micro handlers and implements business logic
- `model`
  - Persistence layer built on `bun`
  - Defines tables and basic CRUD/query helpers
- `protobuf`
  - `account.proto`, `storyline.proto`, plus generated Go/micro stubs
  - API and service communicate through these contracts
- `config`
  - Shared TOML config loader
  - Both API and service read `config.toml` from their own working directory
- `utils`
  - Password hashing and JWT helpers
- `deploy`
  - Kubernetes manifests and related deployment resources
- `service/schemes/postgre`
  - SQL schema bootstrap scripts for PostgreSQL

## Runtime Model

- `wego.New(...)` in both entry points initializes:
  - etcd registry client
  - PostgreSQL `bun.DB`
  - Redis client
  - HTTP server
  - go-micro server
  - global Mesa runtime
- The unified API process:
  - runs from `api/storyline/main.go`
  - configures JWT-based auth middleware for both `account` and `storyline`
  - registers Gin routes from both `api/account/ginhandler` and `api/storyline/ginhandler`
  - creates service clients from both `api/account/serviceclient` and `api/storyline/serviceclient`
- The service process:
  - configures the microservice identity from `[service]`
  - registers RPC handlers in `service/account/servicehandler/registry.go` or `service/storyline/servicehandler/registry.go`

## Account Module

- HTTP routes:
  - `POST /account/login`
  - `POST /account/register`
  - `GET /account/hello`
  - Swagger is mounted at `/swagger/*any`
- RPC contract:
  - `protobuf/pb/account.proto`
  - methods: `Hello`, `Login`, `Register`
- Persistence:
  - `model/account.go`
  - table DDL: `service/schemes/postgre/accounts.sql`
- Auth flow:
  - login validates the bcrypt password hash
  - service generates JWT via `utils.GenerateToken`
  - token is cached in Redis as `token:<accountID>`
  - API middleware checks `Authorization: Bearer <token>` plus request header `account`
- Important caveat:
  - `utils.GenerateToken` sets JWT expiry to 24 hours, but login writes the Redis token with a 7 day TTL. Actual validity is therefore bounded by the JWT expiry unless token parsing behavior is changed elsewhere.

## Storyline Module

- HTTP routes:
  - `GET /storylines`
  - `GET /storylines/:id`
  - `POST /storylines`
  - `PUT /storylines/:id`
- RPC contract:
  - `protobuf/pb/storyline.proto`
  - methods: `ListStorylines`, `GetStoryline`, `CreateStoryline`, `UpdateStoryline`
- Persistence:
  - `model/storyline.go`
  - table DDL: `service/schemes/postgre/storylines.sql`
- Data model:
  - one `storylines` row per storyline
  - many `storyline_nodes` rows per storyline
  - both storyline id and node id are stored as UUID strings to match frontend `String` ids
- Auth and ownership:
  - requests authenticate by Bearer JWT only
  - API middleware parses `account_id` from JWT and validates the token against Redis
  - all storyline reads and writes are filtered by `account_id`
- Behavior contract:
  - API JSON shape must match `we_client/lib/features/storyline/models/storyline_models.dart`
  - node order is normalized by `sortOrder`
  - list sorting follows latest node date first, then `updatedAt`
  - `coverPhotoUri` and `photoUri` are stored as raw strings, including data URIs

## Config Files

- API template: `api/config.toml.template`
- Service template: `service/config.toml.template`
- Typical required dependencies:
  - PostgreSQL
  - Redis
  - etcd
- Important fields:
  - `[http].excludeAuthPaths`: unauthenticated API routes
  - `[service]`: microservice name/version/port for service registration
  - `[services].account`: target service name used by the API client
  - `[services].storyline`: target service name used by the API client
  - `[profile].name`: runtime profile / logger identity

## Local Run Notes

- Simple local start from the repo root:
  - `./run.sh`
- That script starts:
  - unified `api/storyline` gateway
  - `service/account`
  - `service/storyline`
- Each process expects a local `config.toml` in its own directory, so copy from:
  - `api/storyline/config.toml.template` -> `api/storyline/config.toml`
  - `service/config.toml.template` -> `service/account/config.toml`
  - `service/storyline/config.toml.template` -> `service/storyline/config.toml`
- Database schema setup:
  - run `service/schemes/postgre/aa_init.sql`
  - then run `service/schemes/postgre/accounts.sql`
  - then run `service/schemes/postgre/storylines.sql`

## Change Guidance

- When adding or changing an API for an existing module, update in this order:
  - `protobuf/pb/*.proto`
  - regenerate code via `protobuf/proto.sh`
  - `service/.../servicehandler`
  - `api/.../serviceclient`
  - `api/.../ginhandler`
  - `model/...` if storage changes are needed
- When adding a brand new module, mirror the existing `account` shape:
  - add a new proto service definition
  - add a new `service/<module>` entry point and handler registry
  - add a new `api/<module>` route registry and service client
  - add schema/model files
  - add config entries for the downstream service name
- Prefer keeping transport and business concerns separate:
  - HTTP validation/response shaping stays in `ginhandler`
  - business logic stays in `servicehandler`
  - DB logic stays in `model`

## Testing Notes

- Current meaningful automated tests are limited:
  - `service/account/servicehandler/account_test.go` checks logger fallback behavior
  - `service/storyline/servicehandler/storyline_test.go` checks logger fallback and basic validation
  - `api/storyline/ginhandler/storyline_test.go` checks request validation
- `model/account_test.go` is integration-style and depends on a real PostgreSQL instance with hard-coded DSNs; treat it as manual or environment-specific unless it is rewritten
- When changing login/register behavior, add tests close to the service handler first
- When changing storyline request/response shape, validate against the frontend Dart models before changing backend field names

## Working Notes

- `wego` is not just a helper library; it is the backend runtime foundation for server startup, transport wiring, registry, DB/Redis lifecycle, and shared middleware behavior
- The repo currently looks like an evolving template plus one implemented business module; inspect for stale references before assuming a feature is live
- Safe default for new work:
  - trace the request through `api/storyline/main.go` -> feature `ginhandler` -> feature `serviceclient` -> feature `service/.../main.go` -> `servicehandler` -> `model`
- Be careful with service naming:
  - the API client resolves the downstream service by `cfg.Services.Account`
  - the service registers itself by `cfg.Service.Name`
  - these two names must match in practice
